package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/frolmr/metrics/internal/server/decryptor"
)

func WithDecrypt(d *decryptor.Decryptor) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(res http.ResponseWriter, req *http.Request) {
			if d == nil || d.PrivateKey == nil {
				next.ServeHTTP(res, req)
				return
			}

			encryptedData, err := io.ReadAll(req.Body)
			if err != nil {
				http.Error(res, "failed to read request body", http.StatusBadRequest)
				return
			}
			defer req.Body.Close()

			decryptedData, err := d.DecryptData(encryptedData)
			if err != nil {
				http.Error(res, "failed to decrypt data", http.StatusBadRequest)
				return
			}

			req.Body = io.NopCloser(bytes.NewReader(decryptedData))
			req.ContentLength = int64(len(decryptedData))

			next.ServeHTTP(res, req)
		}
		return http.HandlerFunc(fn)
	}
}
