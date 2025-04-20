package config

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFlags(t *testing.T) {
	type want struct {
		host           string
		reportInterval time.Duration
		pollInterval   time.Duration
	}
	tests := []struct {
		args []string
		want want
	}{
		{
			args: []string{"-a", "localhost:8081", "-r", "20", "-p", "5"},
			want: want{
				host:           "localhost:8081",
				reportInterval: 20 * time.Second,
				pollInterval:   5 * time.Second,
			},
		},
		{
			args: []string{"-a", "localhost:8081"},
			want: want{
				host:           "localhost:8081",
				reportInterval: 10 * time.Second,
				pollInterval:   2 * time.Second,
			},
		},
		{
			args: []string{"-r", "20"},
			want: want{
				host:           "localhost:8080",
				reportInterval: 20 * time.Second,
				pollInterval:   2 * time.Second,
			},
		},
		{
			args: []string{"-p", "5"},
			want: want{
				host:           "localhost:8080",
				reportInterval: 10 * time.Second,
				pollInterval:   5 * time.Second,
			},
		},
	}

	for _, test := range tests {
		os.Args = append([]string{"cmd"}, test.args...)

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		config, _ := NewConfig()

		assert.Equal(t, test.want.host, config.HTTPAddress)
		assert.Equal(t, test.want.reportInterval, config.ReportInterval)
		assert.Equal(t, test.want.pollInterval, config.PollInterval)
	}
}

func TestEnvVariables(t *testing.T) {
	type want struct {
		host           string
		reportInterval time.Duration
		pollInterval   time.Duration
	}
	tests := []struct {
		args     []string
		envName  string
		envValue string
		want     want
	}{
		{
			args:     []string{"-a", "localhost:8081", "-r", "20", "-p", "5"},
			envName:  "ADDRESS",
			envValue: "localhost:8090",
			want: want{
				host:           "localhost:8090",
				reportInterval: 20 * time.Second,
				pollInterval:   5 * time.Second,
			},
		},
		{
			args:     []string{"-r", "20", "-p", "5"},
			envName:  "ADDRESS",
			envValue: "localhost:8090",
			want: want{
				host:           "localhost:8090",
				reportInterval: 20 * time.Second,
				pollInterval:   5 * time.Second,
			},
		},
		{
			args:     []string{"-r", "20", "-p", "5"},
			envName:  "SOME_VAR",
			envValue: "SOME_VAL",
			want: want{
				host:           "localhost:8080",
				reportInterval: 20 * time.Second,
				pollInterval:   5 * time.Second,
			},
		},
		{
			args:     []string{"-r", "20"},
			envName:  "REPORT_INTERVAL",
			envValue: "30",
			want: want{
				host:           "localhost:8080",
				reportInterval: 30 * time.Second,
				pollInterval:   2 * time.Second,
			},
		},
		{
			args:     []string{},
			envName:  "REPORT_INTERVAL",
			envValue: "30",
			want: want{
				host:           "localhost:8080",
				reportInterval: 30 * time.Second,
				pollInterval:   2 * time.Second,
			},
		},
		{
			args:     []string{},
			envName:  "SOME_VAR",
			envValue: "SOME_VAL",
			want: want{
				host:           "localhost:8080",
				reportInterval: 10 * time.Second,
				pollInterval:   2 * time.Second,
			},
		},
		{
			args:     []string{"-p", "5"},
			envName:  "POLL_INTERVAL",
			envValue: "8",
			want: want{
				host:           "localhost:8080",
				reportInterval: 10 * time.Second,
				pollInterval:   8 * time.Second,
			},
		},
		{
			args:     []string{},
			envName:  "POLL_INTERVAL",
			envValue: "8",
			want: want{
				host:           "localhost:8080",
				reportInterval: 10 * time.Second,
				pollInterval:   8 * time.Second,
			},
		},
		{
			args:     []string{},
			envName:  "SOME_VAL",
			envValue: "SOME_VAR",
			want: want{
				host:           "localhost:8080",
				reportInterval: 10 * time.Second,
				pollInterval:   2 * time.Second,
			},
		},
	}

	for _, test := range tests {
		os.Setenv(test.envName, test.envValue)
		os.Args = append([]string{"cmd"}, test.args...)

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		config, _ := NewConfig()
		os.Unsetenv(test.envName)

		assert.Equal(t, test.want.host, config.HTTPAddress)
		assert.Equal(t, test.want.reportInterval, config.ReportInterval)
		assert.Equal(t, test.want.pollInterval, config.PollInterval)
	}
}

func TestParseKeyFlag(t *testing.T) {
	type want struct {
		key string
	}
	tests := []struct {
		args     []string
		envName  string
		envValue string
		want     want
	}{
		{
			args:     []string{"-k", "super_secret_key"},
			envName:  "KEY",
			envValue: "not_so_secret",
			want: want{
				key: "not_so_secret",
			},
		},
		{
			args:     []string{},
			envName:  "KEY",
			envValue: "secret_key",
			want: want{
				key: "secret_key",
			},
		},
		{
			args:     []string{""},
			envName:  "SOME_VAR",
			envValue: "SOME_VAL",
			want: want{
				key: "",
			},
		},
	}

	for _, test := range tests {
		os.Setenv(test.envName, test.envValue)
		os.Args = append([]string{"cmd"}, test.args...)

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		config, _ := NewConfig()
		os.Clearenv()

		assert.Equal(t, test.want.key, config.Key)
	}
}

func TestParseRateLimitFlag(t *testing.T) {
	type want struct {
		rateLimit int
	}
	tests := []struct {
		args     []string
		envName  string
		envValue string
		want     want
	}{
		{
			args:     []string{"-l", "10"},
			envName:  "RATE_LIMIT",
			envValue: "8",
			want: want{
				rateLimit: 8,
			},
		},
		{
			args:     []string{"-l", "10"},
			envName:  "SOME_VAR",
			envValue: "25",
			want: want{
				rateLimit: 10,
			},
		},
		{
			args:     []string{},
			envName:  "RATE_LIMIT",
			envValue: "9",
			want: want{
				rateLimit: 9,
			},
		},
		{
			args:     []string{""},
			envName:  "SOME_VAR",
			envValue: "SOME_VAL",
			want: want{
				rateLimit: 5,
			},
		},
	}

	for _, test := range tests {
		os.Setenv(test.envName, test.envValue)
		os.Args = append([]string{"cmd"}, test.args...)

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		config, _ := NewConfig()
		os.Clearenv()

		assert.Equal(t, test.want.rateLimit, config.RateLimit)
	}
}

func TestParseCryptoKeyFlag(t *testing.T) {
	pubKey := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAu1SU1LfVLPHCozMxH2Mo
4lgOEePzNm0tRgeLezV6ffAt0gunVTLw7onLRnrq0/IzW7yWR7QkrmBL7jTKEn5u
+qKhbwKfBstIs+bMY2Zkp18gnTxKLxoS2tFczGkPLPgizskuemMghRniWaoLcyeh
kd3qqGElvW/VDL5AaWTg0nLVkjRo9z+40RQzuVaE8AkAFmxZzow3x+VJkAE/Ag+Z
cL5HBPpE5oVuAfQwF1/7+9VP3Mp9v6sED6bFiPQ0NdwCYp6j6X7WQ8CJ7M5kQ+7J
9Z6MCQD5qjU1fXg9JwZw5V5Z0X6J+ZQ0C3c0yW0q5fYDP6wUcJb6MnN4B7pXwJ2d
zQIDAQAB
-----END PUBLIC KEY-----`

	tmpFile, err := os.CreateTemp("", "public_key_*.pem")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(pubKey); err != nil {
		t.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	type want struct {
		cryptoKeyExists bool
	}
	tests := []struct {
		name     string
		args     []string
		envName  string
		envValue string
		want     want
		wantErr  bool
	}{
		{
			name:    "no key specified",
			args:    []string{},
			want:    want{cryptoKeyExists: false},
			wantErr: false,
		},
		{
			name:    "key from flag",
			args:    []string{"-crypto-key", tmpFile.Name()},
			want:    want{cryptoKeyExists: true},
			wantErr: false,
		},
		{
			name:     "key from env overrides flag",
			args:     []string{"-crypto-key", "nonexistent.pem"},
			envName:  "CRYPTO_KEY",
			envValue: tmpFile.Name(),
			want:     want{cryptoKeyExists: true},
			wantErr:  false,
		},
		{
			name:    "invalid key path",
			args:    []string{"-crypto-key", "nonexistent.pem"},
			want:    want{cryptoKeyExists: false},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.envName != "" {
				os.Setenv(test.envName, test.envValue)
				defer os.Unsetenv(test.envName)
			}

			os.Args = append([]string{"cmd"}, test.args...)
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

			config, err := NewConfig()
			if (err != nil) != test.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if test.wantErr {
				return
			}

			if test.want.cryptoKeyExists {
				if config.CryptoKey == nil {
					t.Error("Expected crypto key to be loaded, but got nil")
				}
			} else {
				if config.CryptoKey != nil {
					t.Error("Expected no crypto key, but got one")
				}
			}
		})
	}
}

func TestLoadPublicKey(t *testing.T) {
	validPubKey := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAu1SU1LfVLPHCozMxH2Mo
4lgOEePzNm0tRgeLezV6ffAt0gunVTLw7onLRnrq0/IzW7yWR7QkrmBL7jTKEn5u
+qKhbwKfBstIs+bMY2Zkp18gnTxKLxoS2tFczGkPLPgizskuemMghRniWaoLcyeh
kd3qqGElvW/VDL5AaWTg0nLVkjRo9z+40RQzuVaE8AkAFmxZzow3x+VJkAE/Ag+Z
cL5HBPpE5oVuAfQwF1/7+9VP3Mp9v6sED6bFiPQ0NdwCYp6j6X7WQ8CJ7M5kQ+7J
9Z6MCQD5qjU1fXg9JwZw5V5Z0X6J+ZQ0C3c0yW0q5fYDP6wUcJb6MnN4B7pXwJ2d
zQIDAQAB
-----END PUBLIC KEY-----`

	validTmpFile, err := os.CreateTemp("", "valid_public_key_*.pem")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(validTmpFile.Name())

	if _, WriteErr := validTmpFile.WriteString(validPubKey); err != nil {
		t.Fatal(WriteErr)
	}
	if CloseErr := validTmpFile.Close(); err != nil {
		t.Fatal(CloseErr)
	}

	invalidPubKey := `-----BEGIN PUBLIC KEY-----
INVALID KEY DATA
-----END PUBLIC KEY-----`

	invalidTmpFile, err := os.CreateTemp("", "invalid_public_key_*.pem")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(invalidTmpFile.Name())

	if _, err := invalidTmpFile.WriteString(invalidPubKey); err != nil {
		t.Fatal(err)
	}
	if err := invalidTmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		keyPath string
		wantKey bool
		wantErr bool
	}{
		{
			name:    "empty path",
			keyPath: "",
			wantKey: false,
			wantErr: false,
		},
		{
			name:    "valid key",
			keyPath: validTmpFile.Name(),
			wantKey: true,
			wantErr: false,
		},
		{
			name:    "invalid key format",
			keyPath: invalidTmpFile.Name(),
			wantKey: false,
			wantErr: true,
		},
		{
			name:    "nonexistent file",
			keyPath: "nonexistent.pem",
			wantKey: false,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := loadPublicKey(test.keyPath)
			if (err != nil) != test.wantErr {
				t.Errorf("loadPublicKey() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if (got != nil) != test.wantKey {
				t.Errorf("loadPublicKey() got key = %v, want key %v", got != nil, test.wantKey)
			}
		})
	}
}

func TestAgentConfigPriority(t *testing.T) {
	tmpDir := t.TempDir()

	configContent := `{
		"address": "json:8080",
		"report_interval": 5,
		"poll_interval": 3,
		"rate_limit": 7,
		"key": "json_key"
	}`

	configPath := filepath.Join(tmpDir, "config.json")
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	tests := []struct {
		name       string
		args       []string
		envVars    map[string]string
		configFile string
		expected   Config
	}{
		{
			name:       "json config only",
			configFile: configPath,
			expected: Config{
				HTTPAddress:    "json:8080",
				ReportInterval: 5 * time.Second,
				PollInterval:   3 * time.Second,
				RateLimit:      7,
				Key:            "json_key",
			},
		},
		{
			name:       "flag overrides json",
			configFile: configPath,
			args:       []string{"-a", "flag:8080", "-r", "3", "-k", "flag_key"},
			expected: Config{
				HTTPAddress:    "flag:8080",
				ReportInterval: 3 * time.Second,
				PollInterval:   3 * time.Second,
				RateLimit:      7,
				Key:            "flag_key",
			},
		},
		{
			name:       "env override all",
			configFile: configPath,
			envVars: map[string]string{
				"ADDRESS":         "env:8080",
				"REPORT_INTERVAL": "10",
				"KEY":             "env_key",
			},
			args: []string{"-a", "flag:8080", "-r", "3", "-k", "flag_key"},
			expected: Config{
				HTTPAddress:    "env:8080",
				ReportInterval: 10 * time.Second,
				PollInterval:   3 * time.Second,
				RateLimit:      7,
				Key:            "env_key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			}()

			args := append([]string{"-config", tt.configFile}, tt.args...)
			os.Args = append([]string{"cmd"}, args...)
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

			cfg, err := NewConfig()
			require.NoError(t, err)

			assert.Equal(t, tt.expected.HTTPAddress, cfg.HTTPAddress)
			assert.Equal(t, tt.expected.ReportInterval, cfg.ReportInterval)
			assert.Equal(t, tt.expected.PollInterval, cfg.PollInterval)
			assert.Equal(t, tt.expected.RateLimit, cfg.RateLimit)
			assert.Equal(t, tt.expected.Key, cfg.Key)
		})
	}
}
