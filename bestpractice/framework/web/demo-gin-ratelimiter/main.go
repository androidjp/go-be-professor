package main

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
)

// RateLimitMiddleware 封装为 gin中间件的限流函数
func RateLimitMiddleware(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			// 请求获取令牌数为1，如果返回被消耗的令牌数小于1，说明无法再拿到令牌了。
			// 制定响应内容，并执行中断操作
			c.String(http.StatusForbidden, "rate limit...")
			c.Abort()
			return
		}
		// 正常情况，走下一步
		c.Next()
	}
}

func main() {
	r := gin.Default()
	gin.ForceConsoleColor()
	// 使用 开源库的 限流器，容量是5，每秒钟产生1个令牌
	r.Use(RateLimitMiddleware(time.Second, 5, 1))
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "golang ~")
	})
	r.Run(":8080")
}
