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
	ReportCounterMetric()
	ReportGaugeMetric()
}

func (mc *MetricsCollection) ReportCounterMetrics() {
	for key, value := range mc.CounterMetrics {
		metric := domain.Metrics{
			ID:    key,
			MType: domain.CounterType,
			Delta: &value,
		}
		mc.reportMetric(metric)
	}
}

func (mc *MetricsCollection) ReportGaugeMetrics() {
	for key, value := range mc.GaugeMetrics {
		metric := domain.Metrics{
			ID:    key,
			MType: domain.GaugeType,
			Value: &value,
		}
		mc.reportMetric(metric)
	}
}

func (mc *MetricsCollection) reportMetric(metric domain.Metrics) {
	compressedData, err := mc.compressPayload(metric)
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
		Post("{serverScheme}://{serverHost}/update")

	if err != nil {
		log.Println("request failure ", err.Error())
		return
	}

	log.Println(resp)
}

func (mc *MetricsCollection) compressPayload(metric domain.Metrics) (*bytes.Buffer, error) {
	metricJSON, err := json.Marshal(metric)
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
