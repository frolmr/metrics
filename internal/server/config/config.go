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

	"github.com/frolmr/metrics/pkg/fileconfig"
	"github.com/frolmr/metrics/pkg/formatter"
)

const (
	maxParamCount = 4

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
	schemeValues := make([]string, 0, maxParamCount)
	addressValues := make([]string, 0, maxParamCount)
	databaseValues := make([]string, 0, maxParamCount)
	filePathValues := make([]string, 0, maxParamCount)

	storeIntervalValues := make([]int, 0, maxParamCount)

	restoreValues := make([]bool, 0, maxParamCount)

	keyValues := make([]string, 0, maxParamCount)
	cryptoKeyValues := make([]string, 0, maxParamCount)

	var (
		serverScheme      string
		serverHTTPAddress string
		databaseDsn       string
		storeIntervalSec  int
		fileStoragePath   string
		restore           string
		key               string
		cryptoKeyPath     string
		profile           bool
		configFile        string
	)

	schemeValues = append(schemeValues, defaultScheme)
	addressValues = append(addressValues, defaultAddress)
	storeIntervalValues = append(storeIntervalValues, defaultStoreInterval)
	filePathValues = append(filePathValues, defaultFileStoragePath)
	restoreValues = append(restoreValues, defaultRestore)

	flag.StringVar(&serverScheme, "s", "", "server scheme: http or https")
	flag.StringVar(&serverHTTPAddress, "a", "", "address and port of the server")
	flag.StringVar(&databaseDsn, "d", "", "DB DSN")
	flag.IntVar(&storeIntervalSec, "i", 0, "snapshot data interval")
	flag.StringVar(&fileStoragePath, "f", "", "snapshot file path")
	flag.StringVar(&restore, "r", "", "bool flag for set snapshoting")
	flag.StringVar(&key, "k", "", "encryption key")
	flag.StringVar(&cryptoKeyPath, "crypto-key", "", "path to private key for decryption")
	flag.BoolVar(&profile, "p", profile, "bool flag for app profiling")
	flag.StringVar(&configFile, "config", "", "path to config file")
	flag.Parse()

	if configFile != "" {
		fileCfg, err := fileconfig.ReadServerConfig(configFile)
		if err != nil {
			return nil, err
		}
		if fileCfg != nil {
			if fileCfg.Address != "" {
				addressValues = append(addressValues, fileCfg.Address)
			}
			if fileCfg.Scheme != "" {
				schemeValues = append(schemeValues, fileCfg.Scheme)
			}
			if fileCfg.Key != "" {
				keyValues = append(keyValues, fileCfg.Key)
			}
			if fileCfg.CryptoKey != "" {
				cryptoKeyValues = append(cryptoKeyValues, fileCfg.CryptoKey)
			}
			if fileCfg.Restore {
				restoreValues = append(restoreValues, fileCfg.Restore)
			}
			if fileCfg.StoreIntervalSec != 0 {
				storeIntervalValues = append(storeIntervalValues, fileCfg.StoreIntervalSec)
			}
			if fileCfg.StoreFile != "" {
				filePathValues = append(filePathValues, fileCfg.StoreFile)
			}
			if fileCfg.DatabaseDSN != "" {
				databaseValues = append(databaseValues, fileCfg.DatabaseDSN)
			}
		}
	}

	if serverScheme != "" {
		schemeValues = append(schemeValues, serverScheme)
	}

	if serverHTTPAddress != "" {
		addressValues = append(addressValues, serverHTTPAddress)
	}

	if storeIntervalSec != 0 {
		storeIntervalValues = append(storeIntervalValues, storeIntervalSec)
	}

	if databaseDsn != "" {
		databaseValues = append(databaseValues, databaseDsn)
	}

	if fileStoragePath != "" {
		filePathValues = append(filePathValues, fileStoragePath)
	}

	if restore != "" {
		if restoreKey, err := strconv.ParseBool(restore); err == nil {
			restoreValues = append(restoreValues, restoreKey)
		}
	}

	if key != "" {
		keyValues = append(keyValues, key)
	}

	if cryptoKeyPath != "" {
		cryptoKeyValues = append(cryptoKeyValues, cryptoKeyPath)
	}

	if serverSchemeEnv := os.Getenv(schemeEnvName); serverSchemeEnv != "" {
		schemeValues = append(schemeValues, serverSchemeEnv)
	}

	if serverHTTPAddressEnv := os.Getenv(addressEnvName); serverHTTPAddressEnv != "" {
		addressValues = append(addressValues, serverHTTPAddressEnv)
	}

	if databaseDsnEnv := os.Getenv(databaseDsnEnv); databaseDsnEnv != "" {
		databaseValues = append(databaseValues, databaseDsnEnv)
	}

	if storeIntervalSecEnv, err := strconv.Atoi(os.Getenv(storeIntervalEnv)); storeIntervalSecEnv != 0 {
		if err == nil {
			storeIntervalValues = append(storeIntervalValues, storeIntervalSecEnv)
		}
	}

	if fileStoragePathEnv := os.Getenv(fileStoragePathEnv); fileStoragePathEnv != "" {
		filePathValues = append(filePathValues, fileStoragePathEnv)
	}

	if restoreEnv, err := strconv.ParseBool(os.Getenv(restoreEnv)); err == nil {
		restoreValues = append(restoreValues, restoreEnv)
	}

	if keyEnv := os.Getenv(keyEnv); keyEnv != "" {
		keyValues = append(keyValues, keyEnv)
	}

	if cryptoKeyEnv := os.Getenv(cryptoKeyEnvName); cryptoKeyEnv != "" {
		cryptoKeyValues = append(cryptoKeyValues, cryptoKeyEnv)
	}

	schemeConfig := schemeValues[len(schemeValues)-1]
	if err := formatter.CheckSchemeFormat(schemeConfig); err != nil {
		return nil, err
	}

	addressConfig := addressValues[len(addressValues)-1]
	if err := formatter.CheckAddrFormat(addressConfig); err != nil {
		return nil, err
	}

	storeIntervalConfig := storeIntervalValues[len(storeIntervalValues)-1]
	fileStorageConfig := filePathValues[len(filePathValues)-1]
	restoreConfig := restoreValues[len(restoreValues)-1]

	var databaseDSNConfig string
	if len(databaseValues) != 0 {
		databaseDSNConfig = databaseValues[len(databaseValues)-1]
	}

	var keyConfig string
	if len(keyValues) != 0 {
		keyConfig = keyValues[len(keyValues)-1]
	}

	var cryptoKeyConfig string
	if len(cryptoKeyValues) != 0 {
		cryptoKeyConfig = cryptoKeyValues[len(cryptoKeyValues)-1]
	}

	privateKey, err := loadPrivateKey(cryptoKeyConfig)
	if err != nil {
		return nil, err
	}

	return &Config{
		Scheme:          schemeConfig,
		HTTPAddress:     addressConfig,
		DatabaseDSN:     databaseDSNConfig,
		StoreInterval:   time.Duration(storeIntervalConfig) * time.Second,
		FileStoragePath: fileStorageConfig,
		Restore:         restoreConfig,
		Key:             keyConfig,
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
