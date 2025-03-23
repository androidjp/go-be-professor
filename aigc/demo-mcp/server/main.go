package main

import (
	"fmt"
	"mcpsrvdemo/resources"
	"mcpsrvdemo/tools"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"JP MCP Srv Demo",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	// Add readme resource
	resources.AddResourceReadMe(s)

	// Add the calculator tool
	tools.AddToolCalculator(s)

	// Add the HTTP request tool
	tools.AddToolHTTPRequest(s)

	// Add current_time
	tools.AddToolCurrentTime(s)

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	} else {
		fmt.Println("Server stopped")
	}
}
