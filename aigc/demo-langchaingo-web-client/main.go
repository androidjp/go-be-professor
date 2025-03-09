package main

import (
	"context"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "9527"
	}

	log.Printf("Server is running on http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
