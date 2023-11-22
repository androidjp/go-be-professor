package main

import (
	"fmt"
	"runtime"
	"time"
)

func Foo1() {
	fmt.Println("(Foo1)打印1")
	defer fmt.Println("(Foo1)打印2")
	runtime.Goexit() // 加入这行，即可退出协程，且退出之前还能执行到defer函数
	time.Sleep(200 * time.Millisecond)
	fmt.Println("(Foo1)打印3")
}

func main() {
	go Foo1()
	fmt.Println("打印4")
	// 主协程只能用os.Exit退出，不能用 runtime.Goexit()，否则报错找不到协程。
	//os.Exit(0)
	time.Sleep(2000 * time.Millisecond)
}
