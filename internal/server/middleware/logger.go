package middleware

import (
	"net/http"
	"time"

	"github.com/frolmr/metrics.git/internal/server/logger"
	"go.uber.org/zap"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()

		responseData := &logger.ResponseData{
			Status: 0,
			Size:   0,
		}
		lw := logger.LoggingResponseWriter{
			ResponseWriter: res, // встраиваем оригинальный http.ResponseWriter
			ResponseData:   responseData,
		}
		next.ServeHTTP(&lw, req) // внедряем реализацию http.ResponseWriter

		duration := time.Since(start)

		logger.Log.Info("Request info: ", zap.String("uri", req.RequestURI), zap.String("method", req.Method), zap.Duration("duration", duration))
		logger.Log.Info("Response info: ", zap.Int("status", responseData.Status), zap.Int("size", responseData.Size))
	})
}
