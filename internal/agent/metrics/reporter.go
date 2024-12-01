package metrics

import (
	"fmt"

	"github.com/frolmr/metrics.git/internal/agent/config"
	"github.com/frolmr/metrics.git/internal/common/constants"
	"github.com/frolmr/metrics.git/internal/common/utils"
	"github.com/go-resty/resty/v2"
)

func reportMetric(metricType, metricName, metricValue string) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", constants.ContentType).
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

func ReportCounterMetrics(counterMetrics map[string]int64) {
	for key, value := range counterMetrics {
		reportMetric(constants.CounterType, key, utils.IntToString(value))
	}
}

func ReportGaugeMetrics(gaugeMetrics map[string]float64) {
	for key, value := range gaugeMetrics {
		reportMetric(constants.GaugeType, key, utils.FloatToString(value))
	}
}
