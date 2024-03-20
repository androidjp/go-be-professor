package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	chineseChar1 := "你好"
	chineseChar2 := "你好！"
	englishChar := "Hello"

	fmt.Println(len(chineseChar1)) // 6
	fmt.Println(len(chineseChar2)) // 9
	fmt.Println(len(englishChar))  // 5

	fmt.Println(utf8.RuneCountInString(chineseChar1)) // 2
	fmt.Println(utf8.RuneCountInString(chineseChar2)) // 3
	fmt.Println(utf8.RuneCountInString(englishChar))  // 5
}
