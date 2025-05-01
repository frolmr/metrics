package metrics

import (
	"testing"

	"github.com/frolmr/metrics/internal/agent/config"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewMetricsCollection(t *testing.T) {
	tests := []struct {
		name     string
		reporter *resty.Client
		cfg      *config.Config
		want     *MetricsCollection
	}{
		{
			name: "basic initialization",
			want: &MetricsCollection{
				CounterMetrics: make(map[string]int64),
				GaugeMetrics:   make(map[string]float64),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMetricsCollection()

			assert.NotNil(t, got.CounterMetrics)
			assert.NotNil(t, got.GaugeMetrics)

			assert.Empty(t, got.CounterMetrics)
			assert.Empty(t, got.GaugeMetrics)
		})
	}
}
