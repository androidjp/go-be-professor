package main

import (
	"demo-ratelimit/swlimit"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	gin.ForceConsoleColor()
	// 使用 开源库的 限流器，容量是5，每秒钟产生1个令牌
	r.Use(swlimit.LimitFunc)
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "golang ~")
	})
	_ = r.Run(":8080")
}
