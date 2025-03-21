package metrics

import (
	"github.com/frolmr/metrics.git/internal/agent/config"
	"github.com/go-resty/resty/v2"
)

var (
	pollCount int64
)

type MetricsCollection struct {
	CounterMetrics map[string]int64
	GaugeMetrics   map[string]float64

	ReportClinet *resty.Client
	Config       *config.Config
}

// NewMetricsCollection is the constructor function for metrics collector and reporter.
func NewMetricsCollection(reporter *resty.Client, cfg *config.Config) *MetricsCollection {
	return &MetricsCollection{
		CounterMetrics: make(map[string]int64),
		GaugeMetrics:   make(map[string]float64),

		ReportClinet: reporter,
		Config:       cfg,
	}
}
