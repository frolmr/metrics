package storage

import (
	"encoding/json"
	"io"
	"log"

	"github.com/frolmr/metrics/internal/domain"
)

type Snapshooter interface {
	RestoreFromSnapshot(source io.Reader) error
	SaveToSnapshot(destination io.Writer) error
}

func (ms *MemStorage) RestoreFromSnapshot(source io.Reader) error {
	metricsSnap := make([]domain.Metrics, 0)
	if err := json.NewDecoder(source).Decode(&metricsSnap); err != nil {
		return err
	}

	for _, metric := range metricsSnap {
		if metric.MType == domain.GaugeType {
			ms.GaugeMetrics[metric.ID] = *metric.Value
		} else if metric.MType == domain.CounterType {
			ms.CounterMetrics[metric.ID] = *metric.Delta
		} else {
			log.Println("invalid data in snapshot: ", metric.MType)
		}
	}

	return nil
}

func (ms *MemStorage) SaveToSnapshot(destination io.Writer) error {
	var metricsJSON []domain.Metrics

	for name, value := range ms.CounterMetrics {
		metricsJSON = append(metricsJSON, domain.Metrics{ID: name, MType: domain.CounterType, Delta: &value})
	}
	for name, value := range ms.GaugeMetrics {
		metricsJSON = append(metricsJSON, domain.Metrics{ID: name, MType: domain.GaugeType, Value: &value})
	}

	data, err := json.MarshalIndent(metricsJSON, "", " ")
	if err != nil {
		return err
	}

	if _, err := destination.Write(data); err != nil {
		return err
	}

	return nil
}
