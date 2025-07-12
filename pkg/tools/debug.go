package tools

import (
	"context"
	"fmt"
	"soql-mcp/pkg"

	"github.com/mark3labs/mcp-go/mcp"
)

// CreateDebugTool creates a debug tool for returning server configuration information
func CreateDebugTool() mcp.Tool {
	return mcp.NewTool("debug",
		mcp.WithDescription("return server configuration information"),
	)
}

// DebugHandler handles debug tool requests and returns configuration information
func DebugHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Load current config
	config := pkg.LoadConfig()

	// Format config information as a string
	configInfo := fmt.Sprintf("Server configuration information:\n")
	configInfo += fmt.Sprintf("  Server name: %s\n", config.ServerName)
	configInfo += fmt.Sprintf("  Server version: %s\n", config.ServerVersion)
	configInfo += fmt.Sprintf("  Resource path: %s\n", config.ResourcePath)
	configInfo += fmt.Sprintf("  Debug mode: %t\n", config.Debug)
	configInfo += fmt.Sprintf("  Log level: %s\n", config.LogLevel)

	return mcp.NewToolResultText(configInfo), nil
}
