package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Client struct {
	hub  *Hub            // websocket总线
	conn *websocket.Conn // 当前客户端连接
	send chan []byte
}

// 读取数据通道
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c // websocket总线断开注册
		c.conn.Close()        // 关闭当前客户端连接
	}()
	for {
		// 服务端读取此客户端的消息
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		// 将此消息广播给所有客户端
		c.hub.broadcast <- message
	}
}

// 发送数据通道
func (c *Client) writePump() {
	defer func() {
		c.conn.Close() // 关闭当前客户端连接
	}()
	for message := range c.send {
		var data = messageDataS{
			Event:   "message",
			Message: string(message),
		}
		dataJSON, _ := json.Marshal(data)
		c.conn.WriteMessage(websocket.TextMessage, dataJSON)
	}
}

// ws hub
type Hub struct {
	upGrader         *websocket.Upgrader
	clients          sync.Map     // map[*Client]bool
	broadcast        chan []byte  // 广播的信号管道
	register         chan *Client // 注册的信号管道
	unregister       chan *Client // 断开的信号管道
	aliveClientNum   int64        // 活跃客户端连接
	historyClientNum int64        // 历史客户端连接
}

func NewHub() *Hub {
	return &Hub{
		upGrader: &websocket.Upgrader{
			HandshakeTimeout: 60 * time.Second,
			ReadBufferSize:   256,
			WriteBufferSize:  256,
			WriteBufferPool:  nil,
			Subprotocols:     nil,
			Error:            nil,
			//CheckOrigin: func(r *http.Request) bool {
			//	return true
			//},
			CheckOrigin:       nil,
			EnableCompression: false,
		},
		clients:          sync.Map{},
		broadcast:        make(chan []byte),
		register:         make(chan *Client),
		unregister:       make(chan *Client),
		aliveClientNum:   0,
		historyClientNum: 0,
	}
}

// 运行hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register: // 收到client的注册
			// 注册事件处理
			h.clients.Store(client, true)
			// 更新客户端数
			atomic.AddInt64(&h.aliveClientNum, 1)
			atomic.AddInt64(&h.historyClientNum, 1)
			// 广播
			var data = connectedDataS{
				Event:            "connected",
				HistoryClientNum: h.historyClientNum,
				OnlineClientNum:  h.aliveClientNum,
			}
			var dataJSON, _ = json.Marshal(data)
			h.BroadcastAllClient(websocket.TextMessage, dataJSON)
		case client := <-h.unregister: // 收到client的断开
			if _, ok := h.clients.LoadAndDelete(client); ok {
				atomic.AddInt64(&h.aliveClientNum, -1)
				close(client.send)
				var data = connectedDataS{
					Event:            "connected",
					HistoryClientNum: h.historyClientNum,
					OnlineClientNum:  h.aliveClientNum,
				}
				var dataJSON, _ = json.Marshal(data)
				h.BroadcastAllClient(websocket.TextMessage, dataJSON)
			}
		case message := <-h.broadcast: // 广播的[]byte message
			h.clients.Range(func(key, value any) bool {
				client := key.(*Client)
				select {
				case client.send <- message:
				default:
					close(client.send)
					h.clients.Delete(client)
				}
				return true
			})
		}
	}
}

func (h *Hub) BroadcastAllClient(mt int, message []byte) {
	h.clients.Range(func(client, exists any) bool {
		client.(*Client).conn.WriteMessage(mt, message)
		return true
	})
}

func (h *Hub) CreateWS(c *gin.Context) {
	// 创建连接
	conn, err := h.upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// 创建新客户端
	client := &Client{
		hub:  h,
		conn: conn,
		send: make(chan []byte, 256),
	}
	// 异步启动客户端的读写
	go client.writePump()
	go client.readPump()
	// 注册新客户端
	client.hub.register <- client
}
