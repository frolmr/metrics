package config

import (
	"flag"

	"github.com/frolmr/metrics.git/internal/common/utils"
)

var (
	ServerScheme  string
	ServerAddress string
)

func ParseFlags() {
	flag.StringVar(&ServerScheme, "s", "http", "server scheme: http or https")
	flag.StringVar(&ServerAddress, "a", "localhost:8080", "address and port of the server")

	flag.Parse()

	if err := utils.CheckAddrFormat(ServerAddress); err != nil {
		panic(err)
	}

	if err := utils.CheckSchemeFormat(ServerScheme); err != nil {
		panic(err)
	}
}
