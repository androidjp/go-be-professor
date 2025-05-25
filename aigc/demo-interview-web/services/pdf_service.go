package services

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"regexp"
	"strings"

	"demo-interview-web/models"

	"github.com/ledongthuc/pdf"
)

// ExtractTextFromPDF 从PDF文件中提取文本内容
func ExtractTextFromPDF(pdfData []byte) (string, error) {
	// 创建一个读取器
	reader := bytes.NewReader(pdfData)

	// 使用ledongthuc/pdf库读取PDF
	r, err := pdf.NewReader(reader, int64(len(pdfData)))
	if err != nil {
		return "", err
	}

	// 获取页数
	numPages := r.NumPage()
	if numPages == 0 {
		return "", errors.New("PDF文件没有页面")
	}

	// 提取所有页面的文本
	var textBuilder strings.Builder
	for i := 1; i <= numPages; i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}

		text, err := p.GetPlainText(nil)
		if err != nil {
			slog.Info("提取PDF页文本出错", "page", i, "error", err)
			continue
		}

		textBuilder.WriteString(text)
	}

	return textBuilder.String(), nil
}

// ParseResumeFromText 从文本中解析简历信息
func ParseResumeFromText(text string) models.Resume {
	// 创建简历对象
	resume := models.Resume{
		RawText: text,
		Skills:  []string{},
	}

	// 解析基本信息
	resume.BasicInfo = parseBasicInfo(text)

	// 解析教育经历
	resume.Education = parseEducation(text)

	// 解析工作经验
	resume.WorkExperience = parseWorkExperience(text)

	// 解析技能
	resume.Skills = parseSkills(text)

	// 解析项目经验
	resume.Projects = parseProjects(text)

	// 解析证书
	resume.Certificates = parseCertificates(text)

	return resume
}

// parseBasicInfo 解析基本信息
func parseBasicInfo(text string) models.BasicInfo {
	basicInfo := models.BasicInfo{}

	// 尝试提取姓名 (通常在简历开头)
	nameRegex := regexp.MustCompile(`(?i)^\s*([\p{L}\s]{2,30})\s*$`)
	nameMatches := nameRegex.FindStringSubmatch(text)
	if len(nameMatches) > 1 {
		basicInfo.Name = strings.TrimSpace(nameMatches[1])
	}

	// 提取电子邮件
	emailRegex := regexp.MustCompile(`(?i)[\w._%+-]+@[\w.-]+\.[a-zA-Z]{2,}`)
	emailMatches := emailRegex.FindString(text)
	if emailMatches != "" {
		basicInfo.Email = emailMatches
	}

	// 提取电话号码
	phoneRegex := regexp.MustCompile(`(?i)(\+?\d{1,3}[-.\s]?)?\(?\d{3}\)?[-.\s]?\d{3}[-.\s]?\d{4}`)
	phoneMatches := phoneRegex.FindString(text)
	if phoneMatches != "" {
		basicInfo.Phone = phoneMatches
	}

	// 提取位置信息
	locationRegex := regexp.MustCompile(`(?i)(地址|位置|城市|地区|Address|Location|City)[:\s]*([一-龥a-zA-Z0-9\s,.-]+)`)
	locationMatches := locationRegex.FindStringSubmatch(text)
	if len(locationMatches) > 2 {
		basicInfo.Location = strings.TrimSpace(locationMatches[2])
	}

	// 提取职位头衔
	titleRegex := regexp.MustCompile(`(?i)(职位|职称|Title|Position)[:\s]*([\p{L}\s,.-]+)`)
	titleMatches := titleRegex.FindStringSubmatch(text)
	if len(titleMatches) > 2 {
		basicInfo.Title = strings.TrimSpace(titleMatches[2])
	}

	// 提取个人简介
	summaryRegex := regexp.MustCompile(`(?i)(简介|概述|自我介绍|个人简介|Summary|Profile|About)[:\s]*([\s\S]{10,300}?)(?:\n\n|\r\n\r\n)`)
	summaryMatches := summaryRegex.FindStringSubmatch(text)
	if len(summaryMatches) > 2 {
		basicInfo.Summary = strings.TrimSpace(summaryMatches[2])
	}

	return basicInfo
}

// parseEducation 解析教育经历
func parseEducation(text string) []models.Education {
	var educations []models.Education

	// 查找教育部分
	eduSectionRegex := regexp.MustCompile(`(?i)(教育|学历|教育背景|Education)[\s:]*([\s\S]+?)(?:\n\n|\r\n\r\n|工作经验|工作经历|Work Experience|技能|Skills|$)`)
	eduSection := eduSectionRegex.FindStringSubmatch(text)

	if len(eduSection) < 2 {
		return educations
	}

	eduText := eduSection[2]

	// 提取每段教育经历
	eduEntryRegex := regexp.MustCompile(`(?i)([\p{L}\s.-]+)\s*(?:大学|学院|University|College|Institute)?\s*[\n\r]?\s*([^\n\r]+)?\s*[\n\r]?\s*(?:([\d]{4})\s*[-–—至to]+\s*([\d]{4}|至今|现在|Present|Now|Current))?`)
	eduEntries := eduEntryRegex.FindAllStringSubmatch(eduText, -1)

	for _, entry := range eduEntries {
		if len(entry) > 1 {
			edu := models.Education{
				Institution: strings.TrimSpace(entry[1]),
			}

			if len(entry) > 2 && entry[2] != "" {
				// 尝试分离学位和专业
				degreeAndMajor := strings.TrimSpace(entry[2])
				parts := strings.Split(degreeAndMajor, " ")

				if len(parts) > 1 {
					edu.Degree = parts[0]
					edu.Major = strings.Join(parts[1:], " ")
				} else {
					edu.Major = degreeAndMajor
				}
			}

			if len(entry) > 3 && entry[3] != "" {
				edu.StartDate = entry[3]
			}

			if len(entry) > 4 && entry[4] != "" {
				edu.EndDate = entry[4]
			}

			educations = append(educations, edu)
		}
	}

	return educations
}

// parseWorkExperience 解析工作经验
func parseWorkExperience(text string) []models.Experience {
	var experiences []models.Experience

	// 查找工作经验部分
	workSectionRegex := regexp.MustCompile(`(?i)(工作经[验历]|Work Experience)[\s:]*([\s\S]+?)(?:\n\n|\r\n\r\n|教育|学历|Education|项目经验|Projects|技能|Skills|$)`)
	workSection := workSectionRegex.FindStringSubmatch(text)

	if len(workSection) < 2 {
		return experiences
	}

	workText := workSection[2]

	// 提取每段工作经历
	workEntryRegex := regexp.MustCompile(`(?i)([\p{L}\s&.-]+)\s*[\n\r]?\s*([^\n\r]+)?\s*[\n\r]?\s*(?:([\d]{4}[-/][\d]{1,2})\s*[-–—至to]+\s*([\d]{4}[-/][\d]{1,2}|至今|现在|Present|Now|Current))?\s*[\n\r]([\s\S]+?)(?:\n\n|\r\n\r\n|$)`)
	workEntries := workEntryRegex.FindAllStringSubmatch(workText, -1)

	for _, entry := range workEntries {
		if len(entry) > 4 {
			exp := models.Experience{
				Company: strings.TrimSpace(entry[1]),
			}

			if len(entry) > 2 && entry[2] != "" {
				exp.Title = strings.TrimSpace(entry[2])
			}

			if len(entry) > 3 && entry[3] != "" {
				exp.StartDate = entry[3]
			}

			if len(entry) > 4 && entry[4] != "" {
				exp.EndDate = entry[4]
			}

			if len(entry) > 5 && entry[5] != "" {
				// 分割描述为多行
				descLines := strings.Split(strings.TrimSpace(entry[5]), "\n")
				for _, line := range descLines {
					line = strings.TrimSpace(line)
					if line != "" {
						exp.Description = append(exp.Description, line)
					}
				}
			}

			experiences = append(experiences, exp)
		}
	}

	return experiences
}

// parseSkills 解析技能
func parseSkills(text string) []string {
	var skills []string

	// 查找技能部分
	skillSectionRegex := regexp.MustCompile(`(?i)(技能|专业技能|Skills|Technical Skills)[\s:]*([\s\S]+?)(?:\n\n|\r\n\r\n|教育|学历|Education|工作经验|Work Experience|项目经验|Projects|$)`)
	skillSection := skillSectionRegex.FindStringSubmatch(text)

	if len(skillSection) < 2 {
		return skills
	}

	skillText := skillSection[2]

	// 提取技能列表
	// 尝试匹配不同格式的技能列表
	skillPatterns := []string{
		`[•·\-*]\s*([^\n\r•·\-*]+)`, // 带有项目符号的列表
		`([\w\s\+\#]+(?:,|，|;|；))`,  // 逗号或分号分隔的列表
	}

	for _, pattern := range skillPatterns {
		skillRegex := regexp.MustCompile(pattern)
		skillMatches := skillRegex.FindAllStringSubmatch(skillText, -1)

		if len(skillMatches) > 0 {
			for _, match := range skillMatches {
				if len(match) > 1 {
					skill := strings.TrimSpace(match[1])
					skill = strings.TrimRight(skill, ",，;；")
					if skill != "" {
						skills = append(skills, skill)
					}
				}
			}
			break // 如果找到匹配项，就不再尝试其他模式
		}
	}

	// 如果上面的模式都没有匹配到，尝试按行分割
	if len(skills) == 0 {
		lines := strings.Split(skillText, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				skills = append(skills, line)
			}
		}
	}

	return skills
}

// parseProjects 解析项目经验
func parseProjects(text string) []models.Project {
	var projects []models.Project

	// 查找项目经验部分
	projectSectionRegex := regexp.MustCompile(`(?i)(项目经[验历]|Projects)[\s:]*([\s\S]+?)(?:\n\n|\r\n\r\n|教育|学历|Education|工作经验|Work Experience|技能|Skills|$)`)
	projectSection := projectSectionRegex.FindStringSubmatch(text)

	if len(projectSection) < 2 {
		return projects
	}

	projectText := projectSection[2]

	// 提取每个项目
	projectEntryRegex := regexp.MustCompile(`(?i)([\p{L}\s&.-]+)\s*[\n\r]?\s*(?:([\d]{4}[-/][\d]{1,2})\s*[-–—至to]+\s*([\d]{4}[-/][\d]{1,2}|至今|现在|Present|Now|Current))?\s*[\n\r]([\s\S]+?)(?:\n\n|\r\n\r\n|$)`)
	projectEntries := projectEntryRegex.FindAllStringSubmatch(projectText, -1)

	for _, entry := range projectEntries {
		if len(entry) > 1 {
			proj := models.Project{
				Name: strings.TrimSpace(entry[1]),
			}

			if len(entry) > 2 && entry[2] != "" {
				proj.StartDate = entry[2]
			}

			if len(entry) > 3 && entry[3] != "" {
				proj.EndDate = entry[3]
			}

			if len(entry) > 4 && entry[4] != "" {
				// 提取描述
				descText := strings.TrimSpace(entry[4])
				proj.Description = descText

				// 尝试从描述中提取技术栈
				techRegex := regexp.MustCompile(`(?i)(技术栈|使用技术|Technologies|Tech Stack)[:\s]*([^\n\r]+)`)
				techMatch := techRegex.FindStringSubmatch(descText)
				if len(techMatch) > 2 {
					techText := techMatch[2]
					techs := strings.Split(techText, ",")
					for _, tech := range techs {
						tech = strings.TrimSpace(tech)
						if tech != "" {
							proj.Technologies = append(proj.Technologies, tech)
						}
					}
				}
			}

			projects = append(projects, proj)
		}
	}

	return projects
}

// parseCertificates 解析证书
func parseCertificates(text string) []string {
	var certificates []string

	// 查找证书部分
	certSectionRegex := regexp.MustCompile(`(?i)(证书|资格证书|认证|Certifications|Certificates)[\s:]*([\s\S]+?)(?:\n\n|\r\n\r\n|教育|学历|Education|工作经验|Work Experience|技能|Skills|项目经验|Projects|$)`)
	certSection := certSectionRegex.FindStringSubmatch(text)

	if len(certSection) < 2 {
		return certificates
	}

	certText := certSection[2]

	// 提取证书列表
	certPatterns := []string{
		`[•·\-*]\s*([^\n\r•·\-*]+)`, // 带有项目符号的列表
		`([\w\s\+\#]+(?:,|，|;|；))`,  // 逗号或分号分隔的列表
	}

	for _, pattern := range certPatterns {
		certRegex := regexp.MustCompile(pattern)
		certMatches := certRegex.FindAllStringSubmatch(certText, -1)

		if len(certMatches) > 0 {
			for _, match := range certMatches {
				if len(match) > 1 {
					cert := strings.TrimSpace(match[1])
					cert = strings.TrimRight(cert, ",，;；")
					if cert != "" {
						certificates = append(certificates, cert)
					}
				}
			}
			break // 如果找到匹配项，就不再尝试其他模式
		}
	}

	// 如果上面的模式都没有匹配到，尝试按行分割
	if len(certificates) == 0 {
		lines := strings.Split(certText, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				certificates = append(certificates, line)
			}
		}
	}

	return certificates
}

// ReadPdf 从io.Reader读取PDF内容
func ReadPdf(reader io.Reader) (string, error) {
	pdf.DebugOn = false

	// 读取所有内容
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return ExtractTextFromPDF(data)
}
