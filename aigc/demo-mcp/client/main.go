package main

import (
	"github.com/mark3labs/mcp-go/client"
)

func main() {

}

func connectToMCPServer() (*client.MCPClient, error) {

	// stdio 客户端
	// 适用于通过标准输入/输出与 MCP Server 交互的场景，无需手动启动异步通信
	// cli, err := client.NewStdioMCPClient()

	// 适用于需要通过 Server-Sent Events（SSE）进行异步通信的场景，需手动启动通信
	// cli, err := client.NewSSEMCPClient(myBaseURL)
	// err = cli.Start(ctx)

	return nil, nil
}
