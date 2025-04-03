package metrics

import (
	"testing"
	"time"

	"github.com/shirou/gopsutil/mem"
	"github.com/stretchr/testify/assert"
)

func TestMetricsCollection(t *testing.T) {
	t.Run("CollectMetrics", func(t *testing.T) {
		mc := &MetricsCollection{
			CounterMetrics: make(map[string]int64),
			GaugeMetrics:   make(map[string]float64),
		}
		pollCount = 0

		counter, gauge := mc.CollectMetrics()

		assert.Equal(t, int64(1), counter["PollCount"])
		assert.Contains(t, gauge, "Alloc")
		assert.Contains(t, gauge, "BuckHashSys")
		assert.Contains(t, gauge, "Frees")
		assert.Contains(t, gauge, "GCCPUFraction")
		assert.Contains(t, gauge, "GCSys")
		assert.Contains(t, gauge, "HeapAlloc")
		assert.Contains(t, gauge, "HeapIdle")
		assert.Contains(t, gauge, "HeapInuse")
		assert.Contains(t, gauge, "HeapObjects")
		assert.Contains(t, gauge, "HeapReleased")
		assert.Contains(t, gauge, "HeapSys")
		assert.Contains(t, gauge, "LastGC")
		assert.Contains(t, gauge, "Lookups")
		assert.Contains(t, gauge, "MCacheInuse")
		assert.Contains(t, gauge, "MCacheSys")
		assert.Contains(t, gauge, "MSpanInuse")
		assert.Contains(t, gauge, "MSpanSys")
		assert.Contains(t, gauge, "Mallocs")
		assert.Contains(t, gauge, "NextGC")
		assert.Contains(t, gauge, "NumForcedGC")
		assert.Contains(t, gauge, "NumGC")
		assert.Contains(t, gauge, "OtherSys")
		assert.Contains(t, gauge, "PauseTotalNs")
		assert.Contains(t, gauge, "StackInuse")
		assert.Contains(t, gauge, "StackSys")
		assert.Contains(t, gauge, "Sys")
		assert.Contains(t, gauge, "TotalAlloc")
		assert.NotContains(t, gauge, "SomeOtherMetric")
	})

	t.Run("CollectAdditionalMetrics", func(t *testing.T) {
		mc := &MetricsCollection{
			CounterMetrics: make(map[string]int64),
			GaugeMetrics:   make(map[string]float64),
		}

		origMemVirtualMemory := memVirtualMemory
		origCPUPercent := cpuPercent

		memVirtualMemory = func() (*mem.VirtualMemoryStat, error) {
			return &mem.VirtualMemoryStat{
				Total: 100,
				Free:  50,
			}, nil
		}

		defer func() { memVirtualMemory = origMemVirtualMemory }()

		cpuPercent = func(interval time.Duration, percpu bool) ([]float64, error) {
			return []float64{10.5, 20.3, 30.7}, nil
		}

		defer func() { cpuPercent = origCPUPercent }()

		mc.CollectAdditionalMetrics()

		assert.Equal(t, float64(100), mc.GaugeMetrics["TotalMemory"])
		assert.Equal(t, float64(50), mc.GaugeMetrics["FreeMemory"])
		assert.Equal(t, 10.5, mc.GaugeMetrics["CPUutilization0"])
		assert.Equal(t, 20.3, mc.GaugeMetrics["CPUutilization1"])
		assert.Equal(t, 30.7, mc.GaugeMetrics["CPUutilization2"])
	})

	t.Run("CollectAdditionalMetrics_ErrorHandling", func(t *testing.T) {
		mc := &MetricsCollection{
			CounterMetrics: make(map[string]int64),
			GaugeMetrics:   make(map[string]float64),
		}

		origMemVirtualMemory := memVirtualMemory
		origCPUPercent := cpuPercent

		memVirtualMemory = func() (*mem.VirtualMemoryStat, error) {
			return nil, assert.AnError
		}
		defer func() { memVirtualMemory = origMemVirtualMemory }()

		cpuPercent = func(interval time.Duration, percpu bool) ([]float64, error) {
			return nil, assert.AnError
		}
		defer func() { cpuPercent = origCPUPercent }()

		mc.CollectAdditionalMetrics()

		_, ok := mc.GaugeMetrics["TotalMemory"]
		assert.False(t, ok)
		_, ok = mc.GaugeMetrics["FreeMemory"]
		assert.False(t, ok)
	})
}
