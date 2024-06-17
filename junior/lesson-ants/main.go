package main

import (
	"errors"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"time"
)

func main() {
	// 初始化
	runTask := func(i interface{}) {
		if i == 20 {
			panic(errors.New("dfdf"))
		}
		// 模拟任务处理逻辑
		fmt.Printf("Running task: %d\n", i.(int))
		time.Sleep(time.Second)
	}

	// 创建一个具有10个并发量的池(默认是阻塞模式，即所有任务完成后才退出)
	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		runTask(i)
	})
	defer p.Release()

	// 提交任务
	for i := 0; i < 30; i++ {
		_ = p.Invoke(i)
	}

	// 等待所有任务完成
	p.Waiting()
	fmt.Println("All tasks are done.")
}
