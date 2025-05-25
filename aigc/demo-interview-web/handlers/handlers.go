package handlers

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"demo-interview-web/models"
	"demo-interview-web/services"

	"github.com/gin-gonic/gin"
)

// ExtractPDFHandler 处理PDF内容提取请求
func ExtractPDFHandler(c *gin.Context) {
	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "未找到上传的文件",
		})
		return
	}
	defer file.Close()

	// 检查文件类型
	if header.Header.Get("Content-Type") != "application/pdf" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "只支持PDF文件",
		})
		return
	}

	// 读取文件内容
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "读取文件失败",
		})
		return
	}

	// 提取PDF内容
	text, err := services.ExtractTextFromPDF(fileBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("提取PDF内容失败: %v", err),
		})
		return
	}

	// 解析简历信息
	resume := services.ParseResumeFromText(text)

	// 返回结构化的简历信息
	c.JSON(http.StatusOK, resume)
}

// SummarizeUserInfoHandler 处理用户信息总结请求
func SummarizeUserInfoHandler(c *gin.Context) {
	// 解析请求体
	var request models.SummaryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求格式",
		})
		return
	}

	// 使用Ollama总结简历信息
	summary, err := services.SummarizeResume(request.ResumeData)
	if err != nil {
		slog.Error("总结简历信息失败", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("总结简历信息失败: %v", err),
		})
		return
	}

	slog.Info("总结简历信息成功", "summary", summary)
	// 返回总结结果
	c.JSON(http.StatusOK, models.SummaryResponse{
		Summary: summary,
	})
}

// GetInterviewQuestionsHandler 处理面试问题生成请求（普通HTTP POST接口）
func GetInterviewQuestionsHandler(c *gin.Context) {
	// 解析请求体
	var request models.InterviewRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式"})
		return
	}

	// 打印所有入参
	slog.Info("/api/v1/interview/info/get_all 入参", "ResumeSummary", request.ResumeSummary, "JobTitle", request.JobTitle, "CompanyName", request.CompanyName, "JobRequirements", request.JobRequirements)

	// 生成面试问题
	questions, err := services.GenerateInterviewQuestions(
		request.ResumeSummary,
		request.JobTitle,
		request.CompanyName,
		request.JobRequirements,
	)
	if err != nil {
		slog.Error("生成面试问题失败", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("生成面试问题失败: %v", err)})
		return
	}

	slog.Info("生成面试问题成功", "questions", questions)

	c.JSON(http.StatusOK, questions)
}

// GenQuestionHandler 生成单个面试问题
func GenQuestionHandler(c *gin.Context) {
	var req models.GenQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式或缺少领域"})
		return
	}
	// 支持传入上一个问题内容，避免重复
	prevQ := c.Query("prev_question")
	question, err := services.GenerateSingleQuestion(req.Domain, prevQ)
	if err != nil {
		slog.Error("生成单题失败", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("生成单题失败: %v", err)})
		return
	}
	slog.Info("生成单题成功", "question", question)
	c.JSON(http.StatusOK, question)
}

// /api/v1/answer_and_next 答题并生成下一个问题
func AnswerAndNextHandler(c *gin.Context) {
	var req models.AnswerAndNextRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Domain == "" || req.Answer == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求格式或缺少参数"})
		return
	}
	// 传递上一个问题内容，避免重复
	score, nextQ, scoreReason, referenceAnswer, err := services.AnswerAndNext(req.Domain, req.Question, req.Answer)
	if err != nil {
		slog.Error("答题与生成下一个问题失败", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := models.AnswerAndNextResponse{
		Score:           score,
		NextQuestion:    nextQ,
		ScoreReason:     scoreReason,
		ReferenceAnswer: referenceAnswer,
	}
	slog.Info("答题分析结果", "score_reason", scoreReason, "reference_answer", referenceAnswer)
	c.JSON(http.StatusOK, resp)
}
