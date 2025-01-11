package storage

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testSnapFile = "test_snap"
)

func TestSnaphots(t *testing.T) {
	ms := MemStorage{
		CounterMetrics: map[string]int64{"cm1": 1},
		GaugeMetrics:   map[string]float64{"gm1": 0.2},
	}

	defer os.Remove(testSnapFile)

	t.Run("save and restore from snapshot", func(r *testing.T) {
		file, _ := os.Open(testSnapFile)
		defer file.Close()
		reader := bufio.NewReader(file)
		writer := bufio.NewWriter(file)

		_ = ms.UpdateCounterMetric("cm1", 7)
		_ = ms.UpdateGaugeMetric("gm1", 8.8)
		_ = ms.SaveToSnapshot(writer)
		_ = ms.RestoreFromSnapshot(reader)
		assert.Equal(t, ms.CounterMetrics["cm1"], int64(8))
		assert.Equal(t, ms.GaugeMetrics["gm1"], float64(8.8))
	})
}
