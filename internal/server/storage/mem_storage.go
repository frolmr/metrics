package storage

import (
	"errors"

	"github.com/frolmr/metrics.git/internal/domain"
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

func (ms MemStorage) Ping() error {
	return nil
}

func (ms MemStorage) UpdateCounterMetric(name string, value int64) error {
	ms.CounterMetrics[name] += value
	return nil
}

func (ms MemStorage) UpdateGaugeMetric(name string, value float64) error {
	ms.GaugeMetrics[name] = value
	return nil
}

func (ms MemStorage) UpdateMetrics(metrics []domain.Metrics) error {
	for _, v := range metrics {
		if v.MType == domain.CounterType {
			if err := ms.UpdateCounterMetric(v.ID, *v.Delta); err != nil {
				return err
			}
		} else {
			if err := ms.UpdateGaugeMetric(v.ID, *v.Value); err != nil {
				return err
			}
		}
	}
	return nil
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

func (ms MemStorage) GetCounterMetrics() (map[string]int64, error) {
	return ms.CounterMetrics, nil
}

func (ms MemStorage) GetGaugeMetrics() (map[string]float64, error) {
	return ms.GaugeMetrics, nil
}
