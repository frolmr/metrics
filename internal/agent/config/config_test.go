package config

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
