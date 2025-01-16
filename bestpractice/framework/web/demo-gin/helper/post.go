package helper

import (
	"github.com/gin-gonic/gin"
	"io"
	"reflect"
)

type Empty struct {
}

// POST 对gin的POST路由进行二次封装
func POST(g *gin.Engine, relativePath string, handler any) {
	g.POST(relativePath, func(c *gin.Context) {
		method := reflect.ValueOf(handler)
		// 1）拿到handler函数的第二个入参（第一个入参是 *gin.Context）
		parameter := method.Type().In(1)
		req := reflect.New(parameter.Elem()).Interface()
		err := c.ShouldBindJSON(req)
		if err != io.EOF && err != nil {
			// 2）入参解析不到，直接报错
			Fail(c, err)
			return
		}
		// 3）直接通过反射的方式，调用handler函数
		in := make([]reflect.Value, 0)
		in = append(in, reflect.ValueOf(c))
		in = append(in, reflect.ValueOf(req))
		call := method.Call(in)
		if !call[1].IsNil() {
			// 4）如果存在第二个返回参数，直接拿到转换为error并返回
			callErr := call[1].Interface().(error)
			Fail(c, callErr)
			return
		}
		// 5）否则说明没有error，再尝试解析第一个参数response对象
		if !call[0].IsNil() {
			// 如果存在此resp对象，直接正常返回即可
			Ok(c, call[0].Interface())
		}
	})
}
