package scheduler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
)

func TestScheduler_Start_And_Stop(t *testing.T) {
	// 启动一个模拟的 Redis 服务器
	mrs, err := miniredis.Run()
	if err != nil {
		t.Fatalf("无法启动模拟 Redis 服务器: %v", err)
	}
	defer mrs.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: mrs.Addr(),
	})

	httpClient := &http.Client{}

	timeWheel := NewRTimeWheel(redisClient, httpClient, "test_tasks", 2, 10*time.Second)

	// 添加一个定时任务
	task1 := &RTaskElement{
		// 每秒钟执行一次
		CronExpr: "*/1 * * * * *",
	}
	timeWheel.AddTask(context.Background(), "task1", task1)

	timeWheel.Start()

	// 添加一次性任务，立即执行
	task2 := &RTaskElement{
		ExecuteAt: time.Now(),
	}
	timeWheel.AddTask(context.Background(), "task2", task2)

	// 添加一次性任务：3秒后执行
	task3 := &RTaskElement{
		ExecuteAt: time.Now().Add(3 * time.Second),
	}
	timeWheel.AddTask(context.Background(), "task3", task3)

	// 等待5秒，确保task1至少执行了4次，task2和task3也都执行了
	time.Sleep(4500 * time.Millisecond)

	// 获取Redis中的任务数据，验证task1是否被重新添加
	ctx := context.Background()
	tasks, err := redisClient.ZRange(ctx, "test_tasks", 0, -1).Result()
	if err != nil {
		t.Fatalf("获取Redis任务失败: %v", err)
	}

	// 应该至少有一个task1的下一次执行时间在Redis中
	if len(tasks) == 0 {
		t.Error("没有找到task1的下一次执行时间")
	}

	log.Println("5秒后停止定时器")
	timeWheel.Stop()
}

func TestMultiSchedulerConsistency(t *testing.T) {
	// 启动一个模拟的 Redis 服务器
	mrs, err := miniredis.Run()
	if err != nil {
		t.Fatalf("无法启动模拟 Redis 服务器: %v", err)
	}
	defer mrs.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: mrs.Addr(),
	})

	// 创建5个scheduler实例
	const numSchedulers = 5
	schedulers := make([]*RTimeWheel, numSchedulers)
	for i := 0; i < numSchedulers; i++ {
		httpClient := &http.Client{}
		schedulers[i] = NewRTimeWheel(redisClient, httpClient, "test_tasks", 2, 10*time.Second)
	}

	// 启动所有scheduler
	for _, s := range schedulers {
		s.Start()
	}
	defer func() {
		for _, s := range schedulers {
			s.Stop()
		}
	}()

	// 添加多个一次性任务
	const numTasks = 10
	for i := 0; i < numTasks; i++ {
		task := &RTaskElement{
			ExecuteAt: time.Now().Add(time.Duration(i) * time.Second),
		}
		err := schedulers[0].AddTask(context.Background(), fmt.Sprintf("task%d", i), task)
		if err != nil {
			t.Fatalf("添加任务失败: %v", err)
		}
	}

	// 等待足够的时间让所有任务执行完成
	time.Sleep(15 * time.Second)

	// 验证Redis中是否所有任务都已被执行（已被移除）
	ctx := context.Background()
	tasks, err := redisClient.ZRange(ctx, "test_tasks", 0, -1).Result()
	if err != nil {
		t.Fatalf("获取Redis任务失败: %v", err)
	}

	// 由于是一次性任务，执行完后应该被从Redis中删除
	if len(tasks) > 0 {
		t.Errorf("预期所有任务都已执行完成，但仍有%d个任务在Redis中", len(tasks))
	}
}
