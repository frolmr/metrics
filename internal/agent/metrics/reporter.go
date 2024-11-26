package metrics

import (
	"fmt"

	"github.com/frolmr/metrics.git/internal/agent/config"
	"github.com/frolmr/metrics.git/internal/domain"
	"github.com/frolmr/metrics.git/pkg/utils"
	"github.com/go-resty/resty/v2"
)

type MetricsReporter interface {
	ReportCounterMetric(client *resty.Client)
	ReportGaugeMetric(client *resty.Client)
}

func (mc *MetricsCollection) reportMetric(metricType, metricName, metricValue string, client *resty.Client) {
	resp, err := client.R().
		SetHeader("Content-Type", domain.ContentType).
		SetPathParams(map[string]string{
			"serverScheme": config.ServerScheme,
			"serverHost":   config.ServerAddress,
			"metricType":   metricType,
			"metricName":   metricName,
			"metricValue":  metricValue,
		}).Post("{serverScheme}://{serverHost}/update/{metricType}/{metricName}/{metricValue}")

	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	fmt.Println(resp)
}

func (mc *MetricsCollection) ReportCounterMetrics(client *resty.Client) {
	for key, value := range mc.CounterMetrics {
		mc.reportMetric(domain.CounterType, key, utils.IntToString(value), client)
	}
}

func (mc *MetricsCollection) ReportGaugeMetrics(client *resty.Client) {
	for key, value := range mc.GaugeMetrics {
		mc.reportMetric(domain.GaugeType, key, utils.FloatToString(value), client)
	}
}
