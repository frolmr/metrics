package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/frolmr/metrics/internal/server/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggingMiddleware(t *testing.T) {
	core, recorded := observer.New(zapcore.InfoLevel)
	mockLogger := &logger.Logger{
		SugaredLogger: *zap.New(core).Sugar(),
	}

	r := chi.NewRouter()
	r.Use(WithLog(mockLogger))
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})
	r.Get("/error", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error response"))
	})

	tests := []struct {
		name         string
		path         string
		expectedCode int
		expectedSize int
	}{
		{
			name:         "successful request",
			path:         "/test",
			expectedCode: http.StatusOK,
			expectedSize: len("test response"),
		},
		{
			name:         "error request",
			path:         "/error",
			expectedCode: http.StatusInternalServerError,
			expectedSize: len("error response"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			rec := httptest.NewRecorder()

			recorded.TakeAll()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, rec.Code)
			}

			logs := recorded.All()
			if len(logs) != 1 {
				t.Fatalf("Expected 1 log message, got %d", len(logs))
			}

			log := logs[0]
			if log.Level != zapcore.InfoLevel {
				t.Errorf("Expected log level Info, got %v", log.Level)
			}

			expectedFields := map[string]interface{}{
				"uri":      tt.path,
				"method":   "GET",
				"status":   tt.expectedCode,
				"size":     tt.expectedSize,
				"duration": time.Duration(0),
			}

			for _, field := range log.Context {
				switch field.Key {
				case "uri":
					if field.String != expectedFields["uri"] {
						t.Errorf("Expected uri %v, got %v", expectedFields["uri"], field.String)
					}
				case "method":
					if field.String != expectedFields["method"] {
						t.Errorf("Expected method %v, got %v", expectedFields["method"], field.String)
					}
				case "status":
					if field.Integer != int64(expectedFields["status"].(int)) {
						t.Errorf("Expected status %v, got %v", expectedFields["status"], field.Integer)
					}
				case "size":
					if field.Integer != int64(expectedFields["size"].(int)) {
						t.Errorf("Expected size %v, got %v", expectedFields["size"], field.Integer)
					}
				case "duration":
					if field.Type != zapcore.DurationType {
						t.Error("Expected duration field to be of DurationType")
					}
				}
			}
		})
	}
}

func TestLoggingResponseWriter(t *testing.T) {
	t.Run("Write sets default status if not set", func(t *testing.T) {
		rec := httptest.NewRecorder()
		lrw := &loggingResponseWriter{
			ResponseWriter: rec,
			responseData:   &responseData{},
		}

		data := []byte("test")
		size, err := lrw.Write(data)
		if err != nil {
			t.Fatalf("Write failed: %v", err)
		}

		if size != len(data) {
			t.Errorf("Expected size %d, got %d", len(data), size)
		}
		if lrw.responseData.status != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, lrw.responseData.status)
		}
		if lrw.responseData.size != len(data) {
			t.Errorf("Expected size %d, got %d", len(data), lrw.responseData.size)
		}
		if !lrw.wroteHeader {
			t.Error("Expected wroteHeader to be true after Write")
		}
	})

	t.Run("WriteHeader captures status code", func(t *testing.T) {
		rec := httptest.NewRecorder()
		lrw := &loggingResponseWriter{
			ResponseWriter: rec,
			responseData:   &responseData{},
		}

		testStatus := http.StatusNotFound
		lrw.WriteHeader(testStatus)

		if lrw.responseData.status != testStatus {
			t.Errorf("Expected status %d, got %d", testStatus, lrw.responseData.status)
		}
		if !lrw.wroteHeader {
			t.Error("Expected wroteHeader to be true")
		}
	})

	t.Run("Write after WriteHeader doesn't change status", func(t *testing.T) {
		rec := httptest.NewRecorder()
		lrw := &loggingResponseWriter{
			ResponseWriter: rec,
			responseData:   &responseData{},
		}

		testStatus := http.StatusForbidden
		lrw.WriteHeader(testStatus)
		_, _ = lrw.Write([]byte("test"))

		if lrw.responseData.status != testStatus {
			t.Errorf("Expected status to remain %d, got %d", testStatus, lrw.responseData.status)
		}
		if lrw.responseData.size != len("test") {
			t.Errorf("Expected size %d, got %d", len("test"), lrw.responseData.size)
		}
	})
}

func TestMiddlewareWithEmptyResponse(t *testing.T) {
	core, recorded := observer.New(zapcore.InfoLevel)
	mockLogger := &logger.Logger{
		SugaredLogger: *zap.New(core).Sugar(),
	}

	r := chi.NewRouter()
	r.Use(WithLog(mockLogger))
	r.Get("/empty", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest("GET", "/empty", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	logs := recorded.All()
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log message, got %d", len(logs))
	}

	for _, field := range logs[0].Context {
		switch field.Key {
		case "status":
			if field.Integer != http.StatusNoContent {
				t.Errorf("Expected status %d, got %d", http.StatusNoContent, field.Integer)
			}
		case "size":
			if field.Integer != 0 {
				t.Errorf("Expected size 0, got %d", field.Integer)
			}
		}
	}
}
