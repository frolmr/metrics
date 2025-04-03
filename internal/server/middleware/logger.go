package middleware

import (
	"net/http"
	"time"

	"github.com/frolmr/metrics/internal/server/logger"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
		wroteHeader  bool
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	if r.wroteHeader {
		return
	}
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
	r.wroteHeader = true
}

func WithLog(l *logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(res http.ResponseWriter, req *http.Request) {
			start := time.Now()

			responseData := &responseData{
				status: 0,
				size:   0,
			}
			lw := loggingResponseWriter{
				ResponseWriter: res,
				responseData:   responseData,
			}
			next.ServeHTTP(&lw, req)

			duration := time.Since(start)

			l.SugaredLogger.Infoln(
				"uri", req.RequestURI,
				"method", req.Method,
				"status", responseData.status,
				"duration", duration,
				"size", responseData.size,
			)
		}

		return http.HandlerFunc(fn)
	}
}
