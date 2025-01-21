package middleware

import (
	"bytes"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/frolmr/metrics.git/internal/domain"
	"github.com/frolmr/metrics.git/pkg/formatter"
)

type (
	responseBody struct {
		body []byte
	}

	signingResponseWriter struct {
		http.ResponseWriter
		responseBody *responseBody
		signKey      string
		wroteHeader  bool
	}
)

func (r *signingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	if r.signKey != "" {
		r.responseBody.body = append(r.responseBody.body, b...)
	}
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	return size, err
}

func (r *signingResponseWriter) WriteHeader(statusCode int) {
	if r.wroteHeader {
		return
	}

	defer r.ResponseWriter.WriteHeader(statusCode)

	r.wroteHeader = true

	if r.signKey != "" && statusCode == http.StatusOK {
		respSignature := formatter.SignPayloadWithKey(r.responseBody.body, []byte(r.signKey))
		r.ResponseWriter.Header().Add(domain.SignatureHeader, hex.EncodeToString(respSignature))
	}
}

func WithSignature(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(res http.ResponseWriter, req *http.Request) {
			signature := req.Header.Get(domain.SignatureHeader)
			if key != "" && signature != "" {
				bodyBytes, _ := io.ReadAll(req.Body)
				defer req.Body.Close()

				reqSignature := formatter.SignPayloadWithKey(bodyBytes, []byte(key))

				if hex.EncodeToString(reqSignature) == signature {
					responseBody := &responseBody{
						body: make([]byte, 0),
					}

					lw := signingResponseWriter{
						ResponseWriter: res,
						responseBody:   responseBody,
						signKey:        key,
					}

					req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

					next.ServeHTTP(&lw, req)
				} else {
					http.Error(res, "invalid signature", http.StatusBadGateway)
					return
				}
			} else {
				next.ServeHTTP(res, req)
			}
		}

		return http.HandlerFunc(fn)
	}
}
