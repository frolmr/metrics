package handlers

import (
	"net/http"

	"github.com/frolmr/metrics.git/internal/domain"
	"github.com/frolmr/metrics.git/internal/server/storage"
	"github.com/frolmr/metrics.git/pkg/utils"
	"github.com/go-chi/chi/v5"
)

type RequestHandler struct {
	repo storage.Repository
}

func NewRequestHandler(repo storage.Repository) *RequestHandler {
	return &RequestHandler{
		repo: repo,
	}
}

type MetricsRequester interface {
	UpdateMetric() http.HandlerFunc
	UpdateMetricJSON() http.HandlerFunc
	GetMetric() http.HandlerFunc
	GetMetricJSON() http.HandlerFunc
	GetMetrics() http.HandlerFunc
}

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

		res.WriteHeader(http.StatusOK)
		_, _ = res.Write([]byte("Metric: " + metricName + " value: " + metricValue + " has added"))
	}
}

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
				res.WriteHeader(http.StatusOK)
				_, _ = res.Write([]byte(utils.IntToString(value)))
			}
		case domain.GaugeType:
			if value, err := rh.repo.GetGaugeMetric(metricName); err != nil {
				http.Error(res, "Metric Not Found", http.StatusNotFound)
			} else {
				res.WriteHeader(http.StatusOK)
				_, _ = res.Write([]byte(utils.FloatToString(value)))
			}
		default:
			http.Error(res, "Wrong metric type", http.StatusBadRequest)
		}
	}
}

func (rh *RequestHandler) GetMetrics() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("content-type", domain.TextContentType)

		for name, value := range rh.repo.GetCounterMetrics() {
			_, _ = res.Write([]byte(name + " " + utils.IntToString(value) + "\n"))
		}

		for name, value := range rh.repo.GetGaugeMetrics() {
			_, _ = res.Write([]byte(name + " " + utils.FloatToString(value) + "\n"))
		}
	}
}

func (rh *RequestHandler) updateMetric(metricName, metricType, metricValue string) error {
	if metricType == domain.GaugeType {
		value, err := utils.StringToFloat(metricValue)
		if err != nil {
			return err
		}
		rh.repo.UpdateGaugeMetric(metricName, value)
	}

	if metricType == domain.CounterType {
		value, err := utils.StringToInt(metricValue)
		if err != nil {
			return err
		}
		rh.repo.UpdateCounterMetric(metricName, value)
	}
	return nil
}
