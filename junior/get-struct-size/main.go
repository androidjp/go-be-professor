package main

import (
	"fmt"
	"unsafe"
)

type MyStruct struct {
	Tag       [8]byte
	Version   [4]byte
	Author    [32]byte
	CopyRight [16]byte // 到这里应该长度是60
	Name      string   // 20
	Age       int64    // 8
}

func (m *MyStruct) GetLength() int {
	return int(unsafe.Sizeof(*m))
}

func main() {
	m1 := MyStruct{}
	fmt.Println(m1.GetLength()) // 88
	m2 := MyStruct{Name: "abc"}
	fmt.Println(m2.GetLength()) // 88
	m3 := MyStruct{Name: "小明"}
	fmt.Println(m3.GetLength()) // 88
}
