package main

import (
	"fmt"
	"time"
)

func main() {
	//------------------------------------------------
	// time.After 某些情况会发生内存泄露，替换为：Timer
	//------------------------------------------------
	fmt.Println("A")
	after := time.After(2 * time.Second)

	select {
	// 这里会等待2秒，然后输出B
	case <-after:
		fmt.Println("B")
	}

	fmt.Println("C")
	t := time.NewTimer(2 * time.Second)
	select {
	// 这里会等待2秒，然后输出B
	case <-t.C:
		fmt.Println("D")
	}
}
