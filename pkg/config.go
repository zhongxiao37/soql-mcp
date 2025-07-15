package pkg

import (
	"fmt"
	"os"
	"strconv"
)

// Build-time variables set via ldflags
var (
	Version   = "" // Default version for development
	Commit    = "" // Default commit hash
	BuildDate = "" // Default build time
)

// Config holds all configuration values
type Config struct {
	ServerName    string
	ServerVersion string
	Commit        string
	BuildDate     string
	ResourcePath  string
	Debug         bool
	LogLevel      string
	// Salesforce configuration
	SalesforceURL           string
	SalesforceClientID      string
	SalesforceClientSecret  string
	SalesforceUsername      string
	SalesforcePassword      string
	SalesforceSecurityToken string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		ServerName:    GetEnvWithDefault("MCP_SERVER_NAME", "Demo 🚀"),
		ServerVersion: GetEnvWithDefault("MCP_SERVER_VERSION", Version),
		Commit:        GetEnvWithDefault("MCP_COMMIT", Commit),
		BuildDate:     GetEnvWithDefault("MCP_BUILD_DATE", BuildDate),
		ResourcePath:  GetEnvWithDefault("MCP_RESOURCE_PATH", ""),
		Debug:         getEnvBool("MCP_DEBUG", false),
		LogLevel:      GetEnvWithDefault("MCP_LOG_LEVEL", "info"),
		// Salesforce configuration
		SalesforceURL:           GetEnvWithDefault("SALESFORCE_URL", "https://login.salesforce.com"),
		SalesforceClientID:      GetEnvWithDefault("SALESFORCE_CLIENT_ID", ""),
		SalesforceClientSecret:  GetEnvWithDefault("SALESFORCE_CLIENT_SECRET", ""),
		SalesforceUsername:      GetEnvWithDefault("SALESFORCE_USERNAME", ""),
		SalesforcePassword:      GetEnvWithDefault("SALESFORCE_PASSWORD", ""),
		SalesforceSecurityToken: GetEnvWithDefault("SALESFORCE_SECURITY_TOKEN", ""),
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		fmt.Printf("Configuration error: %v\n", err)
		os.Exit(1)
	}

	return config
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.ServerName == "" {
		return fmt.Errorf("server name cannot be empty")
	}
	if c.ServerVersion == "" {
		return fmt.Errorf("server version cannot be empty")
	}
	if c.ResourcePath == "" {
		return fmt.Errorf("resource path cannot be empty")
	}
	return nil
}

// Print outputs the configuration for debugging
func (c *Config) Print() {
	fmt.Printf("Configuration:\n")
	fmt.Printf("  Server Name: %s\n", c.ServerName)
	fmt.Printf("  Server Version: %s\n", c.ServerVersion)
	fmt.Printf("  Commit: %s\n", c.Commit)
	fmt.Printf("  Resource Path: %s\n", c.ResourcePath)
	fmt.Printf("  Debug: %t\n", c.Debug)
	fmt.Printf("  Log Level: %s\n", c.LogLevel)
	fmt.Printf("  Salesforce URL: %s\n", c.SalesforceURL)
	fmt.Printf("  Salesforce Client ID: %s\n", c.SalesforceClientID)
	fmt.Printf("  Salesforce Username: %s\n", c.SalesforceUsername)
}

// GetEnvWithDefault returns the value of an environment variable or a default value if not set
func GetEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool returns a boolean environment variable value
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
