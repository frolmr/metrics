package storage

import (
	"database/sql"

	"github.com/frolmr/metrics.git/internal/domain"
)

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

func (ds DBStorage) UpdateMetrics(metrics []domain.Metrics) error {
	tx, err := ds.db.Begin()
	if err != nil {
		return err
	}

	for _, m := range metrics {
		if m.MType == domain.CounterType {
			stmt, err := ds.insertCounterMetricStatement()
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			_, err = tx.Stmt(stmt).Exec(m.ID, *m.Delta)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
		} else {
			stmt, err := ds.insertGaugeMetricStatement()
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			_, err = tx.Stmt(stmt).Exec(m.ID, *m.Value)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit()
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

func (ds DBStorage) GetCounterMetrics() map[string]int64 {
	vals := make(map[string]int64, 0)
	stmt, err := ds.db.Prepare("SELECT name, value FROM counter_metrics")
	if err != nil {
		return nil
	}

	return getMetrics(stmt, vals)
}

func (ds DBStorage) GetGaugeMetrics() map[string]float64 {
	vals := make(map[string]float64, 0)
	stmt, err := ds.db.Prepare("SELECT name, value FROM gauge_metrics")
	if err != nil {
		return nil
	}

	return getMetrics(stmt, vals)
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
