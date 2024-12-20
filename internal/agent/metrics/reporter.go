package metrics

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"log"

	"github.com/frolmr/metrics.git/internal/agent/config"
	"github.com/frolmr/metrics.git/internal/domain"
	"github.com/go-resty/resty/v2"
)

type MetricsReporter interface {
	ReportCounterMetric(client *resty.Client)
	ReportGaugeMetric(client *resty.Client)
}

func (mc *MetricsCollection) ReportCounterMetrics(client *resty.Client) {
	for key, value := range mc.CounterMetrics {
		metric := domain.Metrics{
			ID:    key,
			MType: domain.CounterType,
			Delta: &value,
		}
		mc.reportMetric(metric, client)
	}
}

func (mc *MetricsCollection) ReportGaugeMetrics(client *resty.Client) {
	for key, value := range mc.GaugeMetrics {
		metric := domain.Metrics{
			ID:    key,
			MType: domain.GaugeType,
			Value: &value,
		}
		mc.reportMetric(metric, client)
	}
}

func (mc *MetricsCollection) reportMetric(metric domain.Metrics, client *resty.Client) {
	compressedData, err := mc.compressPayload(metric)
	if err != nil {
		log.Println("compression failure ", err.Error())
		return
	}

	resp, err := client.R().
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
