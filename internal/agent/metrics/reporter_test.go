package metrics

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"syscall"
	"testing"

	"github.com/frolmr/metrics/internal/agent/config"
	"github.com/frolmr/metrics/internal/domain"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricsReporter(t *testing.T) {
	tests := []struct {
		name         string
		metrics      *MetricsCollection
		mockResponse *http.Response
		mockError    error
	}{
		{
			name: "successful report",
			metrics: &MetricsCollection{
				GaugeMetrics: map[string]float64{
					"test_gauge": 123.45,
				},
				CounterMetrics: map[string]int64{
					"test_counter": 67,
				},
				ReportClient: resty.New(),
				Config: &config.Config{
					Scheme:      "http",
					HTTPAddress: "localhost:8080",
				},
			},
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       httpmock.NewRespBodyFromString("OK"),
			},
		},
		{
			name: "server error - no retry",
			metrics: &MetricsCollection{
				GaugeMetrics: map[string]float64{
					"test_gauge": 123.45,
				},
				ReportClient: resty.New(),
				Config: &config.Config{
					Scheme:      "http",
					HTTPAddress: "localhost:8080",
				},
			},
			mockResponse: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       httpmock.NewRespBodyFromString("Server Error"),
			},
		},
		{
			name: "empty metrics - no request",
			metrics: &MetricsCollection{
				GaugeMetrics:   map[string]float64{},
				CounterMetrics: map[string]int64{},
				ReportClient:   resty.New(),
				Config: &config.Config{
					Scheme:      "http",
					HTTPAddress: "localhost:8080",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.ActivateNonDefault(tt.metrics.ReportClient.GetClient())
			defer httpmock.DeactivateAndReset()

			if tt.mockResponse != nil {
				httpmock.RegisterResponder(
					"POST",
					"http://localhost:8080/updates/",
					httpmock.ResponderFromResponse(tt.mockResponse),
				)
			} else if tt.mockError != nil {
				httpmock.RegisterResponder(
					"POST",
					"http://localhost:8080/updates/",
					httpmock.NewErrorResponder(tt.mockError),
				)
			}

			tt.metrics.ReportMetrics()

			info := httpmock.GetCallCountInfo()
			callCount := info["POST http://localhost:8080/updates/"]

			if len(tt.metrics.GaugeMetrics)+len(tt.metrics.CounterMetrics) == 0 {
				assert.Equal(t, 0, callCount, "expected no request for empty metrics")
				return
			}

			assert.Equal(t, 1, callCount, "expected exactly one request")
		})
	}
}

func TestReportMetricsWithSignature(t *testing.T) {
	cfg := &config.Config{
		Scheme:      "http",
		HTTPAddress: "localhost:8080",
		Key:         "test-key",
	}

	mc := &MetricsCollection{
		GaugeMetrics: map[string]float64{
			"test_gauge": 123.45,
		},
		ReportClient: resty.New(),
		Config:       cfg,
	}

	httpmock.ActivateNonDefault(mc.ReportClient.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"http://localhost:8080/updates/",
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "gzip", req.Header.Get("Content-Encoding"))

			signature := req.Header.Get(domain.SignatureHeader)
			assert.NotEmpty(t, signature, "signature header should be present")

			gz, err := gzip.NewReader(req.Body)
			require.NoError(t, err)
			defer gz.Close()

			var metrics []domain.Metrics
			err = json.NewDecoder(gz).Decode(&metrics)
			require.NoError(t, err)
			assert.Len(t, metrics, 1)
			assert.Equal(t, "test_gauge", metrics[0].ID)
			assert.Equal(t, 123.45, *metrics[0].Value)

			return httpmock.NewJsonResponse(http.StatusOK, map[string]string{"status": "OK"})
		},
	)

	mc.ReportMetrics()

	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info["POST http://localhost:8080/updates/"])
}

func TestCompressPayload(t *testing.T) {
	mc := &MetricsCollection{
		ReportClient: resty.New(),
		Config:       &config.Config{},
	}

	tests := []struct {
		name     string
		input    []byte
		validate func(t *testing.T, original []byte, compressed *bytes.Buffer)
	}{
		{
			name:  "small data",
			input: []byte("test"),
			validate: func(t *testing.T, original []byte, compressed *bytes.Buffer) {
				assert.NotEqual(t, original, compressed.Bytes(), "compressed data should be different")
			},
		},
		{
			name:  "medium data",
			input: []byte("this is a longer string that should compress well"),
			validate: func(t *testing.T, original []byte, compressed *bytes.Buffer) {
				assert.True(t, compressed.Len() < len(original), "compressed data should be smaller")
			},
		},
		{
			name:  "empty data",
			input: []byte{},
			validate: func(t *testing.T, original []byte, compressed *bytes.Buffer) {
				assert.NotNil(t, compressed, "should return buffer even for empty input")
				if compressed.Len() > 0 {
					gz, err := gzip.NewReader(compressed)
					assert.NoError(t, err, "should be valid gzip data")
					defer gz.Close()
					decompressed, err := io.ReadAll(gz)
					assert.NoError(t, err)
					assert.Equal(t, original, decompressed, "decompressed data should match original")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressed, err := mc.compressPayload(tt.input)
			require.NoError(t, err)
			require.NotNil(t, compressed)

			if len(tt.input) > 0 {
				gz, err := gzip.NewReader(compressed)
				require.NoError(t, err)
				defer gz.Close()

				decompressed, err := io.ReadAll(gz)
				require.NoError(t, err)
				assert.Equal(t, tt.input, decompressed, "decompressed data should match original")
			}

			tt.validate(t, tt.input, compressed)
		})
	}
}

func TestIsConnectionRefused(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "connection refused",
			err:      syscall.ECONNREFUSED,
			expected: true,
		},
		{
			name:     "other error",
			err:      errors.New("other error"),
			expected: false,
		},
		{
			name:     "wrapped connection refused",
			err:      &net.OpError{Err: &os.SyscallError{Err: syscall.ECONNREFUSED}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isConnectionRefused(tt.err))
		})
	}
}
