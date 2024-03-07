package errhandle

import (
	"fmt"
	"net"
	"os"
	"runtime"
)

/**
CatchError
一般提供给业务方法进行defer后调用，用于捕获方法的panic情况
*/
func CatchError() {
	var r any = recover()
	switch r.(type) {
	case runtime.Error:
		fmt.Println("运行是错误：", r)
	case error:
		if opError, ok := r.(*net.OpError); ok {
			if syscallError, ok := opError.Err.(*os.SyscallError); ok {
				fmt.Printf("catch error: %v\n", syscallError)
				return
			}
		}
		fmt.Printf("catch error: %v\n", r)
		buf := make([]byte, 1<<12)
		size := runtime.Stack(buf, true)
		fmt.Printf("catch stacktrace: %s", buf[:size])
	default:
		fmt.Println("catch default.....")
	}
}
