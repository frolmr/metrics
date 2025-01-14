package storage

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/frolmr/metrics.git/internal/domain"
)

func TestPing(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPing()

	dbstor := NewDBStorage(db)

	if err := dbstor.Ping(); err != nil {
		t.Errorf("error was not expected while ping: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateCounterMetric(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO counter_metrics")
	mock.ExpectExec("INSERT INTO counter_metrics").WithArgs("test", 1).WillReturnResult(sqlmock.NewResult(1, 1))

	dbstor := NewDBStorage(db)

	if err := dbstor.UpdateCounterMetric("test", 1); err != nil {
		t.Errorf("error was not expected while updating counter metrics: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateGaugeMetric(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO gauge_metrics")
	mock.ExpectExec("INSERT INTO gauge_metrics").WithArgs("test", 1.1).WillReturnResult(sqlmock.NewResult(1, 1))

	dbstor := NewDBStorage(db)

	if err := dbstor.UpdateGaugeMetric("test", 1.1); err != nil {
		t.Errorf("error was not expected while updating gauge metrics: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateMetrics(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbstor := NewDBStorage(db)
	fValue := float64(1.1)
	iValue := int64(11)

	metrics := []domain.Metrics{
		{
			ID:    "tg",
			MType: "gauge",
			Delta: nil,
			Value: &fValue,
		},
		{
			ID:    "tc",
			MType: "counter",
			Delta: &iValue,
			Value: nil,
		},
	}

	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO gauge_metrics")
	mock.ExpectPrepare("INSERT INTO gauge_metrics")
	mock.ExpectExec("INSERT INTO gauge_metrics").WithArgs("tg", 1.1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectPrepare("INSERT INTO counter_metrics")
	mock.ExpectPrepare("INSERT INTO counter_metrics")
	mock.ExpectExec("INSERT INTO counter_metrics").WithArgs("tc", 11).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := dbstor.UpdateMetrics(metrics); err != nil {
		t.Errorf("error was not expected while updating gauge metrics: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetCounterMetric(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	metricsMockRows := sqlmock.NewRows([]string{"value"}).AddRow("1")

	mock.ExpectPrepare("SELECT value FROM counter_metrics")
	mock.ExpectQuery("SELECT value FROM counter_metrics").WithArgs("test").WillReturnRows(metricsMockRows)

	dbstor := NewDBStorage(db)

	if _, err := dbstor.GetCounterMetric("test"); err != nil {
		t.Errorf("error was not expected while getting counter metric: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetCounterMetrics(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	metricsMockRows := sqlmock.NewRows([]string{"value", "name"}).AddRow("1", "test")

	mock.ExpectPrepare("SELECT name, value FROM counter_metrics")
	mock.ExpectQuery("SELECT name, value FROM counter_metrics").WillReturnRows(metricsMockRows)

	dbstor := NewDBStorage(db)

	_ = dbstor.GetCounterMetrics()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
