package storage

import (
	"database/sql"
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
	_, err := ds.GetCounterMetric(name)
	if err != nil {
		_, err := ds.db.Exec("INSERT INTO counter_metrics(name, value) VALUES ($1, $2)", name, value)
		if err != nil {
			return err
		}
	} else {
		_, err := ds.db.Exec("UPDATE counter_metrics SET value = $1 WHERE name = $2", value, name)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ds DBStorage) UpdateGaugeMetric(name string, value float64) error {
	_, err := ds.GetGaugeMetric(name)
	if err != nil {
		_, err := ds.db.Exec("INSERT INTO gauge_metrics(name, value) VALUES ($1, $2)", name, value)
		if err != nil {
			return err
		}
	} else {
		_, err := ds.db.Exec("UPDATE gauge_metrics SET value = $1 WHERE name = $2", value, name)
		if err != nil {
			return err
		}
	}

	return nil
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
