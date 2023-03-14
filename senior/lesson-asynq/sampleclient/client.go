package sampleclient

import (
	"fmt"
	"github.com/hibiken/asynq"
	"lesson-asynq/test_delivery"
	"log"
	"time"
)

func EmailDeliveryTaskAdd(i int) {
	// 连接本地docker的redis客户端
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     "127.0.0.1:16379",
		Password: "",
		DB:       2,
	})
	defer client.Close()

	// 初使货需要传递的数据
	task, err := test_delivery.NewEmailDeliveryTask(42, fmt.Sprintf("some:template:id:%d", i), `{"name":"lisi"}`)
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}
	// 任务入队
	//info, err := client.Enqueue(task)

	//info, err := client.Enqueue(task, time.Now())
	// 延迟执行
	info, err := client.Enqueue(task, asynq.ProcessIn(3*time.Second))
	// MaxRetry 重度次数 Timeout超时时间
	//info, err = client.Enqueue(task, asynq.MaxRetry(10), asynq.Timeout(3*time.Second))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)
}
