package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 设置静态文件服务
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	// 设置服务端口
	port := 9527
	serverAddr := fmt.Sprintf(":%d", port)

	// 启动服务器
	log.Printf("服务器启动在 http://localhost%s", serverAddr)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
}
