package metrics

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

	"github.com/frolmr/metrics/internal/domain"
	"github.com/frolmr/metrics/pkg/signer"
)

type MetricsReporter interface {
	ReportMetrics()
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

// ReportMetrics functions sends http request to server with metrics collected
func (mc *MetricsCollection) ReportMetrics() {
	metrics := make([]domain.Metrics, 0, len(mc.GaugeMetrics)+len(mc.CounterMetrics))

	for key, value := range mc.GaugeMetrics {
		metric := domain.Metrics{
			ID:    key,
			MType: domain.GaugeType,
			Value: &value,
		}
		metrics = append(metrics, metric)
	}
	for key, value := range mc.CounterMetrics {
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
		err := mc.reportMetrics(metrics)
		if err != nil && (errors.Is(err, context.DeadlineExceeded) || isConnectionRefused(err)) {
			time.Sleep(interval * time.Second)
			continue
		} else {
			return
		}
	}
}

func (mc *MetricsCollection) reportMetrics(metrics []domain.Metrics) error {
	metricsJSON, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	encryptedPayload, err := mc.encryptPayload(metricsJSON)
	if err != nil {
		log.Println("encryption failure ", err.Error())
		return err
	}

	compressedData, err := mc.compressPayload(encryptedPayload)
	if err != nil {
		log.Println("compression failure ", err.Error())
		return err
	}

	cl := mc.ReportClient.R().
		SetHeader("Content-Type", domain.JSONContentType).
		SetHeader("Content-Encoding", domain.CompressFormat).
		SetBody(compressedData).
		SetPathParam("serverScheme", mc.Config.Scheme).
		SetPathParam("serverHost", mc.Config.HTTPAddress)

	var signature []byte
	if mc.Config.Key != "" {
		signature = signer.SignPayloadWithKey(metricsJSON, []byte(mc.Config.Key))
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

func (mc *MetricsCollection) compressPayload(payload []byte) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	if _, err := zb.Write(payload); err != nil {
		return nil, err
	}

	zb.Close()

	return buf, nil
}

func (mc *MetricsCollection) encryptPayload(payload []byte) ([]byte, error) {
	var encryptedPayload []byte
	if mc.Config.CryptoKey != nil {
		// Calculate maximum chunk size (for 2048-bit key: 245 bytes)
		maxChunkSize := mc.Config.CryptoKey.Size() - 11

		// Split payload into chunks
		chunks := chunkData(payload, maxChunkSize)

		// Encrypt each chunk
		for _, chunk := range chunks {
			encryptedChunk, err := rsa.EncryptPKCS1v15(rand.Reader, mc.Config.CryptoKey, chunk)
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
