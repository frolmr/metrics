package gzipper

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/frolmr/metrics.git/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCompressWriter(t *testing.T) {
	recorder := httptest.NewRecorder()

	cw := NewCompressWriter(recorder)

	data := []byte("test data")
	n, err := cw.Write(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), n)

	err = cw.Close()
	assert.NoError(t, err)

	assert.Equal(t, domain.CompressFormat, recorder.Header().Get("Content-Encoding"))

	gr, err := gzip.NewReader(recorder.Body)
	assert.NoError(t, err)
	defer gr.Close()

	decompressedData, err := io.ReadAll(gr)
	assert.NoError(t, err)
	assert.Equal(t, data, decompressedData)
}

func TestCompressWriter_WriteHeader(t *testing.T) {
	recorder := httptest.NewRecorder()

	cw := NewCompressWriter(recorder)

	cw.WriteHeader(http.StatusOK)

	assert.Equal(t, domain.CompressFormat, recorder.Header().Get("Content-Encoding"))
}

func TestCompressWriter_EmptyData(t *testing.T) {
	recorder := httptest.NewRecorder()

	cw := NewCompressWriter(recorder)

	n, err := cw.Write([]byte{})
	assert.NoError(t, err)
	assert.Equal(t, 0, n)

	err = cw.Close()
	assert.NoError(t, err)

	assert.Equal(t, domain.CompressFormat, recorder.Header().Get("Content-Encoding"))

	gr, err := gzip.NewReader(recorder.Body)
	assert.NoError(t, err)
	defer gr.Close()

	decompressedData, err := io.ReadAll(gr)
	assert.NoError(t, err)
	assert.Equal(t, []byte{}, decompressedData)
}

func TestCompressReader(t *testing.T) {
	data := []byte("test data")

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, err := gw.Write(data)
	assert.NoError(t, err)
	err = gw.Close()
	assert.NoError(t, err)

	cr, err := NewCompressReader(io.NopCloser(&buf))
	assert.NoError(t, err)

	decompressedData, err := io.ReadAll(cr)
	assert.NoError(t, err)
	assert.Equal(t, data, decompressedData)

	err = cr.Close()
	assert.NoError(t, err)
}

func TestCompressReader_Error(t *testing.T) {
	invalidData := []byte("invalid gzip data")

	_, err := NewCompressReader(io.NopCloser(bytes.NewReader(invalidData)))
	assert.Error(t, err)
}
