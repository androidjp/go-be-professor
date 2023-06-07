package main

import (
	"context"
	"fmt"
	"github.com/carlmjohnson/requests"
)

func main() {
	var s string

	// GET请求
	err := requests.
		URL("https://www.baidu.com").
		ToString(&s).
		Fetch(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(s)
}
