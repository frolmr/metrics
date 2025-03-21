package domain

// Metrics represents a metric with its name, type, and value.
// @Description Metrics request payload for metrics data.
type Metrics struct {
	// ID is the name of the metric.
	// Example: "cpu_usage"
	ID string `json:"id"`

	// MType is the type of the metric (gauge or counter).
	// Example: "gauge"
	MType string `json:"type"`

	// Delta is the value of the metric if it's a counter.
	// Example: 10
	Delta *int64 `json:"delta,omitempty"`

	// Value is the value of the metric if it's a gauge.
	// Example: 3.14
	Value *float64 `json:"value,omitempty"`
}
