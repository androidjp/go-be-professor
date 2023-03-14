package main

import (
	"github.com/hibiken/asynq"
	"lesson-asynq/sampleclient"
	"lesson-asynq/test_delivery"
	"log"
	"time"
)

func main() {
	// 创建服务端，用于接收新的任务 进来
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     "127,.0.0.1:16379",
			Password: "",
			DB:       2,
		},
		asynq.Config{
			// 每个进程并发执行的worker数
			Concurrency: 5,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)
	mux := asynq.NewServeMux()
	mux.HandleFunc(test_delivery.TypeEmailDelivery, test_delivery.HandleEmailDeliveryTask)

	go func() {
		time.Sleep(time.Second)
		// 延迟1秒后，生成几个异步任务出来
		for i := 0; i < 3; i++ {
			sampleclient.EmailDeliveryTaskAdd(i)
			time.Sleep(time.Second * 2)
		}
	}()

	if err := srv.Run(mux); err != nil {
		log.Fatalf("conuld not run server: %v", err)
	}
}
