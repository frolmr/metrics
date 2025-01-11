package config

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
				interval: time.Duration(200) * time.Second,
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
				restore: false,
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
