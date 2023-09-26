//go:build wireinject

package wire

import (
	"demo-wire/handler"
	"demo-wire/ioc"
	"demo-wire/service"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

/*
InitializeApp 函数就是一个注入器，
函数内部通过 wire.Build 函数连接所有的提供者，然后返回 &gin.Engine{}，
该返回值实际上并没有使用到，只是为了满足编译器的要求，避免报错而已，真正的返回值来自 ioc.NewGinEngineAndRegisterRoute
*/
func InitializeApp() *gin.Engine {
	wire.Build(
		handler.NewPostHandler,
		service.NewPostService,
		ioc.NewGinEngineAndRegisterRoute,
	)
	return &gin.Engine{}
}

// v2优化版本，使用提供者分组wireSet
func InitializeAppV2() *gin.Engine {
	wire.Build(
		handler.PostSet,
		ioc.NewGinEngineAndRegisterRoute,
	)
	return &gin.Engine{}
}

// v3优化版本，建立接口和实现类的关系
func InitializeAppV3() *gin.Engine {
	wire.Build(
		handler.NewPostHandler,
		service.NewPostServiceV3,
		ioc.NewGinEngineAndRegisterRoute,
		wire.Bind(new(service.IPostService), new(*service.PostService)),
	)
	return &gin.Engine{}
}

// 值表达式 + 接口类型使用
func InjectPostService() service.IPostService {
	wire.Build(wire.InterfaceValue(new(service.IPostService), &service.PostService{}))
	return nil
}

// v4, 备注注入器语法 , 省去了最后的return &gin.Engine{}
func InitializeGin() *gin.Engine {
	panic(wire.Build(
		handler.NewPostHandler,
		service.NewPostServiceV3,
		ioc.NewGinEngineAndRegisterRoute,
		wire.Bind(new(service.IPostService), new(*service.PostService)),
	))
}
