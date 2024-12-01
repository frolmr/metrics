package main

import (
	"time"

	"github.com/frolmr/metrics.git/internal/agent/metrics"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	gaugeMetrics := make(map[string]float64)
	counterMetrics := make(map[string]int64)

	go func() {
		for {
			counterMetrics, gaugeMetrics = metrics.CollectMetrics(counterMetrics, gaugeMetrics)
			time.Sleep(pollInterval)
		}
	}()

	for {
		metrics.ReportGaugeMetrics(gaugeMetrics)
		metrics.ReportCounterMetrics(counterMetrics)
		time.Sleep(reportInterval)
	}
}