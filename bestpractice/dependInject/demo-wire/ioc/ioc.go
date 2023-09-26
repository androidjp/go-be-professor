package ioc

import (
	"demo-wire/handler"
	"github.com/gin-gonic/gin"
)

func NewGinEngineAndRegisterRoute(postHandler *handler.PostHandler) *gin.Engine {
	engine := gin.Default()
	postHandler.RegisterRoutes(engine)
	return engine
}
