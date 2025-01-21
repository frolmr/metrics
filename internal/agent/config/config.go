package config

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/frolmr/metrics.git/pkg/formatter"
)

const (
	schemeEnvName         = "SCHEME"
	addressEnvName        = "ADDRESS"
	reportIntervalEnvName = "REPORT_INTERVAL"
	pollIntervalEnvName   = "POLL_INTERVAL"
	keyEnv                = "KEY"

	defaultScheme            = "http"
	defaultAddress           = "localhost:8080"
	defaultReportIntervalSec = 10
	defaultPollIntervalSec   = 2
)

type Config struct {
	Scheme      string
	HTTPAddress string

	ReportInterval time.Duration
	PollInterval   time.Duration

	Key string
}

func NewConfig() (*Config, error) {
	serverScheme := os.Getenv(schemeEnvName)
	if serverScheme == "" {
		serverScheme = defaultScheme
	}
	flag.StringVar(&serverScheme, "s", serverScheme, "server scheme: http or https")
	if err := formatter.CheckSchemeFormat(serverScheme); err != nil {
		return nil, err
	}

	serverHTTPAddress := os.Getenv(addressEnvName)
	if serverHTTPAddress == "" {
		serverHTTPAddress = defaultAddress
	}
	flag.StringVar(&serverHTTPAddress, "a", serverHTTPAddress, "address and port of the server")
	if err := formatter.CheckAddrFormat(serverHTTPAddress); err != nil {
		return nil, err
	}

	reportIntervalSec, err := strconv.Atoi(os.Getenv(reportIntervalEnvName))
	if err != nil {
		reportIntervalSec = defaultReportIntervalSec
	}
	flag.IntVar(&reportIntervalSec, "r", reportIntervalSec, "report interval")

	pollIntervalSec, err := strconv.Atoi(os.Getenv(pollIntervalEnvName))
	if err != nil {
		pollIntervalSec = defaultPollIntervalSec
	}
	flag.IntVar(&pollIntervalSec, "p", pollIntervalSec, "poll interval")

	key := os.Getenv(keyEnv)
	flag.StringVar(&key, "k", key, "encryption key")

	flag.Parse()

	return &Config{
		Scheme:         serverScheme,
		HTTPAddress:    serverHTTPAddress,
		ReportInterval: time.Duration(reportIntervalSec) * time.Second,
		PollInterval:   time.Duration(pollIntervalSec) * time.Second,
		Key:            key,
	}, nil
}
