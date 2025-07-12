package resources

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
)

// CreateTermsResource creates a new terms resource
func CreateTermsResource(resourcePath string) mcp.Resource {
	return mcp.NewResource(
		fmt.Sprintf("file://%s", resourcePath),
		"terms",
		mcp.WithResourceDescription("Terms"),
		mcp.WithMIMEType("application/json"),
	)
}

// TermsResourceHandler handles the terms resource read requests
func TermsResourceHandler(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
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
