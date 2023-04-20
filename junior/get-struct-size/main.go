package main

import (
	"fmt"
	"os"
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

	// 实际上，如果要读取系统的KEYA变量，Goland的run，是需要去 配置->Go modules->自己新增自定义key-val的，因为Goland默认就是读取他自己的go环境变量列表
	// 如果使用 命令行直接跑 go run main.go 读取的则是我们本机系统的关键变量，此时，使用 ~/.zshrc 中 'export KEYA=123213123123aaaaa'
	keyA := os.Getenv("KEYA")
	fmt.Println(keyA)
}
