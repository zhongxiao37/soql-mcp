package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SalesforceAuth represents OAuth response from Salesforce
type SalesforceAuth struct {
	AccessToken string `json:"access_token"`
	InstanceURL string `json:"instance_url"`
	ID          string `json:"id"`
	TokenType   string `json:"token_type"`
	IssuedAt    string `json:"issued_at"`
	Signature   string `json:"signature"`
}

// SalesforceQueryResponse represents the response from SOQL query
type SalesforceQueryResponse struct {
	TotalSize int           `json:"totalSize"`
	Done      bool          `json:"done"`
	Records   []interface{} `json:"records"`
}

// SalesforceError represents error response from Salesforce
type SalesforceError struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}

// SalesforceErrorResponse represents error response wrapper
type SalesforceErrorResponse struct {
	Error            string            `json:"error"`
	ErrorDescription string            `json:"error_description"`
	Errors           []SalesforceError `json:"errors"`
}

// SalesforceClient handles Salesforce API operations
type SalesforceClient struct {
	config *Config
	auth   *SalesforceAuth
}

// NewSalesforceClient creates a new Salesforce client
func NewSalesforceClient(config *Config) *SalesforceClient {
	return &SalesforceClient{
		config: config,
	}
}

// ValidateConfig checks if Salesforce configuration is complete
func (sf *SalesforceClient) ValidateConfig() error {
	if sf.config.SalesforceClientID == "" {
		return fmt.Errorf("SALESFORCE_CLIENT_ID is required")
	}
	if sf.config.SalesforceClientSecret == "" {
		return fmt.Errorf("SALESFORCE_CLIENT_SECRET is required")
	}
	if sf.config.SalesforceUsername == "" {
		return fmt.Errorf("SALESFORCE_USERNAME is required")
	}
	if sf.config.SalesforcePassword == "" {
		return fmt.Errorf("SALESFORCE_PASSWORD is required")
	}
	return nil
}

// Authenticate performs OAuth authentication with Salesforce
func (sf *SalesforceClient) Authenticate() error {
	if err := sf.ValidateConfig(); err != nil {
		return err
	}

	// Prepare authentication request
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", sf.config.SalesforceClientID)
	data.Set("client_secret", sf.config.SalesforceClientSecret)
	data.Set("username", sf.config.SalesforceUsername)
	data.Set("password", sf.config.SalesforcePassword+sf.config.SalesforceSecurityToken)

	// Make authentication request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.PostForm(sf.config.SalesforceURL+"/services/oauth2/token", data)
	if err != nil {
		return fmt.Errorf("failed to make authentication request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read authentication response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp SalesforceErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return fmt.Errorf("authentication failed: %s - %s", errorResp.Error, errorResp.ErrorDescription)
		}
		return fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode, string(body))
	}

	var auth SalesforceAuth
	if err := json.Unmarshal(body, &auth); err != nil {
		return fmt.Errorf("failed to parse authentication response: %v", err)
	}

	sf.auth = &auth
	return nil
}

// Query executes a SOQL query against Salesforce
func (sf *SalesforceClient) Query(query string) (*SalesforceQueryResponse, error) {
	if sf.auth == nil {
		return nil, fmt.Errorf("not authenticated, call Authenticate() first")
	}

	// Prepare query URL
	queryURL := fmt.Sprintf("%s/services/data/v57.0/query", sf.auth.InstanceURL)

	// URL encode the query
	params := url.Values{}
	params.Add("q", query)
	fullURL := fmt.Sprintf("%s?%s", queryURL, params.Encode())

	// Create HTTP request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sf.auth.AccessToken))
	req.Header.Set("Content-Type", "application/json")

	// Make request
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read query response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp SalesforceErrorResponse
		if err := json.Unmarshal(body, &errorResp); err == nil && len(errorResp.Errors) > 0 {
			return nil, fmt.Errorf("query failed: %s - %s", errorResp.Errors[0].ErrorCode, errorResp.Errors[0].Message)
		}
		return nil, fmt.Errorf("query failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result SalesforceQueryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse query response: %v", err)
	}

	return &result, nil
}

// FormatAsTable formats query results as a simple table
func FormatAsTable(result *SalesforceQueryResponse) string {
	if result.TotalSize == 0 {
		return "No records found."
	}

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Total Records: %d\n", result.TotalSize))
	buffer.WriteString(fmt.Sprintf("Records Returned: %d\n", len(result.Records)))
	buffer.WriteString(strings.Repeat("-", 50) + "\n")

	for i, record := range result.Records {
		buffer.WriteString(fmt.Sprintf("Record %d:\n", i+1))
		if recordMap, ok := record.(map[string]interface{}); ok {
			for key, value := range recordMap {
				if key != "attributes" { // Skip Salesforce metadata
					buffer.WriteString(fmt.Sprintf("  %s: %v\n", key, value))
				}
			}
		}
		buffer.WriteString(strings.Repeat("-", 30) + "\n")
	}

	return buffer.String()
}

// FormatAsJSON formats query results as JSON
func FormatAsJSON(result *SalesforceQueryResponse) string {
	jsonBytes, _ := json.MarshalIndent(result, "", "  ")
	return string(jsonBytes)
}
