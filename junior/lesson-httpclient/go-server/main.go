package main

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"go-server/limit"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化http客户端
	router := gin.Default()

	// 设置限制最大连接为5(这个其实只是阻塞，但是不合理，最主要还是得用limit)
	router.Use(limit.MaxAllowed(5))

	// 创建一个路由组
	apiGroup := router.Group("/api")

	// 自定义路由
	apiGroup.GET("/hello/:sec", func(c *gin.Context) {
		// 1. 校验
		sec, ok := c.Params.Get("sec")
		if !ok {
			c.JSON(400, gin.H{
				"message": "sec is required",
			})
			return
		}
		secInt, err := strconv.Atoi(sec)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "sec is invalid",
			})
			return
		}

		// 2. 核心逻辑
		time.Sleep(time.Duration(secInt) * time.Second)
		c.JSON(200, gin.H{
			"cost":    secInt,
			"message": "Hello World",
		})
	})

	// // 启动
	// router.Run(":8080")

	// 那么，如果在gin层面，如何限制最大http连接呢？？
	// 答案：自定义server，放入router
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		MaxHeaderBytes:    1 << 20, // 1M
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
		// TLSConfig: &tls.Config{
		// 	// 设置证书文件和密钥文件路径
		// 	Certificates: []tls.Certificate{cert},
		// 	// 启用 TLS 版本 1.2 和 1.3
		// 	MinVersion: tls.VersionTLS12,
		// 	MaxVersion: tls.VersionTLS13,
		// },
		// MaxHeaderBytes: http.DefaultMaxHeaderBytes, // 1MB
		// 设置 ConnState 回调函数
		ConnState: func(conn net.Conn, state http.ConnState) {
			fmt.Printf("Connection state changed: %s -> %s\n", conn.RemoteAddr(), state)
			// 根据连接状态进行特定的处理
			switch state {
			case http.StateNew:
				// 新连接建立时的处理
				fmt.Println("New connection established")
			case http.StateClosed:
				// 连接关闭时的处理
				fmt.Println("Connection closed")
			}
		},
		ErrorLog:    nil, // 默认为空，使用系统默认日志
		BaseContext: nil, // 默认为空，使用 context.Background
		ConnContext: nil, // 默认为空即不启用。作为baseContext的子类，可用于包装一些k-v对
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
