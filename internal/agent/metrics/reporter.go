package metrics

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"log"

	"github.com/frolmr/metrics.git/internal/agent/config"
	"github.com/frolmr/metrics.git/internal/domain"
)

type MetricsReporter interface {
	ReportMetrics()
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

	mc.reportMetric(metrics)
}

func (mc *MetricsCollection) reportMetric(metrics []domain.Metrics) {
	compressedData, err := mc.compressPayload(metrics)
	if err != nil {
		log.Println("compression failure ", err.Error())
		return
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
		return
	}

	log.Println(resp.StatusCode())
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
