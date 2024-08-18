package main

import (
	"context"
	"demo-echo/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()

	// 默认日志级别是Error
	e.Logger.SetLevel(log.DEBUG)

	// 限流能力
	// e.Use(middleware.RateLimitWay1)
	e.Use(middleware.CustomRateLimit)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"msg": "Hello, World!",
		})
	})

	// 最简单的启动方式
	// e.Logger.Fatal(e.Start(":1323"))

	//-----------------------------------------------------
	// 以下是优雅关闭
	//-----------------------------------------------------

	// 定义一个关闭的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 启动服务器
	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatalf("Shutting down the server: %T, %v", err, err.Error())
		}
	}()

	// 监听中断信号（Ctrl+C 或 kill 命令）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 可以添加一些逻辑来等待正在处理的请求完成
	// for {
	// 	if len(e.Router().Routes()) == 0 {
	// 		break
	// 	}
	// 	time.Sleep(100 * time.Millisecond)
	// }
	e.Logger.Info("开始关闭服务")

	// 接收到信号后开始关闭服务器
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Error(err)
	}
	e.Logger.Info("关闭服务完成")
}
