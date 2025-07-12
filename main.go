package main

import (
	"fmt"

	"soql-mcp/pkg"
	"soql-mcp/pkg/resources"
	"soql-mcp/pkg/tools"

	"github.com/mark3labs/mcp-go/server"
)

var (
	// 这些变量会在构建时通过 ldflags 注入
	version = "dev"
	commit  = "unknown"
)

func main() {
	// Load configuration from environment variables
	config := pkg.LoadConfig()

	// Print configuration for debugging
	if config.Debug {
		fmt.Printf("SOQL MCP Server version: %s (commit: %s)\n", version, commit)
		config.Print()
	}

	// Create a new MCP server with resources capability
	s := server.NewMCPServer(
		config.ServerName,
		config.ServerVersion,
		server.WithResourceCapabilities(true, true),
		server.WithToolCapabilities(true),
	)

	// Add tools
	s.AddTool(tools.CreateHelloTool(), tools.HelloHandler)
	s.AddTool(tools.CreateDebugTool(), tools.DebugHandler)
	s.AddTool(tools.CreateQueryTool(), tools.QueryHandler)

	// Add terms resource using the new resources package
	s.AddResource(resources.CreateTermsResource(config.ResourcePath), resources.TermsResourceHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
