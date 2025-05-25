package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 中间件用于记录API请求日志
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latency := endTime.Sub(startTime)

		// 请求方式
		method := c.Request.Method

		// 请求路由
		path := c.Request.URL.Path

		// 状态码
		statusCode := c.Writer.Status()

		// 客户端IP
		clientIP := c.ClientIP()

		// 日志格式
		slog.Info("API请求日志", "status", statusCode, "latency", latency, "clientIP", clientIP, "method", method, "path", path)
	}
}
