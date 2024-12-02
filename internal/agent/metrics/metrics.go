package metrics

var (
	pollCount int64
)

type MetricsCollection struct {
	CounterMetrics map[string]int64
	GaugeMetrics   map[string]float64
}

func NewMetricsCollection() *MetricsCollection {
	return &MetricsCollection{
		CounterMetrics: make(map[string]int64),
		GaugeMetrics:   make(map[string]float64),
	}
}
