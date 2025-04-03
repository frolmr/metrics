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
			name:     "basic initialization",
			reporter: resty.New(),
			cfg: &config.Config{
				Scheme:      "http",
				HTTPAddress: "localhost:8080",
			},
			want: &MetricsCollection{
				CounterMetrics: make(map[string]int64),
				GaugeMetrics:   make(map[string]float64),
				ReportClient:   resty.New(),
				Config: &config.Config{
					Scheme:      "http",
					HTTPAddress: "localhost:8080",
				},
			},
		},
		{
			name:     "nil reporter",
			reporter: nil,
			cfg: &config.Config{
				Scheme:      "https",
				HTTPAddress: "example.com:443",
			},
			want: &MetricsCollection{
				CounterMetrics: make(map[string]int64),
				GaugeMetrics:   make(map[string]float64),
				ReportClient:   nil,
				Config: &config.Config{
					Scheme:      "https",
					HTTPAddress: "example.com:443",
				},
			},
		},
		{
			name:     "nil config",
			reporter: resty.New(),
			cfg:      nil,
			want: &MetricsCollection{
				CounterMetrics: make(map[string]int64),
				GaugeMetrics:   make(map[string]float64),
				ReportClient:   resty.New(),
				Config:         nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMetricsCollection(tt.reporter, tt.cfg)

			assert.NotNil(t, got.CounterMetrics)
			assert.NotNil(t, got.GaugeMetrics)

			assert.Empty(t, got.CounterMetrics)
			assert.Empty(t, got.GaugeMetrics)

			if tt.reporter == nil {
				assert.Nil(t, got.ReportClient)
			} else {
				assert.NotNil(t, got.ReportClient)
			}

			if tt.cfg == nil {
				assert.Nil(t, got.Config)
			} else {
				assert.NotNil(t, got.Config)
				assert.Equal(t, tt.cfg.Scheme, got.Config.Scheme)
				assert.Equal(t, tt.cfg.HTTPAddress, got.Config.HTTPAddress)
			}
		})
	}
}
