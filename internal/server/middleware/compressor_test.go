package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/frolmr/metrics/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCompressor_CompressedResponse(t *testing.T) {
	handler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		_, _ = res.Write([]byte("test data"))
	})

	middleware := Compressor(handler)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Encoding", domain.CompressFormat)

	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	assert.Equal(t, domain.CompressFormat, rec.Header().Get("Content-Encoding"))

	gr, err := gzip.NewReader(rec.Body)
	assert.NoError(t, err)
	defer gr.Close()

	decompressedData, err := io.ReadAll(gr)
	assert.NoError(t, err)
	assert.Equal(t, "test data", string(decompressedData))
}

func TestCompressor_CompressedRequest(t *testing.T) {
	handler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		assert.NoError(t, err)
		_, _ = res.Write(body)
	})

	middleware := Compressor(handler)

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, err := gw.Write([]byte("test data"))
	assert.NoError(t, err)
	gw.Close()

	req := httptest.NewRequest("POST", "/", &buf)
	req.Header.Set("Content-Encoding", domain.CompressFormat)

	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	assert.Equal(t, "test data", rec.Body.String())
}

func TestCompressor_UncompressedResponse(t *testing.T) {
	handler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		_, _ = res.Write([]byte("test data"))
	})

	middleware := Compressor(handler)

	req := httptest.NewRequest("GET", "/", nil)

	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	assert.Empty(t, rec.Header().Get("Content-Encoding"))

	assert.Equal(t, "test data", rec.Body.String())
}

func TestCompressor_DecompressionError(t *testing.T) {
	handler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		_, _ = res.Write([]byte("test data"))
	})

	middleware := Compressor(handler)

	req := httptest.NewRequest("POST", "/", strings.NewReader("invalid gzip data"))
	req.Header.Set("Content-Encoding", domain.CompressFormat)

	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
