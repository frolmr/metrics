package storage

import (
	"testing"

	"github.com/frolmr/metrics.git/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestMemStorage(t *testing.T) {
	ms := MemStorage{
		CounterMetrics: map[string]int64{"cm1": 1},
		GaugeMetrics:   map[string]float64{"gm1": 0.2},
	}

	t.Run("update existing counter metric", func(t *testing.T) {
		_ = ms.UpdateCounterMetric("cm1", 2)
		assert.Equal(t, ms.CounterMetrics["cm1"], int64(3))
	})

	t.Run("update new counter metric", func(t *testing.T) {
		_ = ms.UpdateCounterMetric("cm2", 5)
		assert.Equal(t, ms.CounterMetrics["cm2"], int64(5))
	})

	t.Run("update existing gauge metric", func(t *testing.T) {
		_ = ms.UpdateGaugeMetric("gm1", 0.5)
		assert.Equal(t, ms.GaugeMetrics["gm1"], float64(0.5))
	})

	t.Run("update new counter metric", func(t *testing.T) {
		_ = ms.UpdateGaugeMetric("gm2", 1.2)
		assert.Equal(t, ms.GaugeMetrics["gm2"], float64(1.2))
	})

	t.Run("get counter metric", func(t *testing.T) {
		val, _ := ms.GetCounterMetric("cm1")
		assert.Equal(t, ms.CounterMetrics["cm1"], val)
	})

	t.Run("get gauge metric", func(t *testing.T) {
		val, _ := ms.GetGaugeMetric("gm1")
		assert.Equal(t, ms.GaugeMetrics["gm1"], val)
	})

	t.Run("get all counter metrics", func(t *testing.T) {
		vals, _ := ms.GetCounterMetrics()
		assert.EqualValues(t, ms.CounterMetrics, vals)
	})

	t.Run("get all gaugel metrics", func(t *testing.T) {
		vals, _ := ms.GetGaugeMetrics()
		assert.EqualValues(t, ms.GaugeMetrics, vals)
	})
}

func TestMemStorageGetCounterMetric_Error(t *testing.T) {
	ms := NewMemStorage()

	_, err := ms.GetCounterMetric("nonexistent")
	assert.Error(t, err)
	assert.Equal(t, "value not found", err.Error())
}

func TestMemStorageGetGaugeMetric_Error(t *testing.T) {
	ms := NewMemStorage()

	_, err := ms.GetGaugeMetric("nonexistent")
	assert.Error(t, err)
	assert.Equal(t, "value not found", err.Error())
}

func TestMemStorageUpdateMetrics_EmptyInput(t *testing.T) {
	ms := NewMemStorage()

	err := ms.UpdateMetrics([]domain.Metrics{})
	assert.NoError(t, err)
}
