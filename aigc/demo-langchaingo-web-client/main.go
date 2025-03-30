package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

// 全局LLM客户端
var llmClient *ollama.LLM

// 初始化LLM客户端
func initLLM() error {
	var err error
	llmClient, err = ollama.New(ollama.WithModel("qwen2:7b"))
	if err != nil {
		return fmt.Errorf("failed to initialize LLM: %v", err)
	}
	return nil
}

// 定义响应结构体
type ExpertContext struct {
	SessionID string `json:"session_id"`
	SrvReqID  string `json:"srv_req_id"`
}

type OutputStrategy struct {
	Mode          int    `json:"mode"`
	WithTitle     bool   `json:"with_title"`
	ContentFormat string `json:"content_format"`
}

type RetrievedResult struct {
	Title string `json:"title"`
}

type Value struct {
	Val string `json:"val"`
}

type NeedConfirm struct {
	Key        string  `json:"key"`
	Desc       string  `json:"desc"`
	ChoiceType int     `json:"choice_type"`
	Vals       []Value `json:"vals"`
}

type ResponseData struct {
	ExpertCtx       ExpertContext     `json:"expert_ctx"`
	OutputStrategy  OutputStrategy    `json:"output_strategy"`
	RetrievedResult []RetrievedResult `json:"retrieved_result"`
	NeedConfirm     []NeedConfirm     `json:"need_confirm"`
}

type Response struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data ResponseData `json:"data"`
}

func chatHandler(c *gin.Context) {
	var request struct {
		Question string `json:"question"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if request.Question == "" {
		log.Print("Empty question received")
		c.JSON(http.StatusBadRequest, gin.H{"error": "question cannot be empty"})
		return
	}

	// 构造mock数据
	mockResp := Response{
		Data: ResponseData{
			ExpertCtx: ExpertContext{
				SessionID: "55555",
				SrvReqID:  "请求唯一ID",
			},
			OutputStrategy: OutputStrategy{
				Mode:          0,
				WithTitle:     true,
				ContentFormat: "segment/page",
			},
			RetrievedResult: []RetrievedResult{
				{Title: "旅行计划"},
			},
			NeedConfirm: []NeedConfirm{
				{
					Key:        "sleep_days",
					Desc:       "居住天数",
					ChoiceType: 0,
					Vals: []Value{
						{Val: "5天"},
						{Val: "3天"},
					},
				},
				{
					Key:        "budge_from",
					Desc:       "旅游总预算的下限",
					ChoiceType: 0,
					Vals: []Value{
						{Val: "无下限"},
						{Val: "2000元"},
					},
				},
				{
					Key:        "budge_to",
					Desc:       "旅游总预算的上限",
					ChoiceType: 0,
					Vals: []Value{
						{Val: "10000元"},
					},
				},
			},
		},
	}

	// 检查need_confirm是否存在并设置响应
	if len(mockResp.Data.NeedConfirm) > 0 {
		c.Writer.Header().Set("need_confirm", "true")
		jsonBs, _ := json.Marshal(mockResp.Data.NeedConfirm)
		c.Writer.Header().Set("need_confirm_data", base64.StdEncoding.EncodeToString(jsonBs))
	}

	ctx := context.Background()

	// 设置 SSE 头部
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Flush()

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, request.Question),
	}

	// 调用流式 API
	_, err := llmClient.GenerateContent(ctx, content, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Fprintf(c.Writer, "data: %s\n\n", string(chunk))
		c.Writer.Flush()
		return nil
	}))

	if err != nil {
		log.Printf("Failed to generate content: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get response"})
		return
	}

	fmt.Fprintln(c.Writer, "data: [DONE]\n")
	c.Writer.Flush()
}

// ConfirmRequest 结构体用于解析确认请求体
type ConfirmRequest struct {
	Confirmed []struct {
		Key  string `json:"key"`
		Vals []struct {
			Val string `json:"val"`
		} `json:"vals"`
	} `json:"confirmed"`
}

// confirmHandler 处理 /api/confirm 路由的请求
func confirmHandler(c *gin.Context) {
	var confirmReq ConfirmRequest
	if err := c.ShouldBindJSON(&confirmReq); err != nil {
		log.Printf("Invalid confirm request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// 构造mock数据
	mockResp := Response{
		Data: ResponseData{
			ExpertCtx: ExpertContext{
				SessionID: "55555",
				SrvReqID:  "请求唯一ID",
			},
			OutputStrategy: OutputStrategy{
				Mode:          0,
				WithTitle:     true,
				ContentFormat: "segment/page",
			},
			RetrievedResult: []RetrievedResult{
				{Title: "旅行计划"},
			},
			NeedConfirm: []NeedConfirm{
				{
					Key:        "sleep_days",
					Desc:       "居住天数",
					ChoiceType: 0,
					Vals: []Value{
						{Val: "5天"},
						{Val: "3天"},
					},
				},
				{
					Key:        "budge_from",
					Desc:       "旅游总预算的下限",
					ChoiceType: 0,
					Vals: []Value{
						{Val: "无下限"},
						{Val: "2000元"},
					},
				},
				{
					Key:        "budge_to",
					Desc:       "旅游总预算的上限",
					ChoiceType: 0,
					Vals: []Value{
						{Val: "10000元"},
					},
				},
			},
		},
	}

	// 检查need_confirm是否存在并设置响应
	if len(mockResp.Data.NeedConfirm) > 0 {
		c.Writer.Header().Set("need_confirm", "true")
		jsonBs, _ := json.Marshal(mockResp.Data.NeedConfirm)
		c.Writer.Header().Set("need_confirm_data", base64.StdEncoding.EncodeToString(jsonBs))
	}

	// 设置 SSE 头部
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Flush()

	// 发送 SSE 消息
	fmt.Fprintf(c.Writer, "data: 我成功收到了你的请求，正在运算中。。。\n\n")
	c.Writer.Flush()
	fmt.Fprintln(c.Writer, "data: [DONE]\n")
	c.Writer.Flush()

	// 处理确认数据
	for _, item := range confirmReq.Confirmed {
		log.Printf("Confirmed key: %s, value: %s", item.Key, item.Vals[0].Val)
	}
	// 可以在这里继续处理运算逻辑
}

func main() {
	// 初始化LLM客户端
	if err := initLLM(); err != nil {
		log.Fatalf("Failed to initialize LLM: %v", err)
	}

	// 设置gin模式
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	// 添加静态文件服务
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("static/index.html")
	})
	r.POST("/api/chat", chatHandler)
	r.POST("/api/confirm", confirmHandler) // 添加新的路由

	port := os.Getenv("PORT")
	if port == "" {
		port = "9527"
	}

	log.Printf("Server is running on http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
