// Package fileconfig provides common functionality for reading configuration from files.
package fileconfig

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// CommonConfig contains fields shared between agent and server configurations
type CommonConfig struct {
	Address   string `json:"address"`
	CryptoKey string `json:"crypto_key"`
	Key       string `json:"key"`
	Scheme    string `json:"scheme"`
}

// AgentConfig represents agent-specific configuration from file
type AgentConfig struct {
	CommonConfig
	ReportIntervalSec int `json:"report_interval"`
	PollIntervalSec   int `json:"poll_interval"`
	RateLimit         int `json:"rate_limit"`
}

// ServerConfig represents server-specific configuration from file
type ServerConfig struct {
	CommonConfig
	Restore          bool   `json:"restore"`
	StoreIntervalSec int    `json:"store_interval"`
	StoreFile        string `json:"store_file"`
	DatabaseDSN      string `json:"database_dsn"`
}

// ReadAgentConfig reads agent configuration from JSON file
func ReadAgentConfig(path string) (*AgentConfig, error) {
	if path == "" {
		return nil, nil
	}

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	var cfg AgentConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// ReadServerConfig reads server configuration from JSON file
func ReadServerConfig(path string) (*ServerConfig, error) {
	if path == "" {
		return nil, nil
	}

	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	var cfg ServerConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
