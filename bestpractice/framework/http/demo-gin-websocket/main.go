package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type MyLogic struct {
	Manager *WebSocketManager
}

func (lgc *MyLogic) WSHandle(c *gin.Context) {
	// 创建ws连接客户端
	client, err := lgc.Manager.upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	// 连接成功后存储当前客户端
	lgc.Manager.SaveClient(client)
	defer func() {
		client.Close()
		lgc.Manager.DelClient(client)
	}()
	for {
		// 监听接受信息
		mt, message, err := client.ReadMessage()
		if err == nil {
			lgc.Manager.SendMsg(mt, message)
		} else {
			break
		}
	}
}

// 服务端实现
func main() {
	lgc := &MyLogic{Manager: &WebSocketManager{
		upGrader: &websocket.Upgrader{
			HandshakeTimeout: 0,
			ReadBufferSize:   0,
			WriteBufferSize:  0,
			WriteBufferPool:  nil,
			Subprotocols:     nil,
			Error:            nil,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			EnableCompression: false,
		},
		wsClients:        sync.Map{},
		aliveClientNum:   0,
		historyClientNum: 0,
	}}

	router := gin.Default()

	// 注册html页面文件
	router.LoadHTMLFiles("cli.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "cli.html", nil)
	})

	// 注册websocket路由
	router.GET("/ws", lgc.WSHandle)

	// 启动应用
	err := router.Run(":9999")
	if err != nil {
		panic(err)
	}
}
