package storage

import (
	"database/sql"

	"github.com/frolmr/metrics.git/internal/domain"
)

const insertBatchSize = 100

type DBStorage struct {
	db *sql.DB
}

type Number interface {
	int64 | float64
}

func NewDBStorage(db *sql.DB) *DBStorage {
	return &DBStorage{
		db: db,
	}
}

func (ds DBStorage) Ping() error {
	if err := ds.db.Ping(); err != nil {
		return err
	}
	return nil
}

func (ds DBStorage) UpdateCounterMetric(name string, value int64) error {
	stmt, err := ds.insertCounterMetricStatement()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(name, value)
	if err != nil {
		return err
	}

	return nil
}

func (ds DBStorage) UpdateGaugeMetric(name string, value float64) error {
	stmt, err := ds.insertGaugeMetricStatement()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(name, value)
	if err != nil {
		return err
	}

	return nil
}

func (ds DBStorage) splitInGroups(metrics []domain.Metrics) [][]domain.Metrics {
	mLen := len(metrics)
	if mLen <= insertBatchSize {
		return [][]domain.Metrics{metrics}
	}
	var batchedMetrics [][]domain.Metrics
	for i := 0; i < mLen; i += insertBatchSize {
		end := i + insertBatchSize
		if end > mLen {
			end = mLen
		}

		batchedMetrics = append(batchedMetrics, metrics[i:end])
	}
	return batchedMetrics
}

func (ds DBStorage) UpdateMetrics(metrics []domain.Metrics) error {
	metricsGroups := ds.splitInGroups(metrics)

	counterStmt, err := ds.insertGaugeMetricStatement()
	if err != nil {
		return err
	}

	gaugeStmt, err := ds.insertCounterMetricStatement()
	if err != nil {
		return err
	}

	for _, group := range metricsGroups {
		tx, err := ds.db.Begin()
		if err != nil {
			return err
		}

		for _, m := range group {
			if m.MType == domain.CounterType {
				_, err = tx.Stmt(gaugeStmt).Exec(m.ID, *m.Delta)
				if err != nil {
					_ = tx.Rollback()
					return err
				}
			} else {
				_, err = tx.Stmt(counterStmt).Exec(m.ID, *m.Value)
				if err != nil {
					_ = tx.Rollback()
					return err
				}
			}
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func (ds DBStorage) insertCounterMetricStatement() (*sql.Stmt, error) {
	queryString := "INSERT INTO counter_metrics(name, value) " +
		"VALUES ($1, $2) " +
		"ON CONFLICT (name) DO UPDATE SET value = counter_metrics.value + $2"

	return ds.db.Prepare(queryString)
}

func (ds DBStorage) insertGaugeMetricStatement() (*sql.Stmt, error) {
	queryString := "INSERT INTO gauge_metrics(name, value)" +
		"VALUES ($1, $2) " +
		"ON CONFLICT (name) DO UPDATE SET value = $2"

	return ds.db.Prepare(queryString)
}

func (ds DBStorage) GetCounterMetric(name string) (int64, error) {
	stmt, err := ds.db.Prepare("SELECT value FROM counter_metrics WHERE name = $1")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var val int64
	err = stmt.QueryRow(name).Scan(&val)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (ds DBStorage) GetGaugeMetric(name string) (float64, error) {
	stmt, err := ds.db.Prepare("SELECT value FROM gauge_metrics WHERE name = $1")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var val float64
	err = stmt.QueryRow(name).Scan(&val)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (ds DBStorage) GetCounterMetrics() (map[string]int64, error) {
	vals := make(map[string]int64, 0)
	stmt, err := ds.db.Prepare("SELECT name, value FROM counter_metrics")
	if err != nil {
		return nil, err
	}

	return getMetrics(stmt, vals), nil
}

func (ds DBStorage) GetGaugeMetrics() (map[string]float64, error) {
	vals := make(map[string]float64, 0)
	stmt, err := ds.db.Prepare("SELECT name, value FROM gauge_metrics")
	if err != nil {
		return nil, err
	}

	return getMetrics(stmt, vals), nil
}

func getMetrics[K string, V Number](stmt *sql.Stmt, m map[string]V) map[string]V {
	rows, err := stmt.Query()
	if err != nil {
		return nil
	}
	defer rows.Close()
	var (
		name  string
		value V
	)
	for rows.Next() {
		err := rows.Scan(&name, &value)
		if err != nil {
			return nil
		}
		m[name] = value
	}
	err = rows.Err()
	if err != nil {
		return nil
	}
	defer rows.Close()

	return m
}
