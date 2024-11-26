package storage

type Repository interface {
	UpdateCounterMetric(name string, value int64)
	UpdateGaugeMetric(name string, value float64)

	GetCounterMetric(name string) (int64, error)
	GetGaugeMetric(name string) (float64, error)

	GetCounterMetrics() map[string]int64
	GetGaugeMetrics() map[string]float64
}
