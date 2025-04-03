package storage

import (
	"errors"
	"time"

	"github.com/frolmr/metrics/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

type RetriableStorage struct {
	dbStorage      *DBStorage
	retryIntervals []time.Duration
}

func NewRetriableStorage(dbStorage *DBStorage) *RetriableStorage {
	return &RetriableStorage{
		dbStorage:      dbStorage,
		retryIntervals: []time.Duration{time.Second, time.Second * 2, time.Second * 5},
	}
}

func (rs RetriableStorage) Ping() (err error) {
	for _, interval := range rs.retryIntervals {
		err = rs.dbStorage.Ping()
		if err == nil {
			return nil
		}
		if rs.isRetriable(err) {
			time.Sleep(interval)
		}
	}
	return
}

func (rs RetriableStorage) UpdateCounterMetric(name string, value int64) (err error) {
	for _, interval := range rs.retryIntervals {
		err = rs.dbStorage.UpdateCounterMetric(name, value)
		if err == nil {
			return nil
		}
		if rs.isRetriable(err) {
			time.Sleep(interval)
		}
	}
	return
}

func (rs RetriableStorage) UpdateGaugeMetric(name string, value float64) (err error) {
	for _, interval := range rs.retryIntervals {
		err = rs.dbStorage.UpdateGaugeMetric(name, value)
		if err == nil {
			return nil
		}
		if rs.isRetriable(err) {
			time.Sleep(interval)
		}
	}
	return
}

func (rs RetriableStorage) UpdateMetrics(metrics []domain.Metrics) (err error) {
	for _, interval := range rs.retryIntervals {
		err = rs.dbStorage.UpdateMetrics(metrics)
		if err == nil {
			return nil
		}
		if rs.isRetriable(err) {
			time.Sleep(interval)
		}
	}
	return
}

func (rs RetriableStorage) GetCounterMetric(name string) (res int64, err error) {
	for _, interval := range rs.retryIntervals {
		res, err = rs.dbStorage.GetCounterMetric(name)
		if err == nil {
			return res, nil
		}
		if rs.isRetriable(err) {
			time.Sleep(interval)
		}
	}
	return
}

func (rs RetriableStorage) GetGaugeMetric(name string) (res float64, err error) {
	for _, interval := range rs.retryIntervals {
		res, err = rs.dbStorage.GetGaugeMetric(name)
		if err == nil {
			return res, nil
		}
		if rs.isRetriable(err) {
			time.Sleep(interval)
		}
	}
	return
}

func (rs RetriableStorage) GetCounterMetrics() (res map[string]int64, err error) {
	for _, interval := range rs.retryIntervals {
		res, err = rs.dbStorage.GetCounterMetrics()
		if err == nil {
			return res, nil
		}
		if rs.isRetriable(err) {
			time.Sleep(interval)
		}
	}
	return
}

func (rs RetriableStorage) GetGaugeMetrics() (res map[string]float64, err error) {
	for _, interval := range rs.retryIntervals {
		res, err = rs.dbStorage.GetGaugeMetrics()
		if err == nil {
			return res, nil
		}
		if rs.isRetriable(err) {
			time.Sleep(interval)
		}
	}
	return
}

func (rs RetriableStorage) isRetriable(err error) bool {
	switch {
	case errors.Is(err, &pgconn.ConnectError{}):
		return true
	default:
		return false
	}
}
