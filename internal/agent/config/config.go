// Package config to read flags and env for further agent setup.
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

	schemeEnvName         = "SCHEME"
	addressEnvName        = "ADDRESS"
	reportIntervalEnvName = "REPORT_INTERVAL"
	pollIntervalEnvName   = "POLL_INTERVAL"
	keyEnv                = "KEY"
	rateLimitEnvName      = "RATE_LIMIT"
	cryptoKeyEnvName      = "CRYPTO_KEY"

	defaultScheme            = "http"
	defaultAddress           = "localhost:8080"
	defaultReportIntervalSec = 10
	defaultPollIntervalSec   = 2
	defaultRateLimit         = 5
)

// Config atructure to store agents configuration.
type Config struct {
	Scheme      string
	HTTPAddress string

	ReportInterval time.Duration
	PollInterval   time.Duration

	Key string

	RateLimit int

	CryptoKey *rsa.PublicKey
}

// NewConfig setups agents config: read flags and env variables.
func NewConfig() (*Config, error) {

	schemeValues := make([]string, 0, maxParamCount)
	addressValues := make([]string, 0, maxParamCount)
	reportIntervalValues := make([]int, 0, maxParamCount)
	pollIntervalValues := make([]int, 0, maxParamCount)
	rateLimitValues := make([]int, 0, maxParamCount)

	keyValues := make([]string, 0, maxParamCount)
	cryptoKeyValues := make([]string, 0, maxParamCount)

	var (
		serverScheme      string
		serverHTTPAddress string
		reportIntervalSec int
		pollIntervalSec   int
		key               string
		rateLimit         int
		cryptoKeyPath     string
		configFile        string
	)

	schemeValues = append(schemeValues, defaultScheme)
	addressValues = append(addressValues, defaultAddress)
	reportIntervalValues = append(reportIntervalValues, defaultReportIntervalSec)
	pollIntervalValues = append(pollIntervalValues, defaultPollIntervalSec)
	rateLimitValues = append(rateLimitValues, defaultRateLimit)

	flag.StringVar(&serverScheme, "s", "", "server scheme: http or https")
	flag.StringVar(&serverHTTPAddress, "a", "", "address and port of the server")
	flag.IntVar(&reportIntervalSec, "r", 0, "report interval")
	flag.IntVar(&pollIntervalSec, "p", 0, "poll interval")
	flag.IntVar(&rateLimit, "l", 0, "requests to server rate limit")
	flag.StringVar(&key, "k", "", "encryption key")
	flag.StringVar(&cryptoKeyPath, "crypto-key", "", "public crypto key path")
	flag.StringVar(&configFile, "config", "", "path to config file")
	flag.Parse()

	if configFile != "" {
		fileCfg, err := fileconfig.ReadAgentConfig(configFile)
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
			if fileCfg.ReportIntervalSec != 0 {
				reportIntervalValues = append(reportIntervalValues, fileCfg.ReportIntervalSec)
			}
			if fileCfg.PollIntervalSec != 0 {
				pollIntervalValues = append(pollIntervalValues, fileCfg.PollIntervalSec)
			}
			if fileCfg.RateLimit != 0 {
				rateLimitValues = append(rateLimitValues, fileCfg.RateLimit)
			}
			if fileCfg.Key != "" {
				keyValues = append(keyValues, fileCfg.Key)
			}
			if fileCfg.CryptoKey != "" {
				cryptoKeyValues = append(cryptoKeyValues, fileCfg.CryptoKey)
			}
		}
	}

	if serverScheme != "" {
		schemeValues = append(schemeValues, serverScheme)
	}

	if serverHTTPAddress != "" {
		addressValues = append(addressValues, serverHTTPAddress)
	}

	if reportIntervalSec != 0 {
		reportIntervalValues = append(reportIntervalValues, reportIntervalSec)
	}

	if pollIntervalSec != 0 {
		pollIntervalValues = append(pollIntervalValues, pollIntervalSec)
	}

	if rateLimit != 0 {
		rateLimitValues = append(rateLimitValues, rateLimit)
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

	if reportIntervalSecEnv, err := strconv.Atoi(os.Getenv(reportIntervalEnvName)); reportIntervalSecEnv != 0 {
		if err == nil {
			reportIntervalValues = append(reportIntervalValues, reportIntervalSecEnv)
		}
	}

	if pollIntervalSecEnv, err := strconv.Atoi(os.Getenv(pollIntervalEnvName)); pollIntervalSecEnv != 0 {
		if err == nil {
			pollIntervalValues = append(pollIntervalValues, pollIntervalSecEnv)
		}
	}

	if rateLimitEnv, _ := strconv.Atoi(os.Getenv(rateLimitEnvName)); rateLimitEnv != 0 {
		rateLimitValues = append(rateLimitValues, rateLimitEnv)
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

	reportIntervalConfig := reportIntervalValues[len(reportIntervalValues)-1]
	pollIntervalConfig := pollIntervalValues[len(pollIntervalValues)-1]
	rateLimitConfig := rateLimitValues[len(rateLimitValues)-1]

	var keyConfig string
	if len(keyValues) != 0 {
		keyConfig = keyValues[len(keyValues)-1]
	}

	var cryptoKeyConfig string
	if len(cryptoKeyValues) != 0 {
		cryptoKeyConfig = cryptoKeyValues[len(cryptoKeyValues)-1]
	}

	cryptoKey, err := loadPublicKey(cryptoKeyConfig)
	if err != nil {
		return nil, err
	}

	return &Config{
		Scheme:         schemeConfig,
		HTTPAddress:    addressConfig,
		ReportInterval: time.Duration(reportIntervalConfig) * time.Second,
		PollInterval:   time.Duration(pollIntervalConfig) * time.Second,
		Key:            keyConfig,
		RateLimit:      rateLimitConfig,
		CryptoKey:      cryptoKey,
	}, nil
}

func loadPublicKey(publicKeyPath string) (*rsa.PublicKey, error) {
	if publicKeyPath == "" {
		return nil, nil
	}

	keyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pub.(*rsa.PublicKey), nil
}
