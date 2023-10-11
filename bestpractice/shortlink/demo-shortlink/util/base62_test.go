package util

import (
	"fmt"
	"testing"
)

func TestEncodeBase62(t *testing.T) {
	res := EncodeBase62(1) // 1
	fmt.Println(res)

	res = EncodeBase62(8888888888) //9hYs3c
	fmt.Println(res)
}

func TestDecodeBase62(t *testing.T) {
	res := DecodeBase62("1") // 1
	fmt.Println(res)

	res = DecodeBase62("9hYs3c") //8888888888
	fmt.Println(res)
}
