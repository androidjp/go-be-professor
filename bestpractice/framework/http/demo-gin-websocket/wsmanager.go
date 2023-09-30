package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"sync"
	"sync/atomic"
)

type WebSocketManager struct {
	// ws upGrader
	upGrader *websocket.Upgrader
	// ws客户端管理 -> map[*websocket.Conn]bool
	wsClients sync.Map
	// 客户端存活数量
	aliveClientNum int64
	// 所有历史客户端数量
	historyClientNum int64
}

// 广播消息
func (m *WebSocketManager) BroadcastMsg(mt int, message []byte) {
	m.wsClients.Range(func(client, value any) bool {
		client.(*websocket.Conn).WriteMessage(mt, message)
		return true
	})
}

// 存储当前客户端
func (m *WebSocketManager) SaveClient(cli *websocket.Conn) {
	m.wsClients.Store(cli, true)
	// 更新存活客户端数
	atomic.AddInt64(&m.aliveClientNum, 1)
	// 更新历史客户端数
	m.historyClientNum++
	// 广播进入
	var data = connectedDataS{
		Event:            "connected",
		HistoryClientNum: m.historyClientNum,
		OnlineClientNum:  m.aliveClientNum,
	}
	var dataJSON, _ = json.Marshal(data)
	m.BroadcastMsg(websocket.TextMessage, dataJSON)
}

// 删除当前客户端
func (m *WebSocketManager) DelClient(cli *websocket.Conn) {
	// 断开连接，删除client
	m.wsClients.Delete(cli)
	// 更新存活客户端数
	atomic.AddInt64(&m.aliveClientNum, -1)
	// 广播离开
	var data = connectedDataS{
		Event:            "disconnected",
		HistoryClientNum: m.historyClientNum,
		OnlineClientNum:  m.aliveClientNum,
	}
	var dataJSON, _ = json.Marshal(data)
	m.BroadcastMsg(websocket.TextMessage, dataJSON)
}

// 发送消息
func (m *WebSocketManager) SendMsg(mt int, message []byte) {
	var data = messageDataS{
		Event:   "message",
		Message: string(message),
	}
	dataJson, _ := json.Marshal(data)
	m.BroadcastMsg(mt, dataJson)
}

type connectedDataS struct {
	Event            string `json:"_event"`
	HistoryClientNum int64  `json:"historyClientNum"`
	OnlineClientNum  int64  `json:"onlineClientNum"`
}

type messageDataS struct {
	Event   string `json:"_event"`
	Message string `json:"message"`
}
