// Package config to read flags and env for further server setup.
package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/frolmr/metrics/pkg/formatter"
)

const (
	schemeEnvName      = "SCHEME"
	addressEnvName     = "ADDRESS"
	storeIntervalEnv   = "STORE_INTERVAL"
	fileStoragePathEnv = "FILE_STORAGE_PATH"
	restoreEnv         = "RESTORE"
	databaseDsnEnv     = "DATABASE_DSN"
	keyEnv             = "KEY"
	cryptoKeyEnvName   = "CRYPTO_KEY"
)

const (
	defaultScheme          = "http"
	defaultAddress         = "localhost:8080"
	defaultStoreInterval   = 300
	defaultFileStoragePath = "data_snapshot"
	defaultRestore         = false
)

// Config structure to store server configuration.
type Config struct {
	Scheme      string
	HTTPAddress string
	DatabaseDSN string

	StoreInterval   time.Duration
	FileStoragePath string
	Restore         bool

	Key       string
	CryptoKey *rsa.PrivateKey
	Profiling bool
}

// NewConfig setups server config: read flags and env variables.
func NewConfig() (*Config, error) {
	var (
		serverScheme      string
		serverHTTPAddress string
		databaseDsn       string
		storeIntervalSec  int
		fileStoragePath   string
		restore           bool
		key               string
		cryptoKeyPath     string
		profile           bool
	)

	flag.StringVar(&serverScheme, "s", defaultScheme, "server scheme: http or https")
	flag.StringVar(&serverHTTPAddress, "a", defaultAddress, "address and port of the server")
	flag.StringVar(&databaseDsn, "d", databaseDsn, "DB DSN")
	flag.IntVar(&storeIntervalSec, "i", defaultStoreInterval, "snapshot data interval")
	flag.StringVar(&fileStoragePath, "f", defaultFileStoragePath, "snapshot file path")
	flag.BoolVar(&restore, "r", defaultRestore, "bool flag for set snapshoting")
	flag.StringVar(&key, "k", key, "encryption key")
	flag.StringVar(&cryptoKeyPath, "crypto-key", "", "path to private key for decryption")
	flag.BoolVar(&profile, "p", profile, "bool flag for app profiling")
	flag.Parse()

	if serverSchemeEnv := os.Getenv(schemeEnvName); serverSchemeEnv != "" {
		serverScheme = serverSchemeEnv
	}

	if err := formatter.CheckSchemeFormat(serverScheme); err != nil {
		return nil, err
	}

	if serverHTTPAddressEnv := os.Getenv(addressEnvName); serverHTTPAddressEnv != "" {
		serverHTTPAddress = serverHTTPAddressEnv
	}

	if err := formatter.CheckAddrFormat(serverHTTPAddress); err != nil {
		return nil, err
	}

	if databaseDsnEnv := os.Getenv(databaseDsnEnv); databaseDsnEnv != "" {
		databaseDsn = databaseDsnEnv
	}

	if storeIntervalSecEnv, _ := strconv.Atoi(os.Getenv(storeIntervalEnv)); storeIntervalSecEnv != 0 {
		storeIntervalSec = storeIntervalSecEnv
	}

	if fileStoragePathEnv := os.Getenv(fileStoragePathEnv); fileStoragePathEnv != "" {
		fileStoragePath = fileStoragePathEnv
	}

	if restoreEnv, err := strconv.ParseBool(os.Getenv(restoreEnv)); err == nil {
		restore = restoreEnv
	}

	if keyEnv := os.Getenv(keyEnv); keyEnv != "" {
		key = keyEnv
	}

	if cryptoKeyEnv := os.Getenv(cryptoKeyEnvName); cryptoKeyEnv != "" {
		cryptoKeyPath = cryptoKeyEnv
	}

	privateKey, err := loadPrivateKey(cryptoKeyPath)
	if err != nil {
		return nil, err
	}

	return &Config{
		Scheme:          serverScheme,
		HTTPAddress:     serverHTTPAddress,
		DatabaseDSN:     databaseDsn,
		StoreInterval:   time.Duration(storeIntervalSec) * time.Second,
		FileStoragePath: fileStoragePath,
		Restore:         restore,
		Key:             key,
		CryptoKey:       privateKey,
		Profiling:       profile,
	}, nil
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	if path == "" {
		return nil, nil
	}

	keyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the private key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return key.(*rsa.PrivateKey), nil
	}

	return priv, nil
}
