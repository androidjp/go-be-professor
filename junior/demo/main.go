package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	StrategyLatestChatList = iota // 最新聊天列表
	StrategyFunc2
	StrategyFunc3
	StrategyFunc4
)

func main() {
	var cnt0, cnt1, cnt2, cnt3 int

	var res int
	for i := 0; i < 10000; i++ {
		res = RatioAction()
		if res == StrategyLatestChatList {
			cnt0++
		}
		if res == StrategyFunc2 {
			cnt1++
		}
		if res == StrategyFunc3 {
			cnt2++
		}
		if res == StrategyFunc4 {
			cnt3++
		}

	}

	fmt.Println(cnt0, cnt1, cnt2, cnt3)
}

func RatioAction() int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randNum := rand.Intn(100)

	// 70% 概率走 策略1
	if randNum < 70 {
		return StrategyLatestChatList
	}
	if randNum < 90 {
		return StrategyFunc2
	}
	if randNum < 95 {
		return StrategyFunc3
	}
	return StrategyFunc4
}

type Person struct {
	Name string
	Age  int64
}
