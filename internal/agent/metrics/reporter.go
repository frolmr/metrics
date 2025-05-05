package metrics

type MetricsReporter interface {
	ReportMetrics(MetricsCollection)
	Close() error
}
