package storage

type Repository interface {
	Ping() error
	UpdateCounterMetric(name string, value int64) error
	UpdateGaugeMetric(name string, value float64) error

	GetCounterMetric(name string) (int64, error)
	GetGaugeMetric(name string) (float64, error)

	GetCounterMetrics() map[string]int64
	GetGaugeMetrics() map[string]float64
}
