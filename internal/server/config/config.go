package config

import (
	"flag"
	"os"

	"github.com/frolmr/metrics.git/internal/common/utils"
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

func GetConfig() {
	var err error

	if ServerScheme = os.Getenv(schemeEnvName); ServerScheme == "" {
		ServerScheme = defaultScheme
	}

	if ServerAddress = os.Getenv(addressEnvName); ServerAddress == "" {
		ServerAddress = defaultAddress
	}

	flag.StringVar(&ServerScheme, "s", ServerScheme, "server scheme: http or https")
	flag.StringVar(&ServerAddress, "a", ServerAddress, "address and port of the server")

	flag.Parse()

	if err = utils.CheckSchemeFormat(ServerScheme); err != nil {
		panic(err)
	}

	if err = utils.CheckAddrFormat(ServerAddress); err != nil {
		panic(err)
	}
}
