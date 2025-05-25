package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware 处理跨域请求的中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许所有来源的请求
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 允许的HTTP方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		// 允许的HTTP头
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		// 允许暴露的头
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		// 允许携带凭证
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
