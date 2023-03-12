package main

// 执行： go build -gcflags '-m -m -l'
// 逃逸，因为 fmt.Println 入参的类型是interface{}，编译期无法确定其具体参数类型，所以内存分配到堆中。
//func main() {
//	a := 666
//	fmt.Println(a)
//}

//// 逃逸原因： 由于test函数的变量a在test()执行完毕后栈帧就被销毁了，但是 对象引用 被返回到函数之外（相当于被外部引用），为了避免说 对象引用被拿到使用发现内存已经被释放回收而造成非法内存，于是变量的内存分配放到了堆上
//func main() {
//	_ = test()
//}
//
//func test() *int {
//	a := 10
//	return &a
//}

// 变量内存占用较大 (当长度是100时，整体正常)
//func main() {
//	test()
//}
//
//func test() {
//	a := make([]int, 10000, 10000)
//	for i := 0; i < 10000; i++ {
//		a[i] = i
//	}
//}

//// 内存大小不确定时
//func main() {
//	test()
//}
//
//func test() {
//	l := 1
//	// 由于编译期间，只能知道是变量l，不知道变量l 的具体值
//	a := make([]int, l, l)
//	for i := 0; i < l; i++ {
//		a[i] = i
//	}
//}

// 内存大小不确定时
func main() {
	test()
}

func test() {
	l := 1
	// 由于编译期间，只能知道是变量l，不知道变量l 的具体值
	a := make([]int, l, l)
	for i := 0; i < l; i++ {
		a[i] = i
	}
}
