package metrics

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/frolmr/metrics.git/internal/domain"
	"github.com/frolmr/metrics.git/pkg/signer"
)

type MetricsReporter interface {
	ReportMetrics()
}

func isConnectionRefused(err error) bool {
	var netErr *net.OpError
	if errors.As(err, &netErr) {
		if errors.Is(netErr.Err.(*os.SyscallError), syscall.ECONNREFUSED) {
			return true
		}
	}
	return false
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

	compressedData, err := mc.compressPayload(metricsJSON)
	if err != nil {
		log.Println("compression failure ", err.Error())
		return err
	}

	cl := mc.ReportClinet.R().
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

func (mc *MetricsCollection) compressPayload(metricsJSON []byte) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	if _, err := zb.Write(metricsJSON); err != nil {
		return nil, err
	}

	zb.Close()

	return buf, nil
}
