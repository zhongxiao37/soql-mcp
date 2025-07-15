package main

import (
	"fmt"
	"os"

	"soql-mcp/pkg"
	"soql-mcp/pkg/resources"
	"soql-mcp/pkg/tools"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

var (
	versionFlag bool
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "soql-mcp",
		Short: "SOQL MCP Server",
		Long:  "SOQL MCP Server for Salesforce object queries",
		Run: func(cmd *cobra.Command, args []string) {
			if versionFlag {
				fmt.Printf("soql-mcp version %s (commit: %s, build date: %s)\n", pkg.Version, pkg.Commit, pkg.BuildDate)
				os.Exit(0)
			}
			runServer()
		},
	}

	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Print version information")

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func runServer() {
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

	// Add tools
	s.AddTool(tools.CreateDebugTool(), tools.DebugHandler)
	s.AddTool(tools.CreateQueryTool(), tools.QueryHandler)
	s.AddTool(tools.CreateDescribeTool(), tools.DescribeHandler)

	// Add terms resource using the new resources package
	s.AddResource(resources.CreateTermsResource(config.ResourcePath), resources.TermsResourceHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
