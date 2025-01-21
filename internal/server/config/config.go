package config

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/frolmr/metrics.git/pkg/formatter"
)

const (
	schemeEnvName      = "SCHEME"
	addressEnvName     = "ADDRESS"
	storeIntervalEnv   = "STORE_INTERVAL"
	fileStoragePathEnv = "FILE_STORAGE_PATH"
	restoreEnv         = "RESTORE"
	databaseDsnEnv     = "DATABASE_DSN"
	keyEnv             = "KEY"

	defaultScheme          = "http"
	defaultAddress         = "localhost:8080"
	defaultStoreInterval   = 300
	defaultFileStoragePath = "data_snapshot"
	defaultRestore         = false
)

type Config struct {
	Scheme      string
	HTTPAddress string
	DatabaseDSN string

	StoreInterval   time.Duration
	FileStoragePath string
	Restore         bool

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

	databaseDsn := os.Getenv(databaseDsnEnv)
	flag.StringVar(&databaseDsn, "d", databaseDsn, "DB DSN")

	storeIntervalSec, _ := strconv.Atoi(os.Getenv(storeIntervalEnv))
	flag.IntVar(&storeIntervalSec, "i", storeIntervalSec, "snapshot data interval")

	if storeIntervalSec == 0 {
		storeIntervalSec = defaultStoreInterval
	}

	fileStoragePath := os.Getenv(fileStoragePathEnv)
	if fileStoragePath == "" {
		fileStoragePath = defaultFileStoragePath
	}

	flag.StringVar(&fileStoragePath, "f", fileStoragePath, "snapshot file path")

	restoreString := os.Getenv(restoreEnv)
	flag.StringVar(&restoreString, "r", restoreString, "bool flag for set snapshoting")

	var key string
	flag.StringVar(&key, "k", key, "encryption key")

	flag.Parse()
	key = os.Getenv(keyEnv)

	restore := defaultRestore
	if restoreString == "true" {
		restore = true
	}

	return &Config{
		Scheme:          serverScheme,
		HTTPAddress:     serverHTTPAddress,
		DatabaseDSN:     databaseDsn,
		StoreInterval:   time.Duration(storeIntervalSec) * time.Second,
		FileStoragePath: fileStoragePath,
		Restore:         restore,
		Key:             key,
	}, nil
}
