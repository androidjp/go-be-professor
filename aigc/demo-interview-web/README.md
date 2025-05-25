# 面试问题生成器后端服务

这是一个基于Golang的后端服务，用于处理简历解析、用户信息总结和面试问题生成。该服务与前端面试问题生成器配合使用，提供API接口支持。

## 功能特点

- **PDF简历解析**：提取上传的PDF文件内容，并解析为结构化的简历信息
- **用户信息总结**：使用本地Ollama API对简历信息进行总结，提炼核心技能和经验
- **面试问题生成**：基于简历和岗位要求，使用Ollama API生成针对性的面试问题和参考答案
- **Server-Sent Events**：使用SSE协议实现面试问题的流式生成和推送

## 系统要求

- Go 1.16+
- 本地运行的Ollama服务（默认地址：http://localhost:11434）

## 安装步骤

1. 克隆仓库

```bash
git clone <repository-url>
cd demo-interview-web
```

2. 安装依赖

```bash
go mod download
```

3. 构建项目

```bash
go build -o interview-server
```

4. 运行服务

```bash
./interview-server
```

服务将在 http://localhost:8222 上启动。

## API接口

### 1. PDF内容提取

- **URL**: `/api/dev/v1/pdf/extract`
- **方法**: POST
- **Content-Type**: multipart/form-data
- **参数**: 
  - `file`: PDF文件
- **响应**: 结构化的简历信息（JSON格式）

### 2. 用户信息总结

- **URL**: `/api/dev/v1/userinfo/summary`
- **方法**: POST
- **Content-Type**: application/json
- **请求体**:
  ```json
  {
    "resume": {
      "basicInfo": { ... },
      "education": [ ... ],
      "workExperience": [ ... ],
      "skills": [ ... ],
      "projects": [ ... ],
      "certificates": [ ... ]
    }
  }
  ```
- **响应**:
  ```json
  {
    "summary": "总结内容..."
  }
  ```

### 3. 面试问题生成（SSE）

- **URL**: `/api/v1/interview/info/get_all`
- **方法**: POST
- **Content-Type**: application/json
- **请求体**:
  ```json
  {
    "resumeSummary": "简历总结...",
    "jobTitle": "职位名称",
    "companyName": "公司名称",
    "jobRequirements": "职位要求..."
  }
  ```
- **响应**: Server-Sent Events流，每个事件包含一个面试问题

## 配置说明

默认配置：
- 服务端口：8222
- Ollama API地址：http://localhost:11434/api/generate
- 默认模型：llama3

## 与前端集成

本服务设计为与面试问题生成器前端应用配合使用。前端应用可以：

1. 上传简历PDF文件到`/api/dev/v1/pdf/extract`接口
2. 将解析后的简历信息发送到`/api/dev/v1/userinfo/summary`接口获取总结
3. 将简历总结和岗位信息发送到`/api/v1/interview/info/get_all`接口，通过SSE接收生成的面试问题

## 开发说明

- 项目使用Gin框架构建HTTP服务
- 使用pdfcpu库进行PDF内容提取
- 使用Ollama API进行AI生成和总结
- 实现了SSE协议用于流式传输生成的面试问题

## 许可证

[MIT](LICENSE)