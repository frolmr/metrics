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
		_ = GetConfig()

		assert.Equal(t, test.want.host, ServerAddress)
		assert.Equal(t, test.want.reportInterval, ReportInterval)
		assert.Equal(t, test.want.pollInterval, PollInterval)
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
				host:           "localhost:8081",
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
				reportInterval: 20 * time.Second,
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
				pollInterval:   5 * time.Second,
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
		_ = GetConfig()
		os.Unsetenv(test.envName)

		assert.Equal(t, test.want.host, ServerAddress)
		assert.Equal(t, test.want.reportInterval, ReportInterval)
		assert.Equal(t, test.want.pollInterval, PollInterval)
	}
}
