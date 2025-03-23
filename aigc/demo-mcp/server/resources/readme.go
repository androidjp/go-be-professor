package resources

import (
	"context"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	mcpsrv "github.com/mark3labs/mcp-go/server"
)

//所谓资源：定义如何导出数据给到 LLM
// 可以是静态资源：文件、接口响应体、数据库查询、系统信息等

// 添加静态资源 给 mcp server
func AddResourceReadMe(srv *mcpsrv.MCPServer) {
	resource := mcp.NewResource(
		"docs://readme",
		"Project README",
		mcp.WithResourceDescription("The project's README file"),
		mcp.WithMIMEType("text/markdown"),
	)

	// Add resource with its handler
	srv.AddResource(resource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		content, err := os.ReadFile("README.md")
		if err != nil {
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "docs://readme",
				MIMEType: "text/markdown",
				Text:     string(content),
			},
		}, nil
	})

}

func AddDynamicResource(srv *mcpsrv.MCPServer) {
	// // Dynamic resource example - user profiles by ID
	// template := mcp.NewResourceTemplate(
	// 	"users://{id}/profile",
	// 	"User Profile",
	// 	mcp.WithTemplateDescription("Returns user profile information"),
	// 	mcp.WithTemplateMIMEType("application/json"),
	// )

	// // Add template with its handler
	// srv.AddResourceTemplate(template, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// 	// Extract ID from the URI using regex matching
	// 	// The server automatically matches URIs to templates
	// 	userID := extractIDFromURI(request.Params.URI)

	// 	profile, err := getUserProfile(userID) // Your DB/API call here
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	return []mcp.ResourceContents{
	// 		mcp.TextResourceContents{
	// 			URI:      request.Params.URI,
	// 			MIMEType: "application/json",
	// 			Text:     profile,
	// 		},
	// 	}, nil
	// })

}
