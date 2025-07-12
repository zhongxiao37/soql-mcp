package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"soql-mcp/pkg"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Load configuration from environment variables
	config := pkg.LoadConfig()

	// Print configuration for debugging
	if config.Debug {
		config.Print()
	}

	// Create a new MCP server with resources capability
	s := server.NewMCPServer(
		config.ServerName,
		config.ServerVersion,
		server.WithResourceCapabilities(true, true),
		server.WithToolCapabilities(true),
	)

	// Add existing tool
	tool := mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)

	// Add tool handler
	s.AddTool(tool, helloHandler)

	// Add debug tool to return config information
	debugTool := mcp.NewTool("debug",
		mcp.WithDescription("返回服务器配置信息"),
	)

	// Add debug tool handler
	s.AddTool(debugTool, debugHandler)

	// Add a static resource using environment variable
	resource := mcp.NewResource(
		fmt.Sprintf("file://%s", config.ResourcePath),
		"terms",
		mcp.WithResourceDescription("Terms"),
		mcp.WithMIMEType("application/json"),
	)
	s.AddResource(resource, fileResourceHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}

func debugHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Load current config
	config := pkg.LoadConfig()

	// Format config information as a string
	configInfo := fmt.Sprintf("服务器配置信息:\n")
	configInfo += fmt.Sprintf("  服务器名称: %s\n", config.ServerName)
	configInfo += fmt.Sprintf("  服务器版本: %s\n", config.ServerVersion)
	configInfo += fmt.Sprintf("  资源路径: %s\n", config.ResourcePath)
	configInfo += fmt.Sprintf("  调试模式: %t\n", config.Debug)
	configInfo += fmt.Sprintf("  日志级别: %s\n", config.LogLevel)

	return mcp.NewToolResultText(configInfo), nil
}

// Handler for static terms resource
func fileResourceHandler(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	// Parse the URI to extract the file path
	parsedURI, err := url.Parse(request.Params.URI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URI: %w", err)
	}

	// Extract the file path from the URI
	filePath := parsedURI.Path
	fmt.Println("filePath", filePath)

	// Read the file using the extracted path
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(content),
		},
	}, nil
}
