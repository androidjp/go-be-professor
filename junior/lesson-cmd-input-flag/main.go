package main

import (
	"flag"
	"fmt"
)

// 执行：go run main.go -name Alice -age 20
func main() {
	var name string
	var age int

	flag.StringVar(&name, "name", "", "请输入姓名")
	flag.IntVar(&age, "age", 0, "请输入年龄")

	flag.Parse()

	fmt.Printf("姓名：%s，年龄：%d\n", name, age)
}
