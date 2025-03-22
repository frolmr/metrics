package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/frolmr/metrics.git/internal/domain"
	"github.com/frolmr/metrics.git/internal/server/mocks"
	"github.com/frolmr/metrics.git/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func testJSONRequest(t *testing.T, ts *httptest.Server, method,
	path string, body []byte) ([]byte, int) {
	//nolint:noctx // Context will be added later
	req, err := http.NewRequest(method, ts.URL+path, bytes.NewBuffer(body))
	require.NoError(t, err)

	req.Header.Set("content-type", domain.JSONContentType)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return respBody, resp.StatusCode
}

func prepareBodySlice(metrics []domain.Metrics) []byte {
	result, _ := json.Marshal(metrics)
	return result
}

func TestUpdateJSONMetricHandler(t *testing.T) {
	ms := storage.MemStorage{
		CounterMetrics: make(map[string]int64),
		GaugeMetrics:   make(map[string]float64),
	}

	rh := NewRequestHandler(ms)

	r := chi.NewRouter()
	r.Post("/update", rh.UpdateMetricJSON())

	ts := httptest.NewServer(r)
	defer ts.Close()

	gaugeVal := 1.1
	var counterVal int64 = 1

	type want struct {
		statusCode   int
		responseBody []byte
	}
	tests := []struct {
		name   string
		method string
		body   []byte
		want   want
	}{
		{
			name:   "success gauge request",
			method: http.MethodPost,
			body:   prepareBody(domain.Metrics{ID: "tstGauge", MType: "gauge", Value: &gaugeVal}),
			want: want{
				statusCode:   http.StatusOK,
				responseBody: prepareBody(domain.Metrics{ID: "tstGauge", MType: "gauge", Value: &gaugeVal}),
			},
		},
		{
			name:   "success counter request",
			method: http.MethodPost,
			body:   prepareBody(domain.Metrics{ID: "tstGauge", MType: "counter", Delta: &counterVal}),
			want: want{
				statusCode:   http.StatusOK,
				responseBody: prepareBody(domain.Metrics{ID: "tstGauge", MType: "counter", Delta: &counterVal}),
			},
		},
		{
			name:   "fail wrong method request",
			method: http.MethodGet,
			body:   prepareBody(domain.Metrics{ID: "tstGauge", MType: "gauge", Value: &gaugeVal}),
			want: want{
				statusCode:   http.StatusMethodNotAllowed,
				responseBody: []byte(""),
			},
		},
		{
			name:   "fail wrong metric type request",
			method: http.MethodPost,
			body:   prepareBody(domain.Metrics{ID: "tstGauge", MType: "invalid_type", Value: &gaugeVal}),
			want: want{
				statusCode:   http.StatusBadRequest,
				responseBody: []byte("wrong metric type\n"),
			},
		},
	}

	for _, tt := range tests {
		body, code := testJSONRequest(t, ts, tt.method, "/update", tt.body)
		assert.Equal(t, string(tt.want.responseBody), string(body))
		assert.Equal(t, tt.want.statusCode, code)
	}
}

func TestGetJSONMetricHandler(t *testing.T) {
	gaugeVal := 1.1
	var counterVal int64 = 1

	ms := storage.MemStorage{
		CounterMetrics: map[string]int64{"cTest1": counterVal},
		GaugeMetrics:   map[string]float64{"gTest1": gaugeVal},
	}

	rh := NewRequestHandler(ms)

	r := chi.NewRouter()
	r.Post("/value", rh.GetMetricJSON())

	ts := httptest.NewServer(r)
	defer ts.Close()

	type want struct {
		statusCode   int
		responseBody []byte
	}
	tests := []struct {
		name   string
		method string
		body   []byte
		want   want
	}{
		{
			name:   "success gauge request",
			method: http.MethodPost,
			body:   prepareBody(domain.Metrics{ID: "gTest1", MType: "gauge"}),
			want: want{
				statusCode:   http.StatusOK,
				responseBody: prepareBody(domain.Metrics{ID: "gTest1", MType: "gauge", Value: &gaugeVal}),
			},
		},
		{
			name:   "success counter request",
			method: http.MethodPost,
			body:   prepareBody(domain.Metrics{ID: "cTest1", MType: "counter"}),
			want: want{
				statusCode:   http.StatusOK,
				responseBody: prepareBody(domain.Metrics{ID: "cTest1", MType: "counter", Delta: &counterVal}),
			},
		},
		{
			name:   "fail wrong method request",
			method: http.MethodGet,
			body:   prepareBody(domain.Metrics{ID: "tstGauge", MType: "gauge"}),
			want: want{
				statusCode:   http.StatusMethodNotAllowed,
				responseBody: []byte(""),
			},
		},
		{
			name:   "fail wrong metric type request",
			method: http.MethodPost,
			body:   prepareBody(domain.Metrics{ID: "tstGauge", MType: "invalid_type"}),
			want: want{
				statusCode:   http.StatusBadRequest,
				responseBody: []byte("wrong metric type\n"),
			},
		},
	}

	for _, tt := range tests {
		body, code := testJSONRequest(t, ts, tt.method, "/value", tt.body)
		assert.Equal(t, string(tt.want.responseBody), string(body))
		assert.Equal(t, tt.want.statusCode, code)
	}
}

func prepareBody(metrics domain.Metrics) []byte {
	result, _ := json.Marshal(metrics)

	return result
}

func TestBulkUpdateJSONMetricHandler(t *testing.T) {
	ms := storage.MemStorage{
		CounterMetrics: make(map[string]int64),
		GaugeMetrics:   make(map[string]float64),
	}

	rh := NewRequestHandler(ms)

	r := chi.NewRouter()
	r.Post("/updates", rh.BulkUpdateMetricJSON())

	ts := httptest.NewServer(r)
	defer ts.Close()

	gaugeVal := 1.1
	var counterVal int64 = 1

	type want struct {
		statusCode   int
		responseBody []byte
	}
	tests := []struct {
		name   string
		method string
		body   []byte
		want   want
	}{
		{
			name:   "success bulk update request",
			method: http.MethodPost,
			body: prepareBodySlice([]domain.Metrics{
				{ID: "tstGauge", MType: "gauge", Value: &gaugeVal},
				{ID: "tstCounter", MType: "counter", Delta: &counterVal},
			}),
			want: want{
				statusCode:   http.StatusOK,
				responseBody: []byte(""),
			},
		},
		{
			name:   "fail invalid payload",
			method: http.MethodPost,
			body:   []byte("invalid json"),
			want: want{
				statusCode:   http.StatusBadRequest,
				responseBody: []byte("invalid character 'i' looking for beginning of value\n"),
			},
		},
		{
			name:   "fail wrong method request",
			method: http.MethodGet,
			body: prepareBodySlice([]domain.Metrics{
				{ID: "tstGauge", MType: "gauge", Value: &gaugeVal},
			}),
			want: want{
				statusCode:   http.StatusMethodNotAllowed,
				responseBody: []byte(""),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, code := testJSONRequest(t, ts, tt.method, "/updates", tt.body)
			assert.Equal(t, string(tt.want.responseBody), string(body))
			assert.Equal(t, tt.want.statusCode, code)
		})
	}
}

func TestUpdateJSONMetricHandler_ErrorScenarios(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock repository
	mockRepo := mocks.NewMockRepository(ctrl)

	// Set up expectations for the mock
	mockRepo.EXPECT().
		UpdateGaugeMetric(gomock.Any(), gomock.Any()).
		Return(errors.New("repo error")).
		Times(1)

	// Create the handler with the mock repository
	rh := NewRequestHandler(mockRepo)

	r := chi.NewRouter()
	r.Post("/update", rh.UpdateMetricJSON())

	ts := httptest.NewServer(r)
	defer ts.Close()

	gaugeVal := 1.1

	type want struct {
		statusCode   int
		responseBody []byte
	}
	tests := []struct {
		name   string
		method string
		body   []byte
		want   want
	}{
		{
			name:   "fail repository error for gauge",
			method: http.MethodPost,
			body:   prepareBody(domain.Metrics{ID: "tstGauge", MType: "gauge", Value: &gaugeVal}),
			want: want{
				statusCode:   http.StatusBadRequest,
				responseBody: []byte("error updating metric\n"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, code := testJSONRequest(t, ts, tt.method, "/update", tt.body)
			assert.Equal(t, string(tt.want.responseBody), string(body))
			assert.Equal(t, tt.want.statusCode, code)
		})
	}
}

func TestGetJSONMetricHandler_ErrorScenarios(t *testing.T) {
	ms := storage.MemStorage{
		CounterMetrics: make(map[string]int64),
		GaugeMetrics:   make(map[string]float64),
	}

	rh := NewRequestHandler(ms)

	r := chi.NewRouter()
	r.Post("/value", rh.GetMetricJSON())

	ts := httptest.NewServer(r)
	defer ts.Close()

	type want struct {
		statusCode   int
		responseBody []byte
	}
	tests := []struct {
		name   string
		method string
		body   []byte
		want   want
	}{
		{
			name:   "fail invalid payload",
			method: http.MethodPost,
			body:   []byte("invalid json"),
			want: want{
				statusCode:   http.StatusBadRequest,
				responseBody: []byte("invalid character 'i' looking for beginning of value\n"),
			},
		},
		{
			name:   "fail metric not found",
			method: http.MethodPost,
			body:   prepareBody(domain.Metrics{ID: "nonexistent", MType: "gauge"}),
			want: want{
				statusCode:   http.StatusNotFound,
				responseBody: []byte("metric not found\n"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, code := testJSONRequest(t, ts, tt.method, "/value", tt.body)
			assert.Equal(t, string(tt.want.responseBody), string(body))
			assert.Equal(t, tt.want.statusCode, code)
		})
	}
}
