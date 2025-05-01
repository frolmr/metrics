package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseHostFlags(t *testing.T) {
	type want struct {
		host string
	}
	tests := []struct {
		args []string
		want want
	}{
		{
			args: []string{"-a", "localhost:8081"},
			want: want{
				host: "localhost:8081",
			},
		},
		{
			args: []string{""},
			want: want{
				host: "localhost:8080",
			},
		},
	}

	for _, test := range tests {
		os.Args = append([]string{"cmd"}, test.args...)

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		config, _ := NewConfig()

		assert.Equal(t, test.want.host, config.HTTPAddress)
	}
}

func TestParseIntervalFlags(t *testing.T) {
	type want struct {
		interval time.Duration
	}
	tests := []struct {
		args []string
		want want
	}{
		{
			args: []string{"-i", "200"},
			want: want{
				interval: time.Duration(200) * time.Second,
			},
		},
		{
			args: []string{""},
			want: want{
				interval: time.Duration(300) * time.Second,
			},
		},
	}

	for _, test := range tests {
		os.Args = append([]string{"cmd"}, test.args...)

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		config, _ := NewConfig()

		assert.Equal(t, test.want.interval, config.StoreInterval)
	}
}

func TestParseFileFlags(t *testing.T) {
	type want struct {
		file string
	}
	tests := []struct {
		args []string
		want want
	}{
		{
			args: []string{"-f", "./db/snapshot"},
			want: want{
				file: "./db/snapshot",
			},
		},
		{
			args: []string{""},
			want: want{
				file: "data_snapshot",
			},
		},
	}

	for _, test := range tests {
		os.Args = append([]string{"cmd"}, test.args...)

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		config, _ := NewConfig()

		assert.Equal(t, test.want.file, config.FileStoragePath)
	}
}

func TestParseRestoreFlags(t *testing.T) {
	type want struct {
		restore bool
	}
	tests := []struct {
		args []string
		want want
	}{
		{
			args: []string{"-r", "true"},
			want: want{
				restore: true,
			},
		},
		{
			args: []string{""},
			want: want{
				restore: false,
			},
		},
	}

	for _, test := range tests {
		os.Args = append([]string{"cmd"}, test.args...)

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		config, _ := NewConfig()

		assert.Equal(t, test.want.restore, config.Restore)
	}
}

func TestHostEnvVariables(t *testing.T) {
	type want struct {
		host string
	}
	tests := []struct {
		args     []string
		envName  string
		envValue string
		want     want
	}{
		{
			args:     []string{"-a", "localhost:8081"},
			envName:  "ADDRESS",
			envValue: "localhost:8090",
			want: want{
				host: "localhost:8090",
			},
		},
		{
			args:     []string{},
			envName:  "ADDRESS",
			envValue: "localhost:8090",
			want: want{
				host: "localhost:8090",
			},
		},
		{
			args:     []string{},
			envName:  "SOME_VAR",
			envValue: "SOME_VAL",
			want: want{
				host: "localhost:8080",
			},
		},
	}

	for _, test := range tests {
		os.Setenv(test.envName, test.envValue)
		os.Args = append([]string{"cmd"}, test.args...)

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		config, _ := NewConfig()

		os.Clearenv()

		assert.Equal(t, test.want.host, config.HTTPAddress)
	}
}

func TestHostIntervalVariables(t *testing.T) {
	type want struct {
		interval time.Duration
	}
	tests := []struct {
		args     []string
		envName  string
		envValue string
		want     want
	}{
		{
			args:     []string{"-i", "200"},
			envName:  "STORE_INTERVAL",
			envValue: "400",
			want: want{
				interval: time.Duration(400) * time.Second,
			},
		},
		{
			args:     []string{"-f", "tst"},
			envName:  "STORE_INTERVAL",
			envValue: "400",
			want: want{
				interval: time.Duration(400) * time.Second,
			},
		},
		{
			args:     []string{},
			envName:  "SOME_VAR",
			envValue: "SOME_VAL",
			want: want{
				interval: time.Duration(300) * time.Second,
			},
		},
	}

	for _, test := range tests {
		os.Setenv(test.envName, test.envValue)
		os.Args = append([]string{"cmd"}, test.args...)

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		config, _ := NewConfig()

		os.Clearenv()

		assert.Equal(t, test.want.interval, config.StoreInterval)
	}
}

func TestHostRestoreVariables(t *testing.T) {
	type want struct {
		restore bool
	}
	tests := []struct {
		args     []string
		envName  string
		envValue string
		want     want
	}{
		{
			args:     []string{"-r", "false"},
			envName:  "RESTORE",
			envValue: "true",
			want: want{
				restore: true,
			},
		},
		{
			args:     []string{},
			envName:  "RESTORE",
			envValue: "true",
			want: want{
				restore: true,
			},
		},
		{
			args:     []string{},
			envName:  "RESTORE",
			envValue: "SOME_VAL",
			want: want{
				restore: false,
			},
		},
		{
			args:     []string{},
			envName:  "SOME_VAR",
			envValue: "SOME_VAL",
			want: want{
				restore: false,
			},
		},
	}

	for _, test := range tests {
		os.Setenv(test.envName, test.envValue)
		os.Args = append([]string{"cmd"}, test.args...)

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
		config, _ := NewConfig()

		os.Clearenv()

		assert.Equal(t, test.want.restore, config.Restore)
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

func TestCryptoKeyConfig(t *testing.T) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	privKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	})

	tmpFile, err := os.CreateTemp("", "private_key_*.pem")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(privKeyPEM)
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	type testCase struct {
		name          string
		args          []string
		envVars       map[string]string
		wantKeyLoaded bool
		wantErr       bool
		wantKeyExists bool
	}

	tests := []testCase{
		{
			name:          "flag - no key specified",
			args:          []string{},
			wantKeyLoaded: false,
			wantErr:       false,
			wantKeyExists: false,
		},
		{
			name:          "flag - valid key path",
			args:          []string{"-crypto-key", tmpFile.Name()},
			wantKeyLoaded: true,
			wantErr:       false,
			wantKeyExists: true,
		},
		{
			name:          "flag - invalid key path",
			args:          []string{"-crypto-key", "nonexistent.pem"},
			wantKeyLoaded: false,
			wantErr:       true,
			wantKeyExists: false,
		},
		{
			name:          "env - no key specified",
			envVars:       map[string]string{},
			wantKeyLoaded: false,
			wantErr:       false,
			wantKeyExists: false,
		},
		{
			name:          "env - valid key path",
			envVars:       map[string]string{"CRYPTO_KEY": tmpFile.Name()},
			wantKeyLoaded: true,
			wantErr:       false,
			wantKeyExists: true,
		},
		{
			name:          "env - invalid key path",
			envVars:       map[string]string{"CRYPTO_KEY": "nonexistent.pem"},
			wantKeyLoaded: false,
			wantErr:       true,
			wantKeyExists: false,
		},
		{
			name:          "flag and env - env overrides flag with valid path",
			args:          []string{"-crypto-key", "nonexistent.pem"},
			envVars:       map[string]string{"CRYPTO_KEY": tmpFile.Name()},
			wantKeyLoaded: true,
			wantErr:       false,
			wantKeyExists: true,
		},
		{
			name:          "flag and env - env overrides flag with invalid path",
			args:          []string{"-crypto-key", tmpFile.Name()},
			envVars:       map[string]string{"CRYPTO_KEY": "nonexistent.pem"},
			wantKeyLoaded: false,
			wantErr:       true,
			wantKeyExists: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				for k := range test.envVars {
					os.Unsetenv(k)
				}
			}()

			os.Args = append([]string{"cmd"}, test.args...)
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

			config, err := NewConfig()

			if test.wantErr {
				require.Error(t, err, "Expected error but got none")
				return
			}
			require.NoError(t, err, "Unexpected error")

			if test.wantKeyExists {
				require.NotNil(t, config.CryptoKey, "Expected crypto key to be loaded")
			} else {
				require.Nil(t, config.CryptoKey, "Expected no crypto key to be loaded")
			}
		})
	}
}

func TestLoadPrivateKey(t *testing.T) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	privKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	})

	validTmpFile, err := os.CreateTemp("", "valid_private_key_*.pem")
	require.NoError(t, err)
	defer os.Remove(validTmpFile.Name())

	_, err = validTmpFile.Write(privKeyPEM)
	require.NoError(t, err)
	require.NoError(t, validTmpFile.Close())

	invalidTmpFile, err := os.CreateTemp("", "invalid_private_key_*.pem")
	require.NoError(t, err)
	defer os.Remove(invalidTmpFile.Name())

	_, err = invalidTmpFile.Write([]byte("invalid key data"))
	require.NoError(t, err)
	require.NoError(t, invalidTmpFile.Close())

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
			got, err := loadPrivateKey(test.keyPath)
			if (err != nil) != test.wantErr {
				t.Errorf("loadPrivateKey() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if (got != nil) != test.wantKey {
				t.Errorf("loadPrivateKey() got key = %v, want key %v", got != nil, test.wantKey)
			}
		})
	}
}

func TestServerConfigPriority(t *testing.T) {
	tmpDir := t.TempDir()

	configContent := `{
		"address": "json:8080",
		"store_interval": 10,
		"store_file": "/json/store.db",
		"restore": true,
		"database_dsn": "json_dsn",
		"key": "json_key"
	}`
	configPath := filepath.Join(tmpDir, "config.json")
	err := os.WriteFile(configPath, []byte(configContent), 0600)
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
				HTTPAddress:     "json:8080",
				StoreInterval:   10 * time.Second,
				FileStoragePath: "/json/store.db",
				Restore:         true,
				DatabaseDSN:     "json_dsn",
				Key:             "json_key",
			},
		},
		{
			name:       "keys overrides json",
			configFile: configPath,
			args:       []string{"-a", "flag:8080", "-i", "15", "-k", "flag_key"},
			expected: Config{
				HTTPAddress:     "flag:8080",
				StoreInterval:   15 * time.Second,
				FileStoragePath: "/json/store.db",
				Restore:         true,
				DatabaseDSN:     "json_dsn",
				Key:             "flag_key",
			},
		},
		{
			name:       "envs override all",
			configFile: configPath,
			envVars: map[string]string{
				"ADDRESS":           "env:8080",
				"STORE_INTERVAL":    "60",
				"FILE_STORAGE_PATH": "/env/store.db",
				"KEY":               "env_key",
			},
			args: []string{"-a", "flag:8080", "-i", "15", "-k", "flag_key"},
			expected: Config{
				HTTPAddress:     "env:8080",
				StoreInterval:   60 * time.Second,
				FileStoragePath: "/env/store.db",
				Restore:         true,
				DatabaseDSN:     "json_dsn",
				Key:             "env_key",
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
			assert.Equal(t, tt.expected.StoreInterval, cfg.StoreInterval)
			assert.Equal(t, tt.expected.FileStoragePath, cfg.FileStoragePath)
			assert.Equal(t, tt.expected.Restore, cfg.Restore)
			assert.Equal(t, tt.expected.DatabaseDSN, cfg.DatabaseDSN)
			assert.Equal(t, tt.expected.Key, cfg.Key)
		})
	}
}

func TestParseTrustedSubnetFlag(t *testing.T) {
	type want struct {
		subnet *net.IPNet
	}
	tests := []struct {
		name    string
		args    []string
		want    want
		wantErr bool
	}{
		{
			name: "valid CIDR flag",
			args: []string{"-t", "192.168.1.0/24"},
			want: want{
				subnet: &net.IPNet{
					IP:   net.IPv4(192, 168, 1, 0),
					Mask: net.IPv4Mask(255, 255, 255, 0),
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid CIDR flag",
			args:    []string{"-t", "invalid_cidr"},
			want:    want{subnet: nil},
			wantErr: true,
		},
		{
			name:    "no flag",
			args:    []string{},
			want:    want{subnet: nil},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Args = append([]string{"cmd"}, test.args...)
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

			config, err := NewConfig()

			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			if test.want.subnet == nil {
				assert.Nil(t, config.TrustedSubnet)
			} else {
				require.NotNil(t, config.TrustedSubnet)
				assert.Equal(t, test.want.subnet.String(), config.TrustedSubnet.String())
			}
		})
	}
}

func TestTrustedSubnetEnvVariables(t *testing.T) {
	type want struct {
		subnet *net.IPNet
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
			name:     "valid CIDR env",
			envName:  "TRUSTED_SUBNET",
			envValue: "10.0.0.0/8",
			want: want{
				subnet: &net.IPNet{
					IP:   net.IPv4(10, 0, 0, 0),
					Mask: net.IPv4Mask(255, 0, 0, 0),
				},
			},
			wantErr: false,
		},
		{
			name:     "invalid CIDR env",
			envName:  "TRUSTED_SUBNET",
			envValue: "invalid_cidr",
			want:     want{subnet: nil},
			wantErr:  true,
		},
		{
			name:     "no env",
			envName:  "OTHER_VAR",
			envValue: "value",
			want:     want{subnet: nil},
			wantErr:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Setenv(test.envName, test.envValue)
			defer os.Unsetenv(test.envName)

			os.Args = append([]string{"cmd"}, test.args...)
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

			config, err := NewConfig()

			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			if test.want.subnet == nil {
				assert.Nil(t, config.TrustedSubnet)
			} else {
				require.NotNil(t, config.TrustedSubnet)
				assert.Equal(t, test.want.subnet.String(), config.TrustedSubnet.String())
			}
		})
	}
}

func TestTrustedSubnetConfigFile(t *testing.T) {
	tmpDir := t.TempDir()

	configContent := `{
		"trusted_subnet": "172.16.0.0/12"
	}`
	configPath := filepath.Join(tmpDir, "config.json")
	err := os.WriteFile(configPath, []byte(configContent), 0600)
	require.NoError(t, err)

	os.Args = []string{"cmd", "-config", configPath}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	config, err := NewConfig()
	require.NoError(t, err)

	expected := &net.IPNet{
		IP:   net.IPv4(172, 16, 0, 0),
		Mask: net.IPv4Mask(255, 240, 0, 0),
	}
	require.NotNil(t, config.TrustedSubnet)
	assert.Equal(t, expected.String(), config.TrustedSubnet.String())
}

func TestTrustedSubnetPriority(t *testing.T) {
	tmpDir := t.TempDir()

	configContent := `{
		"trusted_subnet": "10.0.0.0/8"
	}`
	configPath := filepath.Join(tmpDir, "config.json")
	err := os.WriteFile(configPath, []byte(configContent), 0600)
	require.NoError(t, err)

	tests := []struct {
		name     string
		args     []string
		envVars  map[string]string
		expected string
	}{
		{
			name:     "config file only",
			args:     []string{"-config", configPath},
			expected: "10.0.0.0/8",
		},
		{
			name:     "flag overrides config",
			args:     []string{"-config", configPath, "-t", "192.168.1.0/24"},
			expected: "192.168.1.0/24",
		},
		{
			name: "env overrides all",
			args: []string{"-config", configPath, "-t", "192.168.1.0/24"},
			envVars: map[string]string{
				"TRUSTED_SUBNET": "172.16.0.0/12",
			},
			expected: "172.16.0.0/12",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				for k := range test.envVars {
					os.Unsetenv(k)
				}
			}()

			os.Args = append([]string{"cmd"}, test.args...)
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

			config, err := NewConfig()
			require.NoError(t, err)

			if test.expected == "" {
				assert.Nil(t, config.TrustedSubnet)
			} else {
				require.NotNil(t, config.TrustedSubnet)
				assert.Equal(t, test.expected, config.TrustedSubnet.String())
			}
		})
	}
}
