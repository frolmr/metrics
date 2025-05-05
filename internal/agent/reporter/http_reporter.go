package reporter

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/frolmr/metrics/internal/agent/config"
	"github.com/frolmr/metrics/internal/agent/metrics"
	"github.com/frolmr/metrics/internal/domain"
	"github.com/frolmr/metrics/pkg/signer"
	"github.com/go-resty/resty/v2"
)

type HTTPReporter struct {
	config *config.Config
	client *resty.Client
}

func NewHTTPReporter(cfg *config.Config) *HTTPReporter {
	return &HTTPReporter{
		config: cfg,
		client: resty.New(),
	}
}

func (r *HTTPReporter) Close() error {
	if r.client != nil {
		r.client.GetClient().CloseIdleConnections()
	}
	return nil
}

// ReportMetrics functions sends http request to server with metrics collected
func (r *HTTPReporter) ReportMetrics(ms metrics.MetricsCollection) {
	metrics := make([]domain.Metrics, 0, len(ms.GaugeMetrics)+len(ms.CounterMetrics))

	for key, value := range ms.GaugeMetrics {
		metric := domain.Metrics{
			ID:    key,
			MType: domain.GaugeType,
			Value: &value,
		}
		metrics = append(metrics, metric)
	}
	for key, value := range ms.CounterMetrics {
		metric := domain.Metrics{
			ID:    key,
			MType: domain.CounterType,
			Delta: &value,
		}
		metrics = append(metrics, metric)
	}

	if len(metrics) == 0 {
		return
	}

	retryIntervals := []time.Duration{time.Second, time.Second * 2, time.Second * 5}

	for _, interval := range retryIntervals {
		err := r.reportMetrics(metrics)
		if err != nil && (errors.Is(err, context.DeadlineExceeded) || isConnectionRefused(err)) {
			time.Sleep(interval * time.Second)
			continue
		} else {
			return
		}
	}
}

func (r *HTTPReporter) reportMetrics(metrics []domain.Metrics) error {
	metricsJSON, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	encryptedPayload, err := r.encryptPayload(metricsJSON)
	if err != nil {
		log.Println("encryption failure ", err.Error())
		return err
	}

	compressedData, err := r.compressPayload(encryptedPayload)
	if err != nil {
		log.Println("compression failure ", err.Error())
		return err
	}

	hostIP, err := getOutboundIP()
	if err != nil {
		log.Println("failed to get host IP:", err)
		hostIP = "unknown"
	}

	cl := r.client.R().
		SetHeader("Content-Type", domain.JSONContentType).
		SetHeader("Content-Encoding", domain.CompressFormat).
		SetHeader("X-Real-IP", hostIP).
		SetBody(compressedData).
		SetPathParam("serverScheme", r.config.Scheme).
		SetPathParam("serverHost", r.config.HTTPAddress)

	var signature []byte
	if r.config.Key != "" {
		signature = signer.SignPayloadWithKey(metricsJSON, []byte(r.config.Key))
		cl.SetHeader(domain.SignatureHeader, hex.EncodeToString(signature))
	}

	resp, err := cl.
		Post("{serverScheme}://{serverHost}/updates/")

	if err != nil {
		log.Println("request failure ", err.Error())
		return err
	}

	log.Println("got resp code from server: ", resp.StatusCode())

	return nil
}

func (r *HTTPReporter) compressPayload(payload []byte) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	if _, err := zb.Write(payload); err != nil {
		return nil, err
	}

	zb.Close()

	return buf, nil
}

func (r *HTTPReporter) encryptPayload(payload []byte) ([]byte, error) {
	var encryptedPayload []byte
	if r.config.CryptoKey != nil {
		// Calculate maximum chunk size (for 2048-bit key: 245 bytes)
		maxChunkSize := r.config.CryptoKey.Size() - 11

		// Split payload into chunks
		chunks := chunkData(payload, maxChunkSize)

		// Encrypt each chunk
		for _, chunk := range chunks {
			encryptedChunk, err := rsa.EncryptPKCS1v15(rand.Reader, r.config.CryptoKey, chunk)
			if err != nil {
				return nil, fmt.Errorf("RSA encryption failed: %w", err)
			}
			encryptedPayload = append(encryptedPayload, encryptedChunk...)
		}
	} else {
		encryptedPayload = payload
	}

	return encryptedPayload, nil
}

func getOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func chunkData(data []byte, chunkSize int) [][]byte {
	var chunks [][]byte
	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize
		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
	}
	return chunks
}

func isConnectionRefused(err error) bool {
	var netErr *net.OpError
	if errors.As(err, &netErr) {
		var syscallErr *os.SyscallError
		if errors.As(netErr.Err, &syscallErr) {
			return errors.Is(syscallErr.Err, syscall.ECONNREFUSED)
		}
		return errors.Is(netErr.Err, syscall.ECONNREFUSED)
	}

	var syscallErr *os.SyscallError
	if errors.As(err, &syscallErr) {
		return errors.Is(syscallErr.Err, syscall.ECONNREFUSED)
	}

	return errors.Is(err, syscall.ECONNREFUSED)
}
