package distributed

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/demdxx/gocast"
	"lesson-timewheel/libs/redis"
	"lesson-timewheel/libs/xhttp"
	"lesson-timewheel/util"
	"net/http"
	"strings"
	"sync"
	"time"
)

type RedisTimeWheel struct {
	// 单例工具，保证时间轮停止操作只能执行一次
	sync.Once
	// 类比：驱动指针运转的齿轮
	// 时间轮定时器
	ticker *time.Ticker
	// 停止时间轮的 channel
	stopc chan struct{}
	// redis客户端
	redisCli *redis.Client
	// http客户端，在执行定时任务时需要使用到
	httpCli *xhttp.Client
}

type RTaskElement struct {
	// 定时任务全局唯一key
	Key      string       `json:"key"`
	Callback CallbackInfo `json:"callback"`
}

type CallbackInfo struct {
	// 定时任务执行时，回调的http url
	URL string `json:"url"`
	// 回调使用的http方法
	Method string `json:"method"`
	// 回调的req请求参数
	Req interface{} `json:"req"`
	// 回调使用的http请求头
	Header map[string]string `json:"header"`
}

func NewRedisTimeWheel(redisClient *redis.Client, httpClient *xhttp.Client) *RedisTimeWheel {
	r := RedisTimeWheel{
		ticker:   time.NewTicker(time.Second),
		redisCli: redisClient,
		httpCli:  httpClient,
		stopc:    make(chan struct{}),
	}

	go r.run()
	return &r
}

func (r *RedisTimeWheel) Stop() {
	r.Do(func() {
		close(r.stopc)
		r.ticker.Stop()
	})
}

func (r *RedisTimeWheel) AddTask(ctx context.Context, key string, task *RTaskElement, executeAt time.Time) error {
	if err := r.addTaskPrecheck(task); err != nil {
		return err
	}

	task.Key = key
	taskBody, _ := json.Marshal(task)
	_, err := r.redisCli.Eval(ctx, LuaAddTasks, 2, []interface{}{
		// 分钟级 zset 时间片
		r.getMinuteSlice(executeAt),
		// 标识任务删除的集合
		r.getDeleteSetKey(executeAt),
		// 以执行时刻的秒级时间戳作为 zset 中的 score
		executeAt.Unix(),
		// 任务明细
		string(taskBody),
		// 任务 key，用于存放在删除集合中
		key,
	})
	return err
}

func (r *RedisTimeWheel) RemoveTask(ctx context.Context, key string, executeAt time.Time) error {
	// 标识任务已被删除
	_, err := r.redisCli.Eval(ctx, LuaDeleteTask, 1, []interface{}{
		r.getDeleteSetKey(executeAt),
		key,
	})
	return err
}

func (r *RedisTimeWheel) run() {
	for {
		select {
		case <-r.stopc:
			return
		case <-r.ticker.C:
			// 每次 tick 获取任务
			go r.executeTasks()
		}
	}
}

func (r *RedisTimeWheel) executeTasks() {
	defer func() {
		if err := recover(); err != nil {
			// log
		}
	}()

	// 并发控制，30 s
	tctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	tasks, err := r.getExecutableTasks(tctx)
	if err != nil {
		// log
		return
	}

	// 并发执行任务
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		// shadow
		task := task
		go func() {
			defer func() {
				if err := recover(); err != nil {
				}
				wg.Done()
			}()
			if err := r.executeTask(tctx, task); err != nil {
				// log
			}
		}()
	}
	wg.Wait()
}

func (r *RedisTimeWheel) executeTask(ctx context.Context, task *RTaskElement) error {
	return r.httpCli.JSONDo(ctx, task.Callback.Method, task.Callback.URL, task.Callback.Header, task.Callback.Req, nil)
}

func (r *RedisTimeWheel) addTaskPrecheck(task *RTaskElement) error {
	if task.Callback.Method != http.MethodGet && task.Callback.Method != http.MethodPost {
		return fmt.Errorf("invalid method: %s", task.Callback.Method)
	}
	if !strings.HasPrefix(task.Callback.URL, "http://") && !strings.HasPrefix(task.Callback.URL, "https://") {
		return fmt.Errorf("invalid url: %s", task.Callback.URL)
	}
	return nil
}

func (r *RedisTimeWheel) getExecutableTasks(ctx context.Context) ([]*RTaskElement, error) {
	now := time.Now()
	minuteSlice := r.getMinuteSlice(now)
	deleteSetKey := r.getDeleteSetKey(now)
	nowSecond := util.GetTimeSecond(now)
	score1 := nowSecond.Unix()
	score2 := nowSecond.Add(time.Second).Unix()
	rawReply, err := r.redisCli.Eval(ctx, LuaZrangeTasks, 2, []interface{}{
		minuteSlice, deleteSetKey, score1, score2,
	})
	if err != nil {
		return nil, err
	}

	replies := gocast.ToInterfaceSlice(rawReply)
	if len(replies) == 0 {
		return nil, fmt.Errorf("invalid replies: %v", replies)
	}

	deleteds := gocast.ToStringSlice(replies[0])
	deletedSet := make(map[string]struct{}, len(deleteds))
	for _, deleted := range deleteds {
		deletedSet[deleted] = struct{}{}
	}

	tasks := make([]*RTaskElement, 0, len(replies)-1)
	for i := 1; i < len(replies); i++ {
		var task RTaskElement
		if err := json.Unmarshal([]byte(gocast.ToString(replies[i])), &task); err != nil {
			// log
			continue
		}

		if _, ok := deletedSet[task.Key]; ok {
			continue
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (r *RedisTimeWheel) getMinuteSlice(executeAt time.Time) string {
	return fmt.Sprintf("xiaoxu_timewheel_task_{%s}", util.GetTimeMinuteStr(executeAt))
}

func (r *RedisTimeWheel) getDeleteSetKey(executeAt time.Time) string {
	return fmt.Sprintf("xiaoxu_timewheel_delset_{%s}", util.GetTimeMinuteStr(executeAt))
}
