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
		ParseFlags()

		assert.Equal(t, test.want.host, ServerAddress)
		assert.Equal(t, test.want.reportInterval, ReportInterval)
		assert.Equal(t, test.want.pollInterval, PollInterval)
	}
}
