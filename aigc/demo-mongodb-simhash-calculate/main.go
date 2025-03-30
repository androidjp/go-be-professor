package main

import "fmt"

// 实现一个存入 mongodb 的 document 对象，document 对象有个byte 数组字段，用于存储 simhash 值
type SimhashDocument struct {
	// 用于存储 simhash 值
	Simhash []byte `bson:"simhash"`
}

func main() {

	// 模拟计算两个 Document 的 Simhash 的 海明距离
	// 1. 生成两个 Document 的 Simhash
	// 2. 计算两个 Document 的 Simhash 的 海明距离

	a := &SimhashDocument{
		Simhash: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
	}

	b := &SimhashDocument{
		Simhash: []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
	}

	// 海明距离小于某个值，表示两个文档相似度很高
	d := calculateSimhashDistance(a.Simhash, b.Simhash)
	if d < 3 {
		fmt.Println("Similar")
	} else {
		fmt.Println("Not similar")
	}
}

func calculateSimhashDistance(b1, b2 []byte) int {
	if len(b1) != len(b2) {
		return -1
	}

	distance := 0
	for i := 0; i < len(b1); i++ {
		distance += hammingDistance(b1[i], b2[i])
	}

	return distance
}

func hammingDistance(b1, b2 byte) int {
	x := b1 ^ b2
	// 得到这个byte 中有多少个 1
	count := 0
	for x != 0 {
		x = x & (x - 1)
		count++
	}
	return count
}
