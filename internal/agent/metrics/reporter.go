package metrics

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/frolmr/metrics.git/internal/agent/config"
	"github.com/frolmr/metrics.git/internal/domain"
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
	compressedData, err := mc.compressPayload(metrics)
	if err != nil {
		log.Println("compression failure ", err.Error())
		return err
	}

	resp, err := mc.ReportClinet.R().
		SetHeader("Content-Type", domain.JSONContentType).
		SetHeader("Content-Encoding", domain.CompressFormat).
		SetBody(compressedData).
		SetPathParam("serverScheme", config.ServerScheme).
		SetPathParam("serverHost", config.ServerAddress).
		Post("{serverScheme}://{serverHost}/updates/")

	if err != nil {
		log.Println("request failure ", err.Error())
		return err
	}

	log.Println("got resp code from server: ", resp.StatusCode())
	return nil
}

func (mc *MetricsCollection) compressPayload(metrics []domain.Metrics) (*bytes.Buffer, error) {
	metricJSON, err := json.Marshal(metrics)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buf)
	_, err = zb.Write(metricJSON)
	if err != nil {
		return nil, err
	}

	zb.Close()

	return buf, nil
}
