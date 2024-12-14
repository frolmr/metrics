package metrics

import "github.com/go-resty/resty/v2"

var (
	pollCount int64
)

type MetricsCollection struct {
	CounterMetrics map[string]int64
	GaugeMetrics   map[string]float64

	ReportClinet *resty.Client
}

func NewMetricsCollection(reporter *resty.Client) *MetricsCollection {
	return &MetricsCollection{
		CounterMetrics: make(map[string]int64),
		GaugeMetrics:   make(map[string]float64),

		ReportClinet: reporter,
	}
}
