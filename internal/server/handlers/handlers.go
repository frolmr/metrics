// Package handlers hold the business logic of application
package handlers

import (
	"net/http"

	"github.com/frolmr/metrics/internal/domain"
	"github.com/frolmr/metrics/internal/server/storage"
	"github.com/frolmr/metrics/pkg/formatter"
	"github.com/go-chi/chi/v5"
)

type RequestHandler struct {
	repo storage.Repository
}

// NewRequestHandler function is the constructor for handler object that has methods for hadnling requests for app
func NewRequestHandler(repo storage.Repository) *RequestHandler {
	return &RequestHandler{
		repo: repo,
	}
}

type MetricsRequester interface {
	Ping() http.HandlerFunc
	UpdateMetric() http.HandlerFunc
	UpdateMetricJSON() http.HandlerFunc
	GetMetric() http.HandlerFunc
	GetMetricJSON() http.HandlerFunc
	GetMetrics() http.HandlerFunc
}

// Ping godoc
// @Tags Health
// @Summary Request for API health check
// @Success 200
// @Failure 500
// @Router /ping [get]
func (rh *RequestHandler) Ping() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if err := rh.repo.Ping(); err != nil {
			http.Error(res, "DB unavailable", http.StatusInternalServerError)
			return
		}
	}
}

// UpdateMetric updates a metric based on the provided URL parameters.
// @Summary Update a metric
// @Description Updates a metric with the provided type, name, and value in the URL path.
// @Tags metrics
// @Accept plain
// @Produce plain
// @Param type path string true "Type of the metric (gauge or counter)"
// @Param name path string true "Name of the metric"
// @Param value path string true "Value of the metric"
// @Success 200 {string} string "Metric updated successfully"
// @Failure 400 {string} string "Invalid metric type or value"
// @Failure 500 {string} string "Internal server error"
// @Router /update/{type}/{name}/{value} [post]
func (rh *RequestHandler) UpdateMetric() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("content-type", domain.TextContentType)

		metricType := chi.URLParam(req, "type")

		if metricType != domain.GaugeType && metricType != domain.CounterType {
			http.Error(res, "Wrong metric type", http.StatusBadRequest)
			return
		}

		metricName := chi.URLParam(req, "name")
		metricValue := chi.URLParam(req, "value")

		if err := rh.updateMetric(metricName, metricType, metricValue); err != nil {
			http.Error(res, "Wrong metric value", http.StatusBadRequest)
			return
		}

		if _, err := res.Write([]byte("Metric: " + metricName + " value: " + metricValue + " has added")); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// GetMetric handles requests to fetch a metric by type and name.
// @Summary Get a metric by type and name
// @Description Fetches the value of a metric based on the provided type and name.
// @Tags Metrics
// @Accept plain
// @Produce plain
// @Param type path string true "Type of the metric"
// @Param name path string true "Name of the metric"
// @Success 200 {string} string "The value of the metric"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Metric not found"
// @Router /metrics/{type}/{name} [get]
func (rh *RequestHandler) GetMetric() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("content-type", domain.TextContentType)

		metricType := chi.URLParam(req, "type")
		metricName := chi.URLParam(req, "name")

		switch metricType {
		case domain.CounterType:
			if value, err := rh.repo.GetCounterMetric(metricName); err != nil {
				http.Error(res, "Metric Not Found", http.StatusNotFound)
			} else {
				if _, err := res.Write([]byte(formatter.IntToString(value))); err != nil {
					http.Error(res, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		case domain.GaugeType:
			if value, err := rh.repo.GetGaugeMetric(metricName); err != nil {
				http.Error(res, "Metric Not Found", http.StatusNotFound)
			} else {
				if _, err := res.Write([]byte(formatter.FloatToString(value))); err != nil {
					http.Error(res, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		default:
			http.Error(res, "Wrong metric type", http.StatusBadRequest)
		}
	}
}

// GetMetrics returns a list of all metrics (counter and gauge) in plain text.
// @Summary Get all metrics
// @Description Returns a list of all counter and gauge metrics in plain text format.
// @Tags metrics
// @Produce plain
// @Produce html
// @Success 200 {string} string "List of metrics in plain text or HTML format"
// @Failure 500 {string} string "Internal server error"
// @Router / [get]
func (rh *RequestHandler) GetMetrics() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// NOTE: почему-то автотесты 8-й итерации хотят тут "text/html"
		if req.Header.Get("Accept-Encoding") == domain.CompressFormat {
			res.Header().Set("content-type", domain.HTMLContentType)
		} else {
			res.Header().Set("content-type", domain.TextContentType)
		}

		counterMetrics, err := rh.repo.GetCounterMetrics()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		gaugeMetrics, err := rh.repo.GetGaugeMetrics()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		for name, value := range counterMetrics {
			if _, err := res.Write([]byte(name + " " + formatter.IntToString(value) + "\n")); err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		for name, value := range gaugeMetrics {
			if _, err := res.Write([]byte(name + " " + formatter.FloatToString(value) + "\n")); err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}

func (rh *RequestHandler) updateMetric(metricName, metricType, metricValue string) error {
	if metricType == domain.GaugeType {
		value, err := formatter.StringToFloat(metricValue)
		if err != nil {
			return err
		}
		if err := rh.repo.UpdateGaugeMetric(metricName, value); err != nil {
			return err
		}
	}

	if metricType == domain.CounterType {
		value, err := formatter.StringToInt(metricValue)
		if err != nil {
			return err
		}
		if err := rh.repo.UpdateCounterMetric(metricName, value); err != nil {
			return err
		}
	}
	return nil
}
