package fileconfig

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadAgentConfig(t *testing.T) {
	tmpDir := t.TempDir()

	validConfig := `{
		"address": "localhost:9090",
		"report_interval": 5,
		"poll_interval": 1,
		"crypto_key": "/path/to/key.pem",
		"key": "test_key",
		"rate_limit": 10
	}`

	validPath := filepath.Join(tmpDir, "valid_config.json")
	err := os.WriteFile(validPath, []byte(validConfig), 0600)
	require.NoError(t, err)

	invalidConfig := `{ invalid json }`
	invalidPath := filepath.Join(tmpDir, "invalid_config.json")
	err = os.WriteFile(invalidPath, []byte(invalidConfig), 0600)
	require.NoError(t, err)

	tests := []struct {
		name    string
		path    string
		wantCfg *AgentConfig
		wantErr bool
	}{
		{
			name:    "valid config",
			path:    validPath,
			wantErr: false,
			wantCfg: &AgentConfig{
				CommonConfig: CommonConfig{
					Address:   "localhost:9090",
					CryptoKey: "/path/to/key.pem",
					Key:       "test_key",
				},
				ReportIntervalSec: 5,
				PollIntervalSec:   1,
				RateLimit:         10,
			},
		},
		{
			name:    "invalid json",
			path:    invalidPath,
			wantErr: true,
		},
		{
			name:    "non-existent file",
			path:    filepath.Join(tmpDir, "nonexistent.json"),
			wantErr: true,
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: false,
			wantCfg: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := ReadAgentConfig(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantCfg, cfg)
		})
	}
}

func TestReadServerConfig(t *testing.T) {
	tmpDir := t.TempDir()

	validConfig := `{
		"address": "localhost:9090",
		"store_interval": 30,
		"store_file": "/path/to/store.db",
		"restore": true,
		"database_dsn": "postgres://user:pass@localhost:5432/db",
		"crypto_key": "/path/to/key.pem",
		"key": "test_key"
	}`

	validPath := filepath.Join(tmpDir, "valid_config.json")
	err := os.WriteFile(validPath, []byte(validConfig), 0600)
	require.NoError(t, err)

	tests := []struct {
		name    string
		path    string
		wantCfg *ServerConfig
		wantErr bool
	}{
		{
			name:    "valid config",
			path:    validPath,
			wantErr: false,
			wantCfg: &ServerConfig{
				CommonConfig: CommonConfig{
					Address:   "localhost:9090",
					CryptoKey: "/path/to/key.pem",
					Key:       "test_key",
				},
				StoreIntervalSec: 30,
				StoreFile:        "/path/to/store.db",
				Restore:          true,
				DatabaseDSN:      "postgres://user:pass@localhost:5432/db",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := ReadServerConfig(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantCfg, cfg)
		})
	}
}
