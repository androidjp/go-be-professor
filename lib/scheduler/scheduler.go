package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
)

// RTaskElement 定义定时任务结构
type RTaskElement struct {
	Key       string    `json:"key"`
	ExecuteAt time.Time `json:"execute_at"`
	CronExpr  string    `json:"cron_expr"`
}

// RTimeWheel 定义时间轮结构
type RTimeWheel struct {
	ticker       *time.Ticker
	redisClient  *redis.Client
	httpClient   *http.Client
	stopc        chan struct{}
	schedulerKey string
	taskChan     chan RTaskElement
	maxWorkers   int
	timeout      time.Duration
}

// NewRTimeWheel 创建时间轮实例
func NewRTimeWheel(redisClient *redis.Client, httpClient *http.Client, schedulerKey string, maxWorkers int, timeout time.Duration) *RTimeWheel {
	r := &RTimeWheel{
		ticker:       time.NewTicker(time.Second),
		redisClient:  redisClient,
		httpClient:   httpClient,
		stopc:        make(chan struct{}),
		schedulerKey: schedulerKey,
		taskChan:     make(chan RTaskElement, maxWorkers),
		maxWorkers:   maxWorkers,
		timeout:      timeout,
	}
	go r.run()
	return r
}

// run 启动时间轮
func (r *RTimeWheel) run() {
	// 启动工作协程池
	for i := 0; i < r.maxWorkers; i++ {
		go r.worker()
	}

	for {
		select {
		case <-r.ticker.C:
			r.processTasks()
		case <-r.stopc:
			r.ticker.Stop()
			close(r.taskChan)
			return
		}
	}
}

// worker 工作协程，处理任务队列中的任务
func (r *RTimeWheel) worker() {
	for task := range r.taskChan {
		done := make(chan struct{})
		go func(t RTaskElement) {
			r.executeTask(t)
			close(done)
		}(task)

		select {
		case <-done:
			// 任务正常完成
		case <-time.After(r.timeout):
			// 任务超时，放弃该任务
			log.Printf("Task %s timeout after %v", task.Key, r.timeout)
		}
	}
}

// processTasks 处理定时任务
func (r *RTimeWheel) processTasks() {
	ctx := context.Background()
	now := time.Now().UnixNano()

	// 使用Lua脚本原子性地获取并移除到期任务
	result, err := r.redisClient.Eval(ctx, LuaGetAndRemoveTasks,
		[]string{r.schedulerKey},
		fmt.Sprintf("%d", now),
		"10", // 每次最多处理10个任务
	).Result()

	if err != nil {
		log.Printf("Failed to get and remove tasks from Redis: %v", err)
		return
	}

	tasks, ok := result.([]interface{})
	if !ok || len(tasks) == 0 {
		return
	}

	// 处理获取到的任务
	for i := 0; i < len(tasks); i += 2 {
		taskStr, ok := tasks[i].(string)
		if !ok {
			continue
		}

		var task RTaskElement
		if err := json.Unmarshal([]byte(taskStr), &task); err != nil {
			log.Printf("Failed to unmarshal task: %v", err)
			continue
		}

		select {
		case <-r.stopc:
			return
		case r.taskChan <- task: // 将任务发送到协程池处理
		default:
			// 如果任务队列已满，将任务重新加入Redis
			taskBody, _ := json.Marshal(task)
			_, err = r.redisClient.Eval(ctx, LuaAddTasks,
				[]string{r.schedulerKey},
				task.ExecuteAt.UnixNano(),
				string(taskBody),
			).Result()
			if err != nil {
				log.Printf("Failed to readd task %s: %v", task.Key, err)
			}
			log.Printf("Task channel is full, task %s has been readded", task.Key)
		}
	}
}

// executeTask 执行定时任务
func (r *RTimeWheel) executeTask(task RTaskElement) {
	ctx := context.Background()

	// 执行任务前先获取分布式锁
	lockKey := fmt.Sprintf("%s:lock:%s", r.schedulerKey, task.Key)
	lockValue := time.Now().String()
	locked, err := r.redisClient.SetNX(ctx, lockKey, lockValue, r.timeout).Result()
	if err != nil || !locked {
		log.Printf("Failed to acquire lock for task %s: %v", task.Key, err)
		return
	}
	defer r.redisClient.Del(ctx, lockKey)
	log.Printf("Executing task: %s at %v", task.Key, time.Now().Format("2006-01-02 15:04:05"))

	// 这里可以添加具体的任务执行逻辑，例如HTTP请求等
	// 模拟任务执行
	time.Sleep(100 * time.Millisecond)

	// 如果是cron任务，计算下一次执行时间并重新添加到Redis
	if task.CronExpr != "" {
		parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		schedule, err := parser.Parse(task.CronExpr)
		if err != nil {
			log.Printf("Failed to parse cron expression: %v", err)
			return
		}

		// 计算下一次执行时间
		nextTime := schedule.Next(time.Now())
		task.ExecuteAt = nextTime

		// 将任务重新添加到Redis
		taskBody, err := json.Marshal(task)
		if err != nil {
			log.Printf("Failed to marshal task: %v", err)
			return
		}

		// 使用Lua脚本原子性地添加下一次执行时间
		_, err = r.redisClient.Eval(ctx, LuaAddTasks, []string{r.schedulerKey}, nextTime.UnixNano(), string(taskBody)).Result()
		if err != nil {
			log.Printf("Failed to add next cron task to Redis: %v", err)
			// 重试添加任务
			for i := 0; i < 3; i++ {
				time.Sleep(time.Second)
				_, err = r.redisClient.Eval(ctx, LuaAddTasks, []string{r.schedulerKey}, nextTime.UnixNano(), string(taskBody)).Result()
				if err == nil {
					break
				}
			}
			if err != nil {
				log.Printf("Failed to add next cron task after retries: %v", err)
			}
		}
	}
}

// AddTask 添加定时任务
func (r *RTimeWheel) AddTask(ctx context.Context, key string, task *RTaskElement) error {
	if err := r.addTaskPrecheck(task); err != nil {
		return err
	}
	task.Key = key

	// 如果是cron任务，计算第一次执行时间
	if task.CronExpr != "" {
		parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		schedule, err := parser.Parse(task.CronExpr)
		if err != nil {
			log.Printf("invalid cron expression: %v", err)
			return fmt.Errorf("invalid cron expression: %v", err)
		}
		task.ExecuteAt = schedule.Next(time.Now())
	}

	taskBody, err := json.Marshal(task)
	if err != nil {
		return err
	}

	// 使用ExecuteAt作为score
	score := task.ExecuteAt.UnixNano()
	_, err = r.redisClient.Eval(ctx, LuaAddTasks, []string{r.schedulerKey}, score, string(taskBody)).Result()
	if err != nil {
		return err
	}
	return nil
}

// addTaskPrecheck 预检查定时任务参数
func (r *RTimeWheel) addTaskPrecheck(task *RTaskElement) error {
	if task.ExecuteAt.IsZero() && task.CronExpr == "" {
		return fmt.Errorf("either execute_at or cron_expr must be set")
	}
	return nil
}

// LuaAddTasks 添加任务的Lua脚本
const LuaAddTasks = `
local key = KEYS[1]
local score = tonumber(ARGV[1])
local member = ARGV[2]
redis.call("ZADD", key, score, member)
return 1
`

// LuaGetAndRemoveTasks 获取并移除到期任务的Lua脚本
const LuaGetAndRemoveTasks = `
local key = KEYS[1]
local max_score = ARGV[1]
local batch_size = tonumber(ARGV[2])

-- 获取到期的任务
local tasks = redis.call('ZRANGEBYSCORE', key, '-inf', max_score, 'LIMIT', 0, batch_size, 'WITHSCORES')
if #tasks == 0 then
    return {}
end

-- 移除这些任务
for i = 1, #tasks, 2 do
    redis.call('ZREM', key, tasks[i])
end

return tasks
`

// Start 启动调度器
func (r *RTimeWheel) Start() {
	// 不再需要启动cronScheduler
}

// Stop 停止调度器
func (r *RTimeWheel) Stop() {
	r.stopc <- struct{}{}
}
