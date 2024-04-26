package main

import (
	"fmt"
	"time"
)

func main() {
	defer Since(time.Now())
	defer TrackTime("大名")()

	time.Sleep(2 * time.Second)
}

func Since(startT time.Time) {
	elapsed := time.Since(startT)
	fmt.Printf("开始 到 结束 一共：%v\n", elapsed)
}

func TrackTime(key string) func() {
	startT := time.Now()
	return func() {
		endT := time.Now()
		dur := endT.Sub(startT)
		fmt.Printf("%s: 开始（%v） 到 结束（%v）一共：%v\n", key, startT, endT, dur)
	}
}
