package main

import (
	"bytes"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/petermattis/goid"
)

func main() {

	// 方式一（推荐）
	way_1()

	// 方式二
	way_2()
}

func way_2() {
	// 不依赖库，自己拿堆栈信息，效率较低
	// 这种方式的两个弊端：
	// 1. 效率问题：解析堆栈来获取 协程ID，在高并发场景下可能影响性能
	// 2. 不保证稳定性：解析堆栈信息依赖GO内部实现细节，Go版本更新时，堆栈的格式可能发生变化，从而影响协程ID获取。

	log.SetFlags(0) // 禁止时间戳，让输出简单
	log.SetOutput(os.Stdout)

	go func() {
		logWithGoroutineID("Hello from goroutine 1")
	}()
	go func() {
		logWithGoroutineID("Hello from goroutine 2")
	}()
	go func() {
		logWithGoroutineID("Hello from goroutine 3")
	}()

	logWithGoroutineID("Hello from main goroutine")

	time.Sleep(time.Second)
}

func logWithGoroutineID(s string) {
	gid := getGoroutineID()
	log.Printf("goroutine[%d]: %s", gid, s)
}

func getGoroutineID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func way_1() {
	log.SetFlags(0) // 禁止时间戳，让输出简单
	log.SetOutput(os.Stdout)

	// 在不同的goroutine中打印日志
	go func() {
		log.Printf("goroutine[%d]: Hello from goroutine 1", goid.Get())
	}()

	go func() {
		log.Printf("goroutine[%d]: Hello from goroutine 2", goid.Get())
	}()

	go func() {
		log.Printf("goroutine[%d]: Hello from goroutine 3", goid.Get())
	}()

	// 主goroutine打印
	log.Printf("main goroutine[%d]: Hello from main goroutine", goid.Get())

	// 等待所有goroutine结束
	log.Println("Waiting for goroutines to finish...")
	time.Sleep(time.Second)
}
