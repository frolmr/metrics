package config

import (
	"flag"
	"time"

	"github.com/frolmr/metrics.git/internal/common/utils"
)

var (
	ServerScheme  string
	ServerAddress string

	reportIntervalSec int
	pollIntervalSec   int

	ReportInterval time.Duration
	PollInterval   time.Duration
)

func ParseFlags() {
	flag.StringVar(&ServerScheme, "s", "http", "server scheme: http or https")
	flag.StringVar(&ServerAddress, "a", "localhost:8080", "address and port of the server")
	flag.IntVar(&reportIntervalSec, "r", 10, "report interval")
	flag.IntVar(&pollIntervalSec, "p", 2, "poll interval")

	flag.Parse()

	if err := utils.CheckAddrFormat(ServerAddress); err != nil {
		panic(err)
	}

	if err := utils.CheckSchemeFormat(ServerScheme); err != nil {
		panic(err)
	}

	ReportInterval = time.Duration(reportIntervalSec) * time.Second
	PollInterval = time.Duration(pollIntervalSec) * time.Second
}
