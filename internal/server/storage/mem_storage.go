package storage

import (
	"errors"
)

type MemStorage struct {
	CounterMetrics map[string]int64
	GaugeMetrics   map[string]float64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		CounterMetrics: make(map[string]int64),
		GaugeMetrics:   make(map[string]float64),
	}
}

func (ms MemStorage) UpdateCounterMetric(name string, value int64) {
	ms.CounterMetrics[name] += value
}

func (ms MemStorage) UpdateGaugeMetric(name string, value float64) {
	ms.GaugeMetrics[name] = value
}

func (ms MemStorage) GetCounterMetric(name string) (int64, error) {
	if value, exists := ms.CounterMetrics[name]; !exists {
		return 0, errors.New("value not found")
	} else {
		return value, nil
	}
}

func (ms MemStorage) GetGaugeMetric(name string) (float64, error) {
	if value, exists := ms.GaugeMetrics[name]; !exists {
		return 0, errors.New("value not found")
	} else {
		return value, nil
	}
}

func (ms MemStorage) GetCounterMetrics() map[string]int64 {
	return ms.CounterMetrics
}

func (ms MemStorage) GetGaugeMetrics() map[string]float64 {
	return ms.GaugeMetrics
}
