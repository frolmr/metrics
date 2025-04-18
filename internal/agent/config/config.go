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

	"github.com/frolmr/metrics/pkg/formatter"
)

const (
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
	var (
		serverScheme      string
		serverHTTPAddress string
		reportIntervalSec int
		pollIntervalSec   int
		key               string
		rateLimit         int
		cryptoKeyPath     string
	)

	flag.StringVar(&serverScheme, "s", defaultScheme, "server scheme: http or https")
	flag.StringVar(&serverHTTPAddress, "a", defaultAddress, "address and port of the server")
	flag.IntVar(&reportIntervalSec, "r", defaultReportIntervalSec, "report interval")
	flag.IntVar(&pollIntervalSec, "p", defaultPollIntervalSec, "poll interval")
	flag.StringVar(&key, "k", key, "encryption key")
	flag.IntVar(&rateLimit, "l", defaultRateLimit, "requests to server rate limit")
	flag.StringVar(&cryptoKeyPath, "crypto-key", key, "public crypto key path")
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

	if reportIntervalSecEnv, _ := strconv.Atoi(os.Getenv(reportIntervalEnvName)); reportIntervalSecEnv != 0 {
		reportIntervalSec = reportIntervalSecEnv
	}

	if pollIntervalSecEnv, _ := strconv.Atoi(os.Getenv(pollIntervalEnvName)); pollIntervalSecEnv != 0 {
		pollIntervalSec = pollIntervalSecEnv
	}

	if keyEnv := os.Getenv(keyEnv); keyEnv != "" {
		key = keyEnv
	}

	if rateLimitEnv, _ := strconv.Atoi(os.Getenv(rateLimitEnvName)); rateLimitEnv != 0 {
		rateLimit = rateLimitEnv
	}

	if cryptoKeyEnv := os.Getenv(cryptoKeyEnvName); cryptoKeyEnv != "" {
		cryptoKeyPath = cryptoKeyEnv
	}

	cryptoKey, err := loadPublicKey(cryptoKeyPath)
	if err != nil {
		return nil, err
	}

	return &Config{
		Scheme:         serverScheme,
		HTTPAddress:    serverHTTPAddress,
		ReportInterval: time.Duration(reportIntervalSec) * time.Second,
		PollInterval:   time.Duration(pollIntervalSec) * time.Second,
		Key:            key,
		RateLimit:      rateLimit,
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
