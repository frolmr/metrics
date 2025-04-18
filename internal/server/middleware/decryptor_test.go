package middleware

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/frolmr/metrics/internal/server/decryptor"
	"github.com/stretchr/testify/require"
)

func TestWithDecrypt(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	t.Run("successful decryption", func(t *testing.T) {
		d := decryptor.NewDecryptor(privateKey)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			require.Equal(t, "test data", string(data))
			w.WriteHeader(http.StatusOK)
		})

		encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, &privateKey.PublicKey, []byte("test data"))
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/", bytes.NewReader(encrypted))
		req.Header.Set("Content-Type", "application/octet-stream")

		rr := httptest.NewRecorder()

		middleware := WithDecrypt(d)
		middleware(handler).ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("no private key - pass through", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			data, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			require.Equal(t, "test data", string(data))
			w.WriteHeader(http.StatusOK)
		})

		d := &decryptor.Decryptor{PrivateKey: nil}

		req := httptest.NewRequest("POST", "/", bytes.NewReader([]byte("test data")))
		req.Header.Set("Content-Type", "application/octet-stream")

		rr := httptest.NewRecorder()

		middleware := WithDecrypt(d)
		middleware(handler).ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		require.True(t, called, "handler should be called")
	})

	t.Run("nil decryptor - pass through", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			data, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			require.Equal(t, "test data", string(data))
			w.WriteHeader(http.StatusOK)
		})

		req := httptest.NewRequest("POST", "/", bytes.NewReader([]byte("test data")))
		req.Header.Set("Content-Type", "application/octet-stream")

		rr := httptest.NewRecorder()

		middleware := WithDecrypt(nil)
		middleware(handler).ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		require.True(t, called, "handler should be called")
	})

	t.Run("decryption failure", func(t *testing.T) {
		d := decryptor.NewDecryptor(privateKey)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called")
		})

		req := httptest.NewRequest("POST", "/", bytes.NewReader([]byte("invalid encrypted data")))
		req.Header.Set("Content-Type", "application/octet-stream")

		rr := httptest.NewRecorder()

		middleware := WithDecrypt(d)
		middleware(handler).ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("read body failure", func(t *testing.T) {
		d := decryptor.NewDecryptor(privateKey)

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("handler should not be called")
		})

		req := httptest.NewRequest("POST", "/", &failingReader{})
		req.Header.Set("Content-Type", "application/octet-stream")

		rr := httptest.NewRecorder()

		middleware := WithDecrypt(d)
		middleware(handler).ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

type failingReader struct{}

func (r *failingReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("mock read error")
}
