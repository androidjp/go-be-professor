package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
)

// 自己实现的一套令牌桶能力
type TokenBucket struct {
	capacity  int       // 池总大小
	tokens    int       // 实时令牌数
	rate      int       // 速率
	lastToken time.Time // 上次获取令牌时间
}

func NewTokenBucket(capacity, rate int) *TokenBucket {
	tb := &TokenBucket{
		capacity:  capacity,
		tokens:    capacity,
		rate:      rate,
		lastToken: time.Now(),
	}
	return tb
}

func (tb *TokenBucket) TakeToken() bool {
	now := time.Now()
	tokensToAdd := int(now.Sub(tb.lastToken).Seconds() * float64(tb.rate))
	tb.tokens = tb.tokens + tokensToAdd
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
	tb.lastToken = now
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

// 创建令牌桶，每秒 10 个令牌，桶容量为 20
var tokenBucket = NewTokenBucket(2000, 1000)

func CustomRateLimit(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 尝试获取令牌
		if !tokenBucket.TakeToken() {
			return c.JSON(429, nil)
		}
		return next(c)
	}
}
