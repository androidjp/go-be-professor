package main

import (
	"fmt"
	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"
	"strconv"
	"strings"
)

func main() {
	//------------------------------------------------------
	// 数组直接去重
	//------------------------------------------------------
	names := lo.Uniq[string]([]string{"Samuel", "John", "Samuel"})

	fmt.Println(len(names)) // 2
	fmt.Println(names[0])   // 2
	fmt.Println(names[1])   // 2

	//------------------------------------------------------
	// filter 根据条件过滤
	//------------------------------------------------------
	even := lo.Filter[int]([]int{1, 2, 3, 4}, func(x int, index int) bool {
		return x%2 == 0
	})
	// []int{2,4}
	fmt.Println(len(even))
	fmt.Println(even[0], even[1])

	//------------------------------------------------------
	// map函数 直接进行元素的类型转换 int64 转 string
	//------------------------------------------------------
	strs := lo.Map[int64, string]([]int64{1, 2, 3, 4}, func(x int64, index int) string {
		return strconv.FormatInt(x, 10)
	})
	fmt.Println(strs[0], strs[1], strs[2], strs[3])

	//------------------------------------------------------
	// map函数 直接进行元素的类型转换 int64 转 string 【并发：内部函数执行是 goroutine跑的】
	//------------------------------------------------------
	stringsP := lop.Map[int64, string]([]int64{1, 2, 3, 4}, func(x int64, _ int) string {
		return strconv.FormatInt(x, 10)
	})
	fmt.Println(stringsP[0], stringsP[1], stringsP[2], stringsP[3])

	//------------------------------------------------------
	// filterMap函数 同时 做 过滤 + 值转换 的 操作
	//------------------------------------------------------
	fAndM := lo.FilterMap[string, int64]([]string{"cpu", "gpu", "mouse", "keyboard"}, func(x string, _ int) (int64, bool) {
		if strings.HasSuffix(x, "pu") {
			return 111, true
		}
		return 000, false
	})
	// []string{"xpu", "xpu"}
	fmt.Println(fAndM[0], fAndM[1])

}
