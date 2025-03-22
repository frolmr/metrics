package storage

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/frolmr/metrics.git/internal/domain"
	"github.com/stretchr/testify/assert"
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

func TestUpdateCounterMetric_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO counter_metrics").WillReturnError(errors.New("prepare error"))

	dbstor := NewDBStorage(db)

	if err := dbstor.UpdateCounterMetric("test", 1); err == nil {
		t.Error("expected an error, but got nil")
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

func TestUpdateGaugeMetric_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO gauge_metrics").WillReturnError(errors.New("prepare error"))

	dbstor := NewDBStorage(db)

	if err := dbstor.UpdateGaugeMetric("test", 1.1); err == nil {
		t.Error("expected an error, but got nil")
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

	mock.ExpectPrepare("INSERT INTO gauge_metrics")
	mock.ExpectPrepare("INSERT INTO counter_metrics")
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO gauge_metrics").WithArgs("tg", 1.1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO counter_metrics").WithArgs("tc", 11).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := dbstor.UpdateMetrics(metrics); err != nil {
		t.Errorf("error was not expected while updating gauge metrics: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateMetrics_Error(t *testing.T) {
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

	mock.ExpectPrepare("INSERT INTO gauge_metrics")
	mock.ExpectPrepare("INSERT INTO counter_metrics")
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO gauge_metrics").WithArgs("tg", 1.1).WillReturnError(errors.New("exec error"))
	mock.ExpectRollback()

	if err := dbstor.UpdateMetrics(metrics); err == nil {
		t.Error("expected an error, but got nil")
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

func TestGetCounterMetric_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT value FROM counter_metrics").WillReturnError(errors.New("prepare error"))

	dbstor := NewDBStorage(db)

	if _, err := dbstor.GetCounterMetric("test"); err == nil {
		t.Error("expected an error, but got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetGaugeMetric(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	metricsMockRows := sqlmock.NewRows([]string{"value"}).AddRow("1")

	mock.ExpectPrepare("SELECT value FROM gauge_metrics")
	mock.ExpectQuery("SELECT value FROM gauge_metrics").WithArgs("test").WillReturnRows(metricsMockRows)

	dbstor := NewDBStorage(db)

	if _, err := dbstor.GetGaugeMetric("test"); err != nil {
		t.Errorf("error was not expected while getting counter metric: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetGaugeMetric_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT value FROM gauge_metrics").WillReturnError(errors.New("prepare error"))

	dbstor := NewDBStorage(db)

	if _, err := dbstor.GetGaugeMetric("test"); err == nil {
		t.Error("expected an error, but got nil")
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

	_, _ = dbstor.GetCounterMetrics()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetCounterMetrics_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT name, value FROM counter_metrics").WillReturnError(errors.New("prepare error"))

	dbstor := NewDBStorage(db)

	if _, err := dbstor.GetCounterMetrics(); err == nil {
		t.Error("expected an error, but got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetGaugeMetrics(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	metricsMockRows := sqlmock.NewRows([]string{"value", "name"}).AddRow("1", "test")

	mock.ExpectPrepare("SELECT name, value FROM gauge_metrics")
	mock.ExpectQuery("SELECT name, value FROM gauge_metrics").WillReturnRows(metricsMockRows)

	dbstor := NewDBStorage(db)

	_, _ = dbstor.GetGaugeMetrics()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetGaugeMetrics_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT name, value FROM gauge_metrics").WillReturnError(errors.New("prepare error"))

	dbstor := NewDBStorage(db)

	if _, err := dbstor.GetGaugeMetrics(); err == nil {
		t.Error("expected an error, but got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSplitInGroups(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbstor := NewDBStorage(db)

	metrics := []domain.Metrics{}
	groups := dbstor.splitInGroups(metrics)
	assert.Equal(t, 1, len(groups))
	assert.Equal(t, 0, len(groups[0]))

	metrics = []domain.Metrics{{ID: "test1"}, {ID: "test2"}}
	groups = dbstor.splitInGroups(metrics)
	assert.Equal(t, 1, len(groups))
	assert.Equal(t, 2, len(groups[0]))

	metrics = make([]domain.Metrics, insertBatchSize)
	groups = dbstor.splitInGroups(metrics)
	assert.Equal(t, 1, len(groups))
	assert.Equal(t, insertBatchSize, len(groups[0]))

	metrics = make([]domain.Metrics, insertBatchSize+1)
	groups = dbstor.splitInGroups(metrics)
	assert.Equal(t, 2, len(groups))
	assert.Equal(t, insertBatchSize, len(groups[0]))
	assert.Equal(t, 1, len(groups[1]))
}

func TestGetMetrics_Int64(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	metricsMockRows := sqlmock.NewRows([]string{"name", "value"}).AddRow("test1", 1).AddRow("test2", 2)

	mock.ExpectPrepare("SELECT name, value FROM counter_metrics")
	mock.ExpectQuery("SELECT name, value FROM counter_metrics").WillReturnRows(metricsMockRows)

	dbstor := NewDBStorage(db)

	result, err := dbstor.GetCounterMetrics()
	assert.NoError(t, err)
	assert.Equal(t, map[string]int64{"test1": 1, "test2": 2}, result)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetMetrics_Float64(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	metricsMockRows := sqlmock.NewRows([]string{"name", "value"}).AddRow("test1", 1.1).AddRow("test2", 2.2)

	mock.ExpectPrepare("SELECT name, value FROM gauge_metrics")
	mock.ExpectQuery("SELECT name, value FROM gauge_metrics").WillReturnRows(metricsMockRows)

	dbstor := NewDBStorage(db)

	result, err := dbstor.GetGaugeMetrics()
	assert.NoError(t, err)
	assert.Equal(t, map[string]float64{"test1": 1.1, "test2": 2.2}, result)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
