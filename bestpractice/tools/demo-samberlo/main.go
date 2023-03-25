package main

import (
	"fmt"
	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("--------------数组直接去重--------------")
	names := lo.Uniq[string]([]string{"Samuel", "John", "Samuel"})
	fmt.Println(len(names)) // 2
	fmt.Println(names[0])   // 2
	fmt.Println(names[1])   // 2

	fmt.Println("--------------filter 根据条件过滤--------------")
	even := lo.Filter[int]([]int{1, 2, 3, 4}, func(x int, index int) bool {
		return x%2 == 0
	})
	// []int{2,4}
	fmt.Println(len(even))
	fmt.Println(even[0], even[1])

	fmt.Println("--------------map函数 直接进行元素的类型转换 int64 转 string--------------")
	strs := lo.Map[int64, string]([]int64{1, 2, 3, 4}, func(x int64, index int) string {
		return strconv.FormatInt(x, 10)
	})
	fmt.Println(strs[0], strs[1], strs[2], strs[3])

	fmt.Println("--------------map函数 直接进行元素的类型转换 int64 转 string 【并发：内部函数执行是 goroutine跑的】--------------")
	stringsP := lop.Map[int64, string]([]int64{1, 2, 3, 4}, func(x int64, _ int) string {
		return strconv.FormatInt(x, 10)
	})
	fmt.Println(stringsP[0], stringsP[1], stringsP[2], stringsP[3])

	fmt.Println("--------------filterMap函数 同时 做 过滤 + 值转换 的 操作--------------")
	fAndM := lo.FilterMap[string, int64]([]string{"cpu", "gpu", "mouse", "keyboard"}, func(x string, _ int) (int64, bool) {
		if strings.HasSuffix(x, "pu") {
			return 111, true
		}
		return 000, false
	})
	// []string{"xpu", "xpu"}
	fmt.Println(fAndM[0], fAndM[1])

	fmt.Println("--------------flatMap函数 将数组中的元素展开成一个新的数组--------------")
	flat := lo.FlatMap[string, string]([]string{"hello world", "how are you"}, func(x string, _ int) []string {
		return strings.Split(x, " ")
	})
	// []string{"hello", "world", "how", "are", "you"}
	fmt.Println(flat[0], flat[1], flat[2], flat[3], flat[4])

	fmt.Println("--------------reduce函数 对数组中的元素进行累加操作--------------")
	sum := lo.Reduce[int]([]int{1, 2, 3, 4}, func(acc, x int, _ int) int {
		return acc + x
	}, 0)
	fmt.Println(sum) // 10

	sumRight := lo.ReduceRight[int]([]int{1, 2, 3, 4}, func(acc, x int, _ int) int {
		return acc - x
	}, 0)
	fmt.Println(sumRight) // -2

	fmt.Println("--------------forEach函数 对数组中的元素进行遍历操作--------------")
	lo.ForEach[string]([]string{"hello", "world"}, func(x string, _ int) {
		fmt.Println(x)
	})
	fmt.Println("--------------parallelForEach函数 对数组中的元素进行并发遍历操作--------------")
	lop.ForEach[string]([]string{"hello", "world"}, func(x string, _ int) {
		fmt.Println(x)
	})
	fmt.Println("--------------Times函数 重复执行某个操作（可用lop异步）--------------")
	intArr := lo.Times(3, func(i int) int {
		fmt.Println("hello", i)
		return i * 2
	})
	fmt.Println(len(intArr))
	fmt.Println(intArr[0])
	fmt.Println(intArr[1])
	fmt.Println(intArr[2])
	fmt.Println("--------------Uniq函数 字符串直接去重--------------")
	str := lo.Uniq[string]([]string{"Samuel", "John", "Samuel"})
	fmt.Println(len(str)) // 2
	fmt.Println(str[0])   // Samuel
	fmt.Println(str[1])   // John
	fmt.Println("--------------UniqBy函数 根据指定条件去重--------------")
	numbers := lo.UniqBy[int, string]([]int{1, 2, 3, 4}, func(x int) string {
		return strconv.Itoa(x % 2)
	})
	fmt.Println(len(numbers)) // 2
	fmt.Println(numbers[0])   // 1
	fmt.Println(numbers[1])   // 2
	fmt.Println("--------------groupBy函数 根据指定条件分组--------------")
	groups := lo.GroupBy[string, int]([]string{"apple", "banana", "cherry", "date"}, func(x string) int {
		return len(x)
	})
	fmt.Println(groups[5][0]) // apple
	fmt.Println(groups[6][0]) // banana, cherry
	fmt.Println(groups[4][0]) // date
	fmt.Println("--------------Chunk函数 将数组按照指定大小分块--------------")
	chunks := lo.Chunk[int]([]int{1, 2, 3, 4, 5, 6}, 2)
	fmt.Println(len(chunks)) // 3
	fmt.Println(chunks[0])   // []int{1, 2}
	fmt.Println(chunks[1])   // []int{3, 4}
	fmt.Println(chunks[2])   // []int{5, 6}
	fmt.Println("--------------PartitionBy函数 根据指定条件分割数组--------------")
	partitioned := lo.PartitionBy[int]([]int{1, 2, 3, 4, 5, 6}, func(x int) bool {
		return x%2 == 0
	})
	fmt.Println(len(partitioned)) // 2
	fmt.Println(partitioned[0])   // []int{1, 3, 5}
	fmt.Println(partitioned[1])   // []int{2, 4, 6}
	fmt.Println("--------------PartitionByP函数 根据指定条件分割数组，最终结果返回的是按入参顺序的【并发：内部函数执行是 goroutine跑的】--------------")
	partitionedP := lop.PartitionBy[int]([]int{1, 2, 3, 4, 5, 6}, func(x int) bool {
		return x%2 == 0
	})
	fmt.Println(len(partitionedP)) // 2
	fmt.Println(partitionedP[0])   // []int{1, 3, 5}
	fmt.Println(partitionedP[1])   // []int{2, 4, 6}

}
