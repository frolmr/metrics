package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricsUpdateHandler(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name        string
		metricType  string
		metricName  string
		metricValue string
		url         string
		method      string
		contentType string
		want        want
	}{
		{
			name:        "success gauge request",
			metricType:  "gauge",
			metricName:  "test",
			metricValue: "25",
			url:         "/update/gauge/test/25",
			method:      http.MethodPost,
			contentType: "text/plain",
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name:        "success counter request",
			metricType:  "counter",
			metricName:  "test",
			metricValue: "25",
			url:         "/update/counter/test/25",
			method:      http.MethodPost,
			contentType: "text/plain",
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name:        "fail content-type request",
			metricType:  "counter",
			metricName:  "test",
			metricValue: "25",
			url:         "/update/counter/test/25",
			method:      http.MethodPost,
			contentType: "application/json",
			want: want{
				statusCode: http.StatusUnsupportedMediaType,
			},
		},
		{
			name:        "fail wrong method request",
			metricType:  "counter",
			metricName:  "test",
			metricValue: "25",
			url:         "/update/counter/test/25",
			method:      http.MethodGet,
			contentType: "text/plain",
			want: want{
				statusCode: http.StatusMethodNotAllowed,
			},
		},
		{
			name:        "fail wrong metric type request",
			metricType:  "cntr",
			metricName:  "test",
			metricValue: "25",
			url:         "/update/cntr/test/25",
			method:      http.MethodPost,
			contentType: "text/plain",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:        "fail wrong gauge metric value format request",
			metricType:  "gauge",
			metricName:  "test",
			metricValue: "2,5",
			url:         "/update/gauge/test/25",
			method:      http.MethodPost,
			contentType: "text/plain",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:        "fail wrong counter metric value format request",
			metricType:  "counter",
			metricName:  "test",
			metricValue: "2.5",
			url:         "/update/counter/test/25",
			method:      http.MethodPost,
			contentType: "text/plain",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.url, nil)

			request.Header.Set("content-type", tt.contentType)
			request.SetPathValue("type", tt.metricType)
			request.SetPathValue("name", tt.metricName)
			request.SetPathValue("value", tt.metricValue)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(MetricsUpdateHandler)
			h(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			defer result.Body.Close()
		})
	}
}
