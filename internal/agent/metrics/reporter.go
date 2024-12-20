package metrics

import (
	"github.com/frolmr/metrics.git/internal/agent/config"
	"github.com/frolmr/metrics.git/internal/domain"
	"github.com/go-resty/resty/v2"
)

type MetricsReporter interface {
	ReportCounterMetric(client *resty.Client)
	ReportGaugeMetric(client *resty.Client)
}

func (mc *MetricsCollection) reportMetric(metric domain.Metrics, client *resty.Client) {
	_, err := client.R().
		SetHeader("Content-Type", domain.JSONContentType).
		SetBody(metric).
		SetPathParam("serverScheme", config.ServerScheme).
		SetPathParam("serverHost", config.ServerAddress).
		Post("{serverScheme}://{serverHost}/update")

	if err != nil {
		return
	}
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
