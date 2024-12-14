package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlags(t *testing.T) {
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
		_ = GetConfig()

		assert.Equal(t, test.want.host, ServerAddress)
	}
}

func TestEnvVariables(t *testing.T) {
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
				host: "localhost:8081",
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
		_ = GetConfig()

		os.Unsetenv(test.envName)

		assert.Equal(t, test.want.host, ServerAddress)
	}
}
