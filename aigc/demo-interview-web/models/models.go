package models

// Resume 简历数据结构
type Resume struct {
	BasicInfo      BasicInfo    `json:"basic_info"`
	Education      []Education  `json:"education"`
	WorkExperience []Experience `json:"work_experience"`
	Skills         []string     `json:"skills"`
	Projects       []Project    `json:"projects"`
	Certificates   []string     `json:"certificates"`
	RawText        string       `json:"raw_text"` // 原始文本内容
}

// BasicInfo 基本信息
type BasicInfo struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Location string `json:"location"`
	Title    string `json:"title"` // 职位头衔
	Summary  string `json:"summary"`
}

// Education 教育经历
type Education struct {
	Institution string `json:"institution"`
	Degree      string `json:"degree"`
	Major       string `json:"major"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	GPA         string `json:"gpa,omitempty"`
}

// Experience 工作经验
type Experience struct {
	Company     string   `json:"company"`
	Title       string   `json:"title"`
	StartDate   string   `json:"start_date"`
	EndDate     string   `json:"end_date"`
	Description []string `json:"description"`
}

// Project 项目经验
type Project struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Technologies []string `json:"technologies"`
	StartDate    string   `json:"start_date,omitempty"`
	EndDate      string   `json:"end_date,omitempty"`
}

// SummaryRequest 用户信息总结请求
type SummaryRequest struct {
	ResumeData Resume `json:"resume_data"`
}

// SummaryResponse 用户信息总结响应
type SummaryResponse struct {
	Summary string `json:"summary"`
}

// InterviewRequest 面试信息请求
type InterviewRequest struct {
	ResumeSummary   string `json:"resume_summary"`   // 简历摘要
	JobTitle        string `json:"job_title"`        // 职位名称
	CompanyName     string `json:"company_name"`     // 公司名称
	JobRequirements string `json:"job_requirements"` // 职位要求
}

// InterviewQuestion 面试问题
type InterviewQuestion struct {
	Category   string `json:"category"`   // 问题类别
	Difficulty string `json:"difficulty"` // 难度
	Question   string `json:"question"`   // 问题内容
	Answer     string `json:"answer"`     // 参考答案
}

// OllamaRequest Ollama API请求
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// OllamaResponse Ollama API响应
type OllamaResponse struct {
	Model    string `json:"model"`
	Created  string `json:"created_at"`
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// GenQuestionRequest 用于生成单个问题的请求
type GenQuestionRequest struct {
	Domain string `json:"domain"`
}

// 用于答题和生成下一个问题的请求
type AnswerAndNextRequest struct {
	Domain   string            `json:"domain"`
	Question InterviewQuestion `json:"question"`
	Answer   string            `json:"answer"`
}

// 用于答题和生成下一个问题的响应
type AnswerAndNextResponse struct {
	Score           int                `json:"score"`
	NextQuestion    *InterviewQuestion `json:"next_question,omitempty"`
	ScoreReason     string             `json:"score_reason,omitempty"`
	ReferenceAnswer string             `json:"reference_answer,omitempty"`
}
