package config

import (
	"flag"
	"os"

	"github.com/frolmr/metrics.git/pkg/formatter"
)

var (
	ServerScheme  string
	ServerAddress string
)

const (
	schemeEnvName  = "SCHEME"
	addressEnvName = "ADDRESS"

	defaultScheme  = "http"
	defaultAddress = "localhost:8080"
)

func GetConfig() error {
	if ServerScheme = os.Getenv(schemeEnvName); ServerScheme == "" {
		ServerScheme = defaultScheme
	}

	if ServerAddress = os.Getenv(addressEnvName); ServerAddress == "" {
		ServerAddress = defaultAddress
	}

	flag.StringVar(&ServerScheme, "s", ServerScheme, "server scheme: http or https")
	flag.StringVar(&ServerAddress, "a", ServerAddress, "address and port of the server")

	flag.Parse()

	if err := formatter.CheckSchemeFormat(ServerScheme); err != nil {
		return err
	}

	if err := formatter.CheckAddrFormat(ServerAddress); err != nil {
		return err
	}

	return nil
}
