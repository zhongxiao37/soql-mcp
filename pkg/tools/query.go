package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/zhongxiao37/soql-mcp/pkg"
)

// CreateQueryTool creates a new SOQL query tool
func CreateQueryTool() mcp.Tool {
	return mcp.NewTool("query",
		mcp.WithDescription("Execute SOQL queries against Salesforce"),
		mcp.WithString("soql",
			mcp.Required(),
			mcp.Description("The SOQL query to execute (e.g., SELECT Id, Name FROM Account LIMIT 10)"),
		),
		mcp.WithString("format",
			mcp.Description("Output format: 'json' or 'table' (default: json)"),
		),
	)
}

// QueryHandler handles SOQL query requests
func QueryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	soql, err := request.RequireString("soql")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Query parameter is required: %v", err)), nil
	}

	format := request.GetString("format", "json")
	if format != "json" && format != "table" {
		format = "json"
	}

	// Load configuration
	config := pkg.LoadConfig()

	// Get authenticated Salesforce client (with connection reuse)
	clientManager := pkg.GetClientManager(config)
	sfClient, err := clientManager.GetClient()
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Authentication failed: %v", err)), nil
	}

	// Execute SOQL query
	result, err := sfClient.Query(soql)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Query execution failed: %v", err)), nil
	}

	// Format and return results
	var output string
	if format == "table" {
		output = pkg.FormatAsTable(result)
	} else {
		output = pkg.FormatAsJSON(result)
	}

	return mcp.NewToolResultText(output), nil
}
