package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"demo-interview-web/models"
)

const (
	// OllamaAPIURL 本地Ollama API地址
	OllamaAPIURL = "http://localhost:11434/api/generate"

	// DefaultModel 默认使用的模型
	DefaultModel = "qwen2.5:latest"
	// DefaultModel = "llama3.2:latest"
)

// SummarizeResume 使用Ollama总结简历信息
func SummarizeResume(resume models.Resume) (string, error) {
	// 构建提示词
	prompt := buildResumeSummaryPrompt(resume)

	// 调用Ollama API
	response, err := callOllamaAPI(prompt, false)
	if err != nil {
		return "", err
	}

	return response, nil
}

// GenerateInterviewQuestions 生成面试问题（返回数组）
func GenerateInterviewQuestions(resumeSummary, jobTitle, companyName, jobRequirements string) ([]models.InterviewQuestion, error) {
	prompt := buildInterviewQuestionsPrompt(resumeSummary, jobTitle, companyName, jobRequirements)
	response, err := callOllamaAPI(prompt, false)
	if err != nil {
		return nil, err
	}
	// 尝试提取第一个JSON数组
	jsonStart := strings.Index(response, "[")
	jsonEnd := strings.LastIndex(response, "]")
	if jsonStart != -1 && jsonEnd != -1 && jsonEnd > jsonStart {
		jsonStr := response[jsonStart : jsonEnd+1]
		var questions []models.InterviewQuestion
		err = json.Unmarshal([]byte(jsonStr), &questions)
		if err != nil {
			return nil, fmt.Errorf("解析Ollama返回的面试题JSON失败: %w, 原始内容: %s", err, response)
		}
		return questions, nil
	}
	return nil, fmt.Errorf("未找到合法的JSON数组，原始内容: %s", response)
}

// GenerateSingleQuestion 根据知识领域生成一个面试问题
func GenerateSingleQuestion(domain string, prevQuestion ...string) (models.InterviewQuestion, error) {
	prompt := buildSingleQuestionPrompt(domain, prevQuestion...)
	response, err := callOllamaAPI(prompt, false)
	if err != nil {
		return models.InterviewQuestion{}, err
	}
	slog.Info("生成单题完毕，返回的内容是", "response", response)
	// 尝试正则提取第一个JSON对象
	jsonStr := extractJSONObject(response)
	if jsonStr == "" {
		jsonStart := strings.Index(response, "{")
		jsonEnd := strings.LastIndex(response, "}")
		if jsonStart != -1 && jsonEnd != -1 && jsonEnd > jsonStart {
			jsonStr = response[jsonStart : jsonEnd+1]
		}
	}
	if jsonStr != "" {
		jsonStr = cleanJSONStr(jsonStr)
		var question models.InterviewQuestion
		err = json.Unmarshal([]byte(jsonStr), &question)
		if err == nil {
			return question, nil
		}
		// 尝试修复常见格式问题后再次解析
		fixed := fixJSONStr(jsonStr)
		err2 := json.Unmarshal([]byte(fixed), &question)
		if err2 == nil {
			return question, nil
		}
		return models.InterviewQuestion{}, fmt.Errorf("解析Ollama返回的面试题JSON失败: %w, 修正后: %v, 原始内容: %s", err2, fixed, response)
	}
	return models.InterviewQuestion{}, fmt.Errorf("未找到合法的JSON对象，原始内容: %s", response)
}

// extractJSONObject 用正则提取第一个JSON对象
func extractJSONObject(s string) string {
	re := regexp.MustCompile(`(?s)\{.*?\}`)
	match := re.FindString(s)
	return match
}

// cleanJSONStr 去除多余换行、空格、逗号
func cleanJSONStr(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.TrimSpace(s)
	// 去除多余逗号（如最后一个字段后）
	s = regexp.MustCompile(`,\s*}`).ReplaceAllString(s, "}")
	return s
}

// fixJSONStr 尝试修复Ollama返回的JSON字符串中的常见格式问题
func fixJSONStr(s string) string {
	s = strings.ReplaceAll(s, "`", "\"") // 反引号转双引号
	s = strings.ReplaceAll(s, "“", "\"")
	s = strings.ReplaceAll(s, "”", "\"")
	s = strings.ReplaceAll(s, "‘", "\"")
	s = strings.ReplaceAll(s, "’", "\"")
	s = strings.ReplaceAll(s, "，", ",") // 中文逗号转英文
	s = strings.ReplaceAll(s, "：", ":") // 中文冒号转英文
	s = strings.ReplaceAll(s, "。", ".") // 中文句号转英文
	return s
}

// AnswerAndNext 分析答案得分并生成下一个问题
func AnswerAndNext(domain string, question models.InterviewQuestion, answer string) (score int, nextQ *models.InterviewQuestion, scoreReason string, referenceAnswer string, err error) {
	// 1. 让ollama细分评分
	judgePrompt := buildScoreAnswerPrompt(question, answer)
	judgeResp, err := callOllamaAPI(judgePrompt, false)
	if err != nil {
		return 0, nil, "", "", err
	}
	slog.Info("answer_and_next 评分原始返回", "judgeResp", judgeResp)
	score = extractScore(judgeResp)

	// 2. 获取得分理由
	reasonPrompt := buildScoreReasonPrompt(question, answer)
	reasonResp, err := callOllamaAPI(reasonPrompt, false)
	if err != nil {
		scoreReason = ""
	} else {
		scoreReason = strings.TrimSpace(reasonResp)
	}

	// 3. 获取参考答案
	refPrompt := buildReferenceAnswerPrompt(question)
	refResp, err := callOllamaAPI(refPrompt, false)
	if err != nil {
		referenceAnswer = ""
	} else {
		referenceAnswer = strings.TrimSpace(refResp)
	}

	// 4. 生成下一个问题，传入上一个问题内容，避免重复
	next, err := GenerateSingleQuestion(domain, question.Question)
	if err != nil {
		return score, nil, scoreReason, referenceAnswer, nil // 允许最后一个问题没有下一个
	}
	return score, &next, scoreReason, referenceAnswer, nil
}

// buildScoreAnswerPrompt 构建细分评分的提示词
func buildScoreAnswerPrompt(question models.InterviewQuestion, answer string) string {
	return `你是一位专业的面试官，请根据下面的评分规则，对考生的答案进行打分，只输出0、1、2、3或4中的一个数字，不要输出任何解释：\n\n` +
		`评分规则：\n` +
		`得分为0：表示回答完全和问题不相关，且回答本身的理论知识或者结论是错误的。\n` +
		`得分为1：表示回答与问题不相关，但是回答本身的理论知识或者结论是正确的。\n` +
		`得分为2：表示回答与问题相关，但是回答本身的理论知识或者结论是错误的。\n` +
		`得分为3：表示回答与问题相关，且回答本身的理论知识或者结论是正确的，但是且回答是内容不够精炼或者不够高质量。\n` +
		`得分为4：表示回答完全与问题对应，且回答本身的理论知识或者结论是正确的，且回答是内容是精炼且高质量的。\n` +
		`\n问题：` + question.Question + `\n考生答案：` + answer
}

// buildScoreReasonPrompt 构建得分理由提示词
func buildScoreReasonPrompt(question models.InterviewQuestion, answer string) string {
	return `你是一位专业的面试官，请根据评分规则，简要说明考生答案得分的理由，突出相关性、理论正确性、内容质量等关键点，50字以内：\n\n问题：` + question.Question + `\n考生答案：` + answer
}

// buildReferenceAnswerPrompt 构建参考答案提示词
func buildReferenceAnswerPrompt(question models.InterviewQuestion) string {
	return `请给出该面试问题的高质量参考答案，要求内容精炼、准确、专业，150字以内：\n\n问题：` + question.Question
}

// extractScore 从Ollama返回内容中提取0-4分
func extractScore(s string) int {
	re := regexp.MustCompile(`[0-4]`)
	match := re.FindString(s)
	if match != "" {
		return int(match[0] - '0')
	}
	return 0
}

// buildResumeSummaryPrompt 构建简历总结提示词
func buildResumeSummaryPrompt(resume models.Resume) string {
	// 构建提示词
	var sb strings.Builder

	sb.WriteString("你是一位专业的人力资源专家，请根据以下简历信息，提炼出求职者的核心技能、经验和特点，用简洁的语言总结。\n\n")
	sb.WriteString("简历信息:\n")

	// 添加基本信息
	sb.WriteString(fmt.Sprintf("姓名: %s\n", resume.BasicInfo.Name))
	if resume.BasicInfo.Title != "" {
		sb.WriteString(fmt.Sprintf("职位: %s\n", resume.BasicInfo.Title))
	}
	if resume.BasicInfo.Summary != "" {
		sb.WriteString(fmt.Sprintf("个人简介: %s\n", resume.BasicInfo.Summary))
	}

	// 添加教育经历
	if len(resume.Education) > 0 {
		sb.WriteString("\n教育经历:\n")
		for _, edu := range resume.Education {
			sb.WriteString(fmt.Sprintf("- %s, %s %s (%s - %s)\n",
				edu.Institution, edu.Degree, edu.Major, edu.StartDate, edu.EndDate))
		}
	}

	// 添加工作经验
	if len(resume.WorkExperience) > 0 {
		sb.WriteString("\n工作经验:\n")
		for _, exp := range resume.WorkExperience {
			sb.WriteString(fmt.Sprintf("- %s, %s (%s - %s)\n",
				exp.Company, exp.Title, exp.StartDate, exp.EndDate))
			for _, desc := range exp.Description {
				sb.WriteString(fmt.Sprintf("  * %s\n", desc))
			}
		}
	}

	// 添加技能
	if len(resume.Skills) > 0 {
		sb.WriteString("\n技能:\n")
		for _, skill := range resume.Skills {
			sb.WriteString(fmt.Sprintf("- %s\n", skill))
		}
	}

	// 添加项目经验
	if len(resume.Projects) > 0 {
		sb.WriteString("\n项目经验:\n")
		for _, proj := range resume.Projects {
			sb.WriteString(fmt.Sprintf("- %s\n", proj.Name))
			sb.WriteString(fmt.Sprintf("  描述: %s\n", proj.Description))
			if len(proj.Technologies) > 0 {
				sb.WriteString(fmt.Sprintf("  技术: %s\n", strings.Join(proj.Technologies, ", ")))
			}
		}
	}

	// 添加证书
	if len(resume.Certificates) > 0 {
		sb.WriteString("\n证书:\n")
		for _, cert := range resume.Certificates {
			sb.WriteString(fmt.Sprintf("- %s\n", cert))
		}
	}

	sb.WriteString("\n请提炼出这位求职者的核心技能、经验和特点，用简洁的语言总结（不超过300字）。")

	return sb.String()
}

// buildInterviewQuestionsPrompt 构建面试问题生成提示词
func buildInterviewQuestionsPrompt(resumeSummary, jobTitle, companyName, jobRequirements string) string {
	// 构建提示词
	var sb strings.Builder

	sb.WriteString("你是一位经验丰富的技术面试官，需要根据候选人的简历和职位要求，生成20个有针对性的面试问题。\n\n")

	// 添加职位信息
	sb.WriteString(fmt.Sprintf("职位: %s\n", jobTitle))
	sb.WriteString(fmt.Sprintf("公司: %s\n", companyName))
	sb.WriteString(fmt.Sprintf("职位要求:\n%s\n\n", jobRequirements))

	// 添加候选人信息
	sb.WriteString(fmt.Sprintf("候选人简历摘要:\n%s\n\n", resumeSummary))

	sb.WriteString("请严格按照JSON数组格式输出20个面试问题，每个元素为一个对象，包含如下字段：category（问题类别）、difficulty（难度）、question（问题内容）、answer（参考答案）。请保证所有输出尽可能为中文，且question和answer字段不为空，且answer尽可能详细具体，最终结果输出一个没有换行符转义的、压缩后的JSON数组。\n")

	return sb.String()
}

// buildSingleQuestionPrompt 构建单题生成提示词
func buildSingleQuestionPrompt(domain string, prevQuestion ...string) string {
	var prev string
	if len(prevQuestion) > 0 && prevQuestion[0] != "" {
		prev = prevQuestion[0]
	}
	if prev != "" {
		return fmt.Sprintf(`你是一位专业的面试官，请根据以下知识领域，生成一个全新的、有针对性的面试问题，且不要与上一个问题重复或高度相似。\n\n知识领域: %s\n上一个问题: %s\n\n请严格以JSON对象格式输出，包含如下字段：category（问题类别）、difficulty（难度）、question（问题内容）。请保证只输出一个JSON对象，且question字段内容要详细具体、有深度和针对性，避免泛泛而谈，最终结果输出一个没有换行符转义的、压缩后的JSON对象。`, domain, prev)
	}
	return fmt.Sprintf(`你是一位专业的面试官，请根据以下知识领域，生成一个有针对性的面试问题。\n\n知识领域: %s\n\n请严格以JSON对象格式输出，包含如下字段：category（问题类别）、difficulty（难度）、question（问题内容）。请保证只输出一个JSON对象，最终结果输出一个没有换行符转义的、压缩后的JSON对象。`, domain)
}

// callOllamaAPI 调用Ollama API
func callOllamaAPI(prompt string, stream bool) (string, error) {
	// 构建请求体
	request := models.OllamaRequest{
		Model:  DefaultModel,
		Prompt: prompt,
		Stream: stream,
	}

	// 序列化请求体
	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", OllamaAPIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API返回错误状态码 %d: %s", resp.StatusCode, string(body))
	}

	// 处理响应
	if stream {
		// 流式处理逻辑
		return handleStreamResponse(resp.Body)
	} else {
		// 非流式处理逻辑
		return handleNonStreamResponse(resp.Body)
	}
}

// handleStreamResponse 处理流式响应
func handleStreamResponse(body io.ReadCloser) (string, error) {
	var fullResponse strings.Builder
	decoder := json.NewDecoder(body)

	for {
		var response models.OllamaResponse
		err := decoder.Decode(&response)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fullResponse.String(), fmt.Errorf("解析响应失败: %w", err)
		}

		fullResponse.WriteString(response.Response)

		if response.Done {
			break
		}
	}

	return fullResponse.String(), nil
}

// handleNonStreamResponse 处理非流式响应
func handleNonStreamResponse(body io.ReadCloser) (string, error) {
	var response models.OllamaResponse

	// 读取响应体
	respBody, err := io.ReadAll(body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %w", err)
	}

	// 解析JSON
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return "", fmt.Errorf("解析JSON失败: %w", err)
	}

	return response.Response, nil
}
