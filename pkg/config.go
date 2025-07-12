package pkg

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration values
type Config struct {
	ServerName    string
	ServerVersion string
	ResourcePath  string
	Debug         bool
	LogLevel      string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		ServerName:    GetEnvWithDefault("MCP_SERVER_NAME", "Demo 🚀"),
		ServerVersion: GetEnvWithDefault("MCP_SERVER_VERSION", "1.0.0"),
		ResourcePath:  GetEnvWithDefault("MCP_RESOURCE_PATH", "/Users/pzhong/Documents/github/soql-mcp/terms.json"),
		Debug:         getEnvBool("MCP_DEBUG", false),
		LogLevel:      GetEnvWithDefault("MCP_LOG_LEVEL", "info"),
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
	fmt.Printf("  Resource Path: %s\n", c.ResourcePath)
	fmt.Printf("  Debug: %t\n", c.Debug)
	fmt.Printf("  Log Level: %s\n", c.LogLevel)
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
