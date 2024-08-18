package middleware

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

// 创建限速器，每秒允许 10 个请求，令牌桶大小为 20
var limiter = rate.NewLimiter(1, 2)

func RateLimitWay1(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		// 等待获取令牌
		if !limiter.Allow() {
			return c.JSON(429, nil)
		}
		return next(c)
	}

}
