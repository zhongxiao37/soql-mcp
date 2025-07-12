package tools

import (
	"context"
	"fmt"
	"soql-mcp/pkg"

	"github.com/mark3labs/mcp-go/mcp"
)

// CreateDescribeTool creates a new Salesforce object describe tool
func CreateDescribeTool() mcp.Tool {
	return mcp.NewTool("describe",
		mcp.WithDescription("Describe Salesforce objects to get their metadata, fields, and properties"),
		mcp.WithString("object",
			mcp.Required(),
			mcp.Description("The Salesforce object name to describe (e.g., Account, Contact, Opportunity)"),
		),
		mcp.WithString("format",
			mcp.Description("Output format: 'json' or 'table' (default: table)"),
		),
	)
}

// DescribeHandler handles Salesforce object describe requests
func DescribeHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	objectName, err := request.RequireString("object")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Object parameter is required: %v", err)), nil
	}

	format := request.GetString("format", "table")
	if format != "json" && format != "table" {
		format = "table"
	}

	// Load configuration
	config := pkg.LoadConfig()

	// Create Salesforce client
	sfClient := pkg.NewSalesforceClient(config)

	// Authenticate with Salesforce
	if err := sfClient.Authenticate(); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Authentication failed: %v", err)), nil
	}

	// Execute describe operation
	result, err := sfClient.Describe(objectName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Describe operation failed: %v", err)), nil
	}

	// Format and return results
	var output string
	if format == "json" {
		output = pkg.FormatDescribeAsJSON(result)
	} else {
		output = pkg.FormatDescribeAsTable(result)
	}

	return mcp.NewToolResultText(output), nil
}
