package util

import "strings"

const (
	base62chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	scale       = 62
	minLen      = 6
)

func EncodeBase62(num int) string {
	if num == 0 {
		return string(base62chars[0])
	}
	encoded := ""
	for num > 0 {
		remainder := num % 62
		encoded = string(base62chars[remainder]) + encoded
		num = num / 62
	}
	return encoded
}

func DecodeBase62(encoded string) int {
	var decoded int = 0
	for _, char := range encoded {
		decoded = decoded*scale + strings.IndexRune(base62chars, char)
	}
	return decoded
}
