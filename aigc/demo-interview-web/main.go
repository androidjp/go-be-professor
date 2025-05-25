package main

import (
	"log/slog"
	"net/http"
	"os"

	"demo-interview-web/handlers"
	"demo-interview-web/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// slog 设置为JSON方式打印
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})))

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建Gin引擎
	r := gin.Default()

	// 配置CORS
	r.Use(middleware.CORSMiddleware())

	// 使用自定义中间件
	r.Use(middleware.Logger())

	// API路由组
	api := r.Group("/api")
	{
		// 开发API路由组
		dev := api.Group("/dev/v1")
		{
			// PDF提取接口
			dev.POST("/pdf/extract", handlers.ExtractPDFHandler)

			// 用户信息总结接口
			dev.POST("/userinfo/summary", handlers.SummarizeUserInfoHandler)
		}

		// 面试信息接口
		api.POST("/v1/interview/info/get_all", handlers.GetInterviewQuestionsHandler)
		api.POST("/v1/gen_question", handlers.GenQuestionHandler)
		api.POST("/v1/answer_and_next", handlers.AnswerAndNextHandler)
	}

	// 静态文件服务
	r.StaticFile("/", "./index.html")
	r.StaticFile("/ask", "./ask.html")

	// 启动服务器
	serverAddr := ":8222"
	slog.Info("服务器启动", "addr", serverAddr)
	if err := http.ListenAndServe(serverAddr, r); err != nil {
		slog.Error("服务器启动失败", "error", err)
	}
}
