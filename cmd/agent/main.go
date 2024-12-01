package main

import (
	"time"

	"github.com/frolmr/metrics.git/internal/agent/config"
	"github.com/frolmr/metrics.git/internal/agent/metrics"
)

func main() {
	config.GetConfig()

	gaugeMetrics := make(map[string]float64)
	counterMetrics := make(map[string]int64)

	go func() {
		for {
			counterMetrics, gaugeMetrics = metrics.CollectMetrics(counterMetrics, gaugeMetrics)
			time.Sleep(config.PollInterval)
		}
	}()

	for {
		metrics.ReportGaugeMetrics(gaugeMetrics)
		metrics.ReportCounterMetrics(counterMetrics)
		time.Sleep(config.ReportInterval)
	}
}
