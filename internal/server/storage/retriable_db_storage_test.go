package storage

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/frolmr/metrics/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestRetriableStorage_Ping_Success(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPing()

	dbStorage := NewDBStorage(db)
	retriableStorage := NewRetriableStorage(dbStorage)

	err = retriableStorage.Ping()
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetriableStorage_Ping_RetrySuccess(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPing().WillReturnError(&pgconn.ConnectError{})
	mock.ExpectPing()

	dbStorage := NewDBStorage(db)
	retriableStorage := NewRetriableStorage(dbStorage)

	start := time.Now()
	err = retriableStorage.Ping()
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.True(t, duration >= time.Second, "should have waited at least 1 second")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetriableStorage_Ping_NonRetriableError(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPing().WillReturnError(errors.New("non-retriable error"))

	dbStorage := NewDBStorage(db)
	retriableStorage := NewRetriableStorage(dbStorage)

	start := time.Now()
	err = retriableStorage.Ping()
	duration := time.Since(start)

	assert.Error(t, err)
	assert.True(t, duration < time.Second, "should not have waited for retry")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetriableStorage_Ping_AllRetriesFail(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPing().WillReturnError(&pgconn.ConnectError{})
	mock.ExpectPing().WillReturnError(&pgconn.ConnectError{})
	mock.ExpectPing().WillReturnError(&pgconn.ConnectError{})

	dbStorage := NewDBStorage(db)
	retriableStorage := NewRetriableStorage(dbStorage)

	start := time.Now()
	err = retriableStorage.Ping()
	duration := time.Since(start)

	assert.Error(t, err)
	assert.True(t, duration >= 8*time.Second, "should have waited for all retries (1+2+5 seconds)")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetriableStorage_UpdateCounterMetric_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO counter_metrics")
	mock.ExpectExec("INSERT INTO counter_metrics").WithArgs("test", 1).WillReturnResult(sqlmock.NewResult(1, 1))

	dbStorage := NewDBStorage(db)
	retriableStorage := NewRetriableStorage(dbStorage)

	err = retriableStorage.UpdateCounterMetric("test", 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetriableStorage_UpdateCounterMetric_RetrySuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO counter_metrics").WillReturnError(&pgconn.ConnectError{})
	mock.ExpectPrepare("INSERT INTO counter_metrics")
	mock.ExpectExec("INSERT INTO counter_metrics").WithArgs("test", 1).WillReturnResult(sqlmock.NewResult(1, 1))

	dbStorage := NewDBStorage(db)
	retriableStorage := NewRetriableStorage(dbStorage)

	start := time.Now()
	err = retriableStorage.UpdateCounterMetric("test", 1)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.True(t, duration >= time.Second, "should have waited at least 1 second")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetriableStorage_UpdateCounterMetric_NonRetriableError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO counter_metrics").WillReturnError(errors.New("non-retriable error"))

	dbStorage := NewDBStorage(db)
	retriableStorage := NewRetriableStorage(dbStorage)

	start := time.Now()
	err = retriableStorage.UpdateCounterMetric("test", 1)
	duration := time.Since(start)

	assert.Error(t, err)
	assert.True(t, duration < time.Second, "should not have waited for retry")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetriableStorage_UpdateGaugeMetric_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("INSERT INTO gauge_metrics")
	mock.ExpectExec("INSERT INTO gauge_metrics").WithArgs("test", 1.1).WillReturnResult(sqlmock.NewResult(1, 1))

	dbStorage := NewDBStorage(db)
	retriableStorage := NewRetriableStorage(dbStorage)

	err = retriableStorage.UpdateGaugeMetric("test", 1.1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetriableStorage_UpdateMetrics_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	fValue := float64(1.1)
	metrics := []domain.Metrics{
		{
			ID:    "test",
			MType: "gauge",
			Value: &fValue,
		},
	}

	mock.ExpectPrepare("INSERT INTO gauge_metrics")
	mock.ExpectPrepare("INSERT INTO counter_metrics")
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO gauge_metrics").WithArgs("test", 1.1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	dbStorage := NewDBStorage(db)
	retriableStorage := NewRetriableStorage(dbStorage)

	err = retriableStorage.UpdateMetrics(metrics)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetriableStorage_UpdateMetrics_RetrySuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	fValue := float64(1.1)
	metrics := []domain.Metrics{
		{
			ID:    "test",
			MType: "gauge",
			Value: &fValue,
		},
	}

	mock.ExpectPrepare("INSERT INTO gauge_metrics").WillReturnError(&pgconn.ConnectError{})
	mock.ExpectPrepare("INSERT INTO gauge_metrics")
	mock.ExpectPrepare("INSERT INTO counter_metrics")
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO gauge_metrics").WithArgs("test", 1.1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	dbStorage := NewDBStorage(db)
	retriableStorage := NewRetriableStorage(dbStorage)

	start := time.Now()
	err = retriableStorage.UpdateMetrics(metrics)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.True(t, duration >= time.Second, "should have waited at least 1 second")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetriableStorage_GetCounterMetric_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	metricsMockRows := sqlmock.NewRows([]string{"value"}).AddRow("1")

	mock.ExpectPrepare("SELECT value FROM counter_metrics")
	mock.ExpectQuery("SELECT value FROM counter_metrics").WithArgs("test").WillReturnRows(metricsMockRows)

	dbStorage := NewDBStorage(db)
	retriableStorage := NewRetriableStorage(dbStorage)

	val, err := retriableStorage.GetCounterMetric("test")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), val)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetriableStorage_GetCounterMetric_RetrySuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectPrepare("SELECT value FROM counter_metrics").WillReturnError(&pgconn.ConnectError{})
	metricsMockRows := sqlmock.NewRows([]string{"value"}).AddRow("1")
	mock.ExpectPrepare("SELECT value FROM counter_metrics")
	mock.ExpectQuery("SELECT value FROM counter_metrics").WithArgs("test").WillReturnRows(metricsMockRows)

	dbStorage := NewDBStorage(db)
	retriableStorage := NewRetriableStorage(dbStorage)

	start := time.Now()
	val, err := retriableStorage.GetCounterMetric("test")
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), val)
	assert.True(t, duration >= time.Second, "should have waited at least 1 second")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetriableStorage_GetGaugeMetric_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	metricsMockRows := sqlmock.NewRows([]string{"value"}).AddRow("1.1")

	mock.ExpectPrepare("SELECT value FROM gauge_metrics")
	mock.ExpectQuery("SELECT value FROM gauge_metrics").WithArgs("test").WillReturnRows(metricsMockRows)

	dbStorage := NewDBStorage(db)
	retriableStorage := NewRetriableStorage(dbStorage)

	val, err := retriableStorage.GetGaugeMetric("test")
	assert.NoError(t, err)
	assert.Equal(t, 1.1, val)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetriableStorage_GetCounterMetrics_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	metricsMockRows := sqlmock.NewRows([]string{"name", "value"}).AddRow("test", 1)

	mock.ExpectPrepare("SELECT name, value FROM counter_metrics")
	mock.ExpectQuery("SELECT name, value FROM counter_metrics").WillReturnRows(metricsMockRows)

	dbStorage := NewDBStorage(db)
	retriableStorage := NewRetriableStorage(dbStorage)

	metrics, err := retriableStorage.GetCounterMetrics()
	assert.NoError(t, err)
	assert.Equal(t, map[string]int64{"test": 1}, metrics)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetriableStorage_GetGaugeMetrics_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	metricsMockRows := sqlmock.NewRows([]string{"name", "value"}).AddRow("test", 1.1)

	mock.ExpectPrepare("SELECT name, value FROM gauge_metrics")
	mock.ExpectQuery("SELECT name, value FROM gauge_metrics").WillReturnRows(metricsMockRows)

	dbStorage := NewDBStorage(db)
	retriableStorage := NewRetriableStorage(dbStorage)

	metrics, err := retriableStorage.GetGaugeMetrics()
	assert.NoError(t, err)
	assert.Equal(t, map[string]float64{"test": 1.1}, metrics)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetriableStorage_IsRetriable(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbStorage := NewDBStorage(db)
	retriableStorage := NewRetriableStorage(dbStorage)

	connErr := &pgconn.ConnectError{}

	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "ConnectError is retriable",
			err:      connErr,
			expected: true,
		},
		{
			name:     "Other error is not retriable",
			err:      errors.New("some error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := retriableStorage.isRetriable(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
