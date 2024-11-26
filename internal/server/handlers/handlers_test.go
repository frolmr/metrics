package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/frolmr/metrics.git/internal/server/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string, contentType string) (string, int) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	req.Header.Set("content-type", contentType)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.NoError(t, err)

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
			// TODO: тест иногда плавает
			name:        "success list metrics",
			path:        "/",
			method:      http.MethodGet,
			contentType: "text/plain;charset=utf-8",
			want: want{
				statusCode: http.StatusOK,
				response: "cTest1 200\n" +
					"cTest2 128\n" +
					"gTest1 2.12\n" +
					"gTest2 0.54\n",
			},
		},
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
