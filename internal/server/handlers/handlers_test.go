package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/frolmr/metrics/internal/server/mocks"
	"github.com/frolmr/metrics/internal/server/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, contentType string) (string, int) {
	//nolint:noctx // No need for context in tests
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	req.Header.Set("content-type", contentType)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return string(respBody), resp.StatusCode
}

func TestMetricsUpdate(t *testing.T) {
	ms := storage.MemStorage{
		CounterMetrics: make(map[string]int64),
		GaugeMetrics:   make(map[string]float64),
	}

	rh := NewRequestHandler(ms)

	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", rh.UpdateMetric())

	ts := httptest.NewServer(r)
	defer ts.Close()

	type want struct {
		statusCode int
	}
	tests := []struct {
		name        string
		path        string
		method      string
		contentType string
		want        want
	}{
		{
			name:        "success gauge request",
			path:        "/update/gauge/test/25",
			method:      http.MethodPost,
			contentType: "text/plain",
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name:        "success counter request",
			path:        "/update/counter/test/25",
			method:      http.MethodPost,
			contentType: "text/plain",
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name:        "fail wrong method request",
			path:        "/update/counter/test/25",
			method:      http.MethodGet,
			contentType: "text/plain",
			want: want{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
		{
			name:        "fail wrong metric type request",
			path:        "/update/cntr/test/25",
			method:      http.MethodPost,
			contentType: "text/plain",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:        "fail wrong gauge metric value format request",
			path:        "/update/gauge/test/2,5",
			method:      http.MethodPost,
			contentType: "text/plain",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:        "fail wrong counter metric value format request",
			path:        "/update/counter/test/2.5",
			method:      http.MethodPost,
			contentType: "text/plain",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		_, code := testRequest(t, ts, tt.method, tt.path, tt.contentType)
		assert.Equal(t, tt.want.statusCode, code)
	}
}

func TestGetMetricHandler(t *testing.T) {
	ms := storage.MemStorage{
		CounterMetrics: map[string]int64{"cTest1": 200, "cTest2": 128},
		GaugeMetrics:   map[string]float64{"gTest1": 2.12, "gTest2": 0.54},
	}

	rh := NewRequestHandler(ms)

	r := chi.NewRouter()
	r.Get("/value/{type}/{name}", rh.GetMetric())

	ts := httptest.NewServer(r)
	defer ts.Close()

	type want struct {
		statusCode int
		response   string
	}
	tests := []struct {
		name        string
		path        string
		method      string
		contentType string
		want        want
	}{
		{
			name:        "success gauge request",
			path:        "/value/gauge/gTest1",
			method:      http.MethodGet,
			contentType: "text/plain;charset=utf-8",
			want: want{
				statusCode: http.StatusOK,
				response:   "2.12",
			},
		},
		{
			name:        "success counter request",
			path:        "/value/counter/cTest1",
			method:      http.MethodGet,
			contentType: "text/plain;charset=utf-8",
			want: want{
				statusCode: http.StatusOK,
				response:   "200",
			},
		},
		{
			name:        "fail gauge request",
			path:        "/value/gauge/gTest10",
			method:      http.MethodGet,
			contentType: "text/plain;charset=utf-8",
			want: want{
				statusCode: http.StatusNotFound,
				response:   "Metric Not Found\n",
			},
		},
		{
			name:        "fail counter request",
			path:        "/value/counter/cTest10",
			method:      http.MethodGet,
			contentType: "text/plain;charset=utf-8",
			want: want{
				statusCode: http.StatusNotFound,
				response:   "Metric Not Found\n",
			},
		},
		{
			name:        "fail method type",
			path:        "/value/counter/cTest10",
			method:      http.MethodPost,
			contentType: "text/plain;charset=utf-8",
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				response:   "",
			},
		},
		{
			name:        "fail metric type request",
			path:        "/value/cntr/cTest1",
			method:      http.MethodGet,
			contentType: "text/plain;charset=utf-8",
			want: want{
				statusCode: http.StatusBadRequest,
				response:   "Wrong metric type\n",
			},
		},
	}

	for _, tt := range tests {
		body, code := testRequest(t, ts, tt.method, tt.path, tt.contentType)
		assert.Equal(t, tt.want.statusCode, code)
		assert.Equal(t, tt.want.response, body)
	}
}

func TestGetMetricsHandler(t *testing.T) {
	ms := storage.MemStorage{
		CounterMetrics: map[string]int64{"cTest1": 200, "cTest2": 128},
		GaugeMetrics:   map[string]float64{"gTest1": 2.12, "gTest2": 0.54},
	}
	rh := NewRequestHandler(ms)

	r := chi.NewRouter()
	r.Use(middleware.ContentCharset("UTF-8"))
	r.Use(middleware.AllowContentType("text/plain"))
	r.Get("/", rh.GetMetrics())

	ts := httptest.NewServer(r)
	defer ts.Close()

	type want struct {
		statusCode int
		response   string
	}
	tests := []struct {
		name        string
		path        string
		method      string
		contentType string
		want        want
	}{
		{
			name:        "fail method type",
			path:        "/",
			method:      http.MethodPost,
			contentType: "text/plain;charset=utf-8",
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				response:   "",
			},
		},
	}

	for _, tt := range tests {
		body, code := testRequest(t, ts, tt.method, tt.path, tt.contentType)
		assert.Equal(t, tt.want.statusCode, code)
		assert.Equal(t, tt.want.response, body)
	}
}

func TestPingHandler(t *testing.T) {
	ms := storage.MemStorage{
		CounterMetrics: map[string]int64{},
		GaugeMetrics:   map[string]float64{},
	}
	rh := NewRequestHandler(ms)

	r := chi.NewRouter()
	r.Use(middleware.ContentCharset("UTF-8"))
	r.Use(middleware.AllowContentType("text/plain"))
	r.Get("/ping", rh.Ping())

	ts := httptest.NewServer(r)
	defer ts.Close()

	type want struct {
		statusCode int
	}
	tests := []struct {
		want want
	}{
		{
			want: want{
				statusCode: http.StatusOK,
			},
		},
	}

	for _, tt := range tests {
		_, code := testRequest(t, ts, http.MethodGet, "/ping", "text/plain;charset=utf-8")
		assert.Equal(t, tt.want.statusCode, code)
	}
}

func TestUpdateMetric_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		UpdateGaugeMetric(gomock.Any(), gomock.Any()).
		Return(errors.New("repo error")).
		Times(1)

	mockRepo.EXPECT().
		UpdateCounterMetric(gomock.Any(), gomock.Any()).
		Return(errors.New("repo error")).
		Times(1)

	rh := NewRequestHandler(mockRepo)

	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", rh.UpdateMetric())

	ts := httptest.NewServer(r)
	defer ts.Close()

	tests := []struct {
		name        string
		path        string
		method      string
		contentType string
		want        int
	}{
		{
			name:        "repo error for gauge",
			path:        "/update/gauge/test/25",
			method:      http.MethodPost,
			contentType: "text/plain",
			want:        http.StatusBadRequest,
		},
		{
			name:        "repo error for counter",
			path:        "/update/counter/test/25",
			method:      http.MethodPost,
			contentType: "text/plain",
			want:        http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, code := testRequest(t, ts, tt.method, tt.path, tt.contentType)
			assert.Equal(t, tt.want, code)
		})
	}
}

func TestGetMetric_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		GetCounterMetric(gomock.Any()).
		Return(int64(0), errors.New("repo error")).
		Times(1)

	mockRepo.EXPECT().
		GetGaugeMetric(gomock.Any()).
		Return(float64(0), errors.New("repo error")).
		Times(1)

	rh := NewRequestHandler(mockRepo)

	r := chi.NewRouter()
	r.Get("/value/{type}/{name}", rh.GetMetric())

	ts := httptest.NewServer(r)
	defer ts.Close()

	tests := []struct {
		name        string
		path        string
		method      string
		contentType string
		want        int
	}{
		{
			name:        "repo error for gauge",
			path:        "/value/gauge/test",
			method:      http.MethodGet,
			contentType: "text/plain;charset=utf-8",
			want:        http.StatusNotFound,
		},
		{
			name:        "repo error for counter",
			path:        "/value/counter/test",
			method:      http.MethodGet,
			contentType: "text/plain;charset=utf-8",
			want:        http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, code := testRequest(t, ts, tt.method, tt.path, tt.contentType)
			assert.Equal(t, tt.want, code)
		})
	}
}

func TestPing_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		Ping().
		Return(errors.New("repo error")).
		Times(1)

	rh := NewRequestHandler(mockRepo)

	r := chi.NewRouter()
	r.Use(middleware.ContentCharset("UTF-8"))
	r.Use(middleware.AllowContentType("text/plain"))
	r.Get("/ping", rh.Ping())

	ts := httptest.NewServer(r)
	defer ts.Close()

	tests := []struct {
		name        string
		path        string
		method      string
		contentType string
		want        int
	}{
		{
			name:        "repo error",
			path:        "/ping",
			method:      http.MethodGet,
			contentType: "text/plain;charset=utf-8",
			want:        http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, code := testRequest(t, ts, tt.method, tt.path, tt.contentType)
			assert.Equal(t, tt.want, code)
		})
	}
}

func TestGetMetrics_CounterMetricsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		GetCounterMetrics().
		Return(nil, errors.New("repo error")).
		Times(1)

	rh := NewRequestHandler(mockRepo)

	r := chi.NewRouter()
	r.Use(middleware.ContentCharset("UTF-8"))
	r.Use(middleware.AllowContentType("text/plain"))
	r.Get("/", rh.GetMetrics())

	ts := httptest.NewServer(r)
	defer ts.Close()

	_, code := testRequest(t, ts, http.MethodGet, "/", "text/plain;charset=utf-8")
	assert.Equal(t, http.StatusInternalServerError, code)
}

func TestGetMetrics_GaugeMetricsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)

	mockRepo.EXPECT().
		GetCounterMetrics().
		Return(map[string]int64{"test": 123}, nil).
		Times(1)

	mockRepo.EXPECT().
		GetGaugeMetrics().
		Return(nil, errors.New("repo error")).
		Times(1)

	rh := NewRequestHandler(mockRepo)

	r := chi.NewRouter()
	r.Use(middleware.ContentCharset("UTF-8"))
	r.Use(middleware.AllowContentType("text/plain"))
	r.Get("/", rh.GetMetrics())

	ts := httptest.NewServer(r)
	defer ts.Close()

	_, code := testRequest(t, ts, http.MethodGet, "/", "text/plain;charset=utf-8")
	assert.Equal(t, http.StatusInternalServerError, code)
}

// ExampleRequestHandler_UpdateMetric demonstrates how to use the UpdateMetric handler.
func ExampleRequestHandler_UpdateMetric() {
	ms := storage.MemStorage{
		CounterMetrics: make(map[string]int64),
		GaugeMetrics:   make(map[string]float64),
	}

	// Create a new RequestHandler
	rh := NewRequestHandler(ms)

	// Create a new HTTP request with URL parameters
	req := httptest.NewRequest(http.MethodPost, "/update/gauge/cpu_usage/3.14", nil)

	// Create a response recorder to capture the response
	res := httptest.NewRecorder()

	// Create a Chi router and add the handler
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", rh.UpdateMetric())

	// Serve the request
	r.ServeHTTP(res, req)

	fmt.Println(res.Body.String())

	// Output:
	// Metric: cpu_usage value: 3.14 has added
}

// ExampleRequestHandler_GetMetric demonstrates how to use the GetMetric handler.
func ExampleRequestHandler_GetMetric() {
	ms := storage.MemStorage{
		CounterMetrics: make(map[string]int64),
		GaugeMetrics:   map[string]float64{"cpu_usage": 3.14, "memory_usage": 2.71},
	}

	// Create a new RequestHandler
	rh := NewRequestHandler(ms)

	// Create a new HTTP request with URL parameters
	req := httptest.NewRequest(http.MethodGet, "/value/gauge/cpu_usage", nil)

	// Create a response recorder to capture the response
	res := httptest.NewRecorder()

	// Create a Chi router and add the handler
	r := chi.NewRouter()
	r.Get("/value/{type}/{name}", rh.GetMetric())

	// Serve the request
	r.ServeHTTP(res, req)

	fmt.Println(res.Body.String())

	// Output:
	// 3.14
}

// ExampleRequestHandler_GetMetrics demonstrates how to use the GetMetrics handler.
func ExampleRequestHandler_GetMetrics() {
	ms := storage.MemStorage{
		CounterMetrics: make(map[string]int64),
		GaugeMetrics:   map[string]float64{"cpu_usage": 3.14, "memory_usage": 2.71},
	}

	// Create a new RequestHandler
	rh := NewRequestHandler(ms)

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Create a response recorder to capture the response
	res := httptest.NewRecorder()

	// Call the GetMetrics handler
	rh.GetMetrics()(res, req)

	fmt.Println(res.Body.String())

	// Output:
	// cpu_usage 3.14
	// memory_usage 2.71
}
