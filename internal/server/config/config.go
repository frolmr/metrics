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

	storeIntervalString string
	storeIntervalSec    int
	restoreString       string

	StoreInterval   time.Duration
	FileStoragePath string
	Restore         bool
)

const (
	schemeEnvName  = "SCHEME"
	addressEnvName = "ADDRESS"

	storeIntervalEnv   = "STORE_INTERVAL"
	fileStoragePathEnv = "FILE_STORAGE_PATH"
	restoreEnv         = "RESTORE"

	defaultScheme  = "http"
	defaultAddress = "localhost:8080"

	defaultStoreIntervalString = "300"
	defaultStoreInterval       = 300
	defaultFileStoragePath     = "data_snapshot"
	defaultRestoreString       = "true"
)

func GetConfig() error {
	if ServerScheme = os.Getenv(schemeEnvName); ServerScheme == "" {
		ServerScheme = defaultScheme
	}

	if ServerAddress = os.Getenv(addressEnvName); ServerAddress == "" {
		ServerAddress = defaultAddress
	}

	if storeIntervalString = os.Getenv(storeIntervalEnv); storeIntervalString == "" {
		storeIntervalString = defaultStoreIntervalString
	}
	var err error
	if storeIntervalSec, err = strconv.Atoi(storeIntervalString); err != nil {
		storeIntervalSec = 300
	}

	if FileStoragePath = os.Getenv(fileStoragePathEnv); FileStoragePath == "" {
		FileStoragePath = defaultFileStoragePath
	}

	if restoreString = os.Getenv(restoreEnv); restoreString == "" {
		restoreString = defaultRestoreString
	}

	flag.StringVar(&ServerScheme, "s", ServerScheme, "server scheme: http or https")
	flag.StringVar(&ServerAddress, "a", ServerAddress, "address and port of the server")

	flag.IntVar(&storeIntervalSec, "i", storeIntervalSec, "snapshot data interval")
	flag.StringVar(&FileStoragePath, "f", FileStoragePath, "snapshot file path")
	flag.StringVar(&restoreString, "r", restoreString, "bool flag for set snapshoting")

	flag.Parse()

	if err := formatter.CheckSchemeFormat(ServerScheme); err != nil {
		return err
	}

	if err := formatter.CheckAddrFormat(ServerAddress); err != nil {
		return err
	}

	StoreInterval = time.Duration(storeIntervalSec) * time.Second

	Restore = false
	if restoreString == "true" {
		Restore = true
	}

	return nil
}
