package pkg

import (
	"sync"
	"time"
)

// ClientManager manages a singleton Salesforce client with connection reuse
type ClientManager struct {
	client      *SalesforceClient
	config      *Config
	lastAuth    time.Time
	tokenExpiry time.Duration
	mutex       sync.RWMutex
}

var (
	instance *ClientManager
	once     sync.Once
)

// GetClientManager returns the singleton instance of ClientManager
func GetClientManager(config *Config) *ClientManager {
	once.Do(func() {
		instance = &ClientManager{
			config:      config,
			tokenExpiry: 2 * time.Hour, // Salesforce tokens typically expire in 2 hours
		}
	})
	return instance
}

// GetClient returns an authenticated Salesforce client, reusing connection when possible
func (cm *ClientManager) GetClient() (*SalesforceClient, error) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	// Check if we need to authenticate or re-authenticate
	if cm.client == nil || cm.needsReauth() {
		if cm.client == nil {
			cm.client = NewSalesforceClient(cm.config)
		}

		if err := cm.client.Authenticate(); err != nil {
			return nil, err
		}

		cm.lastAuth = time.Now()
	}

	return cm.client, nil
}

// needsReauth checks if re-authentication is needed based on token expiry
func (cm *ClientManager) needsReauth() bool {
	// Re-authenticate if it's been more than 90% of token expiry time
	// This provides a buffer to avoid token expiry during requests
	return time.Since(cm.lastAuth) > time.Duration(float64(cm.tokenExpiry)*0.9)
}

// Reset clears the cached client (useful for testing or configuration changes)
func (cm *ClientManager) Reset() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.client = nil
	cm.lastAuth = time.Time{}
}
