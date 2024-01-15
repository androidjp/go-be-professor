package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	numCodeLen := 7 // 6位数字

	getShuffleNumCodeList(numCodeLen)

}

func getShuffleNumCodeList(numCodeLen int) []string {
	// 10的N次方
	numMax := 1
	for i := 0; i < numCodeLen; i++ {
		numMax *= 10
	}
	format := fmt.Sprintf("%s%dd", "%0", numCodeLen)

	// 100 个字符串的数组
	stringsSlice := make([]string, numMax)
	for i := 0; i < numMax; i++ {
		stringsSlice[i] = fmt.Sprintf(format, i)
	}

	fmt.Println("Before shuffle:", stringsSlice[0:3], ",...,", stringsSlice[numMax-2:numMax])

	// 为了获得更好的随机性，我们需要设置随机数种子
	rand.NewSource(time.Now().UnixNano())
	// 使用math/rand.Shuffle函数打乱字符串切片
	n := len(stringsSlice)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		stringsSlice[i], stringsSlice[j] = stringsSlice[j], stringsSlice[i]
	}
	// 打乱后的字符串切片现在存储在stringsSlice中
	fmt.Println("After shuffle:", stringsSlice[0:3], ",...,", stringsSlice[numMax-2:numMax])
	return stringsSlice
}
