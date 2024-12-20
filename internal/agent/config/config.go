package config

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/frolmr/metrics.git/pkg/formatter"
)

var (
	ServerScheme  string
	ServerAddress string

	reportIntervalSec int
	pollIntervalSec   int

	reportIntervalSecString string
	pollIntervalSecString   string

	ReportInterval time.Duration
	PollInterval   time.Duration
)

const (
	schemeEnvName         = "SCHEME"
	addressEnvName        = "ADDRESS"
	reportIntervalEnvName = "REPORT_INTERVAL"
	pollIntervalEnvName   = "POLL_INTERVAL"

	defaultScheme            = "http"
	defaultAddress           = "localhost:8080"
	defaultReportIntervalSec = 10
	defaultPollIntervalSec   = 2
)

func GetConfig() error {
	var err error

	if ServerScheme = os.Getenv(schemeEnvName); ServerScheme == "" {
		ServerScheme = defaultScheme
	}

	if ServerAddress = os.Getenv(addressEnvName); ServerAddress == "" {
		ServerAddress = defaultAddress
	}

	if reportIntervalSecString = os.Getenv(reportIntervalEnvName); reportIntervalSecString == "" {
		reportIntervalSec = defaultReportIntervalSec
	} else {
		if reportIntervalSec, err = strconv.Atoi(reportIntervalSecString); err != nil {
			reportIntervalSec = defaultReportIntervalSec
		}
	}

	if pollIntervalSecString = os.Getenv(pollIntervalEnvName); pollIntervalSecString == "" {
		pollIntervalSec = defaultPollIntervalSec
	} else {
		if pollIntervalSec, err = strconv.Atoi(pollIntervalSecString); err != nil {
			pollIntervalSec = defaultPollIntervalSec
		}
	}

	flag.StringVar(&ServerScheme, "s", ServerScheme, "server scheme: http or https")
	flag.StringVar(&ServerAddress, "a", ServerAddress, "address and port of the server")
	flag.IntVar(&reportIntervalSec, "r", reportIntervalSec, "report interval")
	flag.IntVar(&pollIntervalSec, "p", pollIntervalSec, "poll interval")

	flag.Parse()

	if err := formatter.CheckSchemeFormat(ServerScheme); err != nil {
		return err
	}

	if err := formatter.CheckAddrFormat(ServerAddress); err != nil {
		return err
	}

	ReportInterval = time.Duration(reportIntervalSec) * time.Second
	PollInterval = time.Duration(pollIntervalSec) * time.Second

	return nil
}
