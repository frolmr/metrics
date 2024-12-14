package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/frolmr/metrics.git/internal/domain"
)

func (rh *RequestHandler) UpdateMetricJSON() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", domain.JSONContentType)

		metricsRequest, err := rh.readPayloadToMetrics(req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		if metricsRequest.MType != domain.GaugeType && metricsRequest.MType != domain.CounterType {
			http.Error(res, "wrong metric type", http.StatusBadRequest)
			return
		}

		if metricsRequest.Delta != nil {
			rh.repo.UpdateCounterMetric(metricsRequest.ID, *metricsRequest.Delta)
		} else {
			rh.repo.UpdateGaugeMetric(metricsRequest.ID, *metricsRequest.Value)
		}

		metricResponse, err := rh.prepareMetricsResponse(metricsRequest.ID, metricsRequest.MType)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(metricResponse)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
		_, _ = res.Write(resp)
	}
}

func (rh *RequestHandler) GetMetricJSON() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", domain.JSONContentType)

		metricsRequest, err := rh.readPayloadToMetrics(req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		if metricsRequest.MType != domain.GaugeType && metricsRequest.MType != domain.CounterType {
			http.Error(res, "wrong metric type", http.StatusBadRequest)
			return
		}

		metricResponse, err := rh.prepareMetricsResponse(metricsRequest.ID, metricsRequest.MType)
		if err != nil {
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		}

		resp, err := json.Marshal(metricResponse)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
		_, _ = res.Write(resp)
	}
}

func (rh *RequestHandler) readPayloadToMetrics(req *http.Request) (domain.Metrics, error) {
	var metrics domain.Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		return domain.Metrics{}, err
	}
	if err := json.Unmarshal(buf.Bytes(), &metrics); err != nil {
		return domain.Metrics{}, err
	}
	return metrics, nil
}

func (rh *RequestHandler) prepareMetricsResponse(metricName, metricType string) (domain.Metrics, error) {
	switch metricType {
	case domain.CounterType:
		metricValue, err := rh.repo.GetCounterMetric(metricName)
		if err != nil {
			return domain.Metrics{}, errors.New("metric not found")
		}
		return domain.Metrics{ID: metricName, MType: metricType, Delta: &metricValue, Value: nil}, nil
	case domain.GaugeType:
		metricValue, err := rh.repo.GetGaugeMetric(metricName)
		if err != nil {
			return domain.Metrics{}, errors.New("metric not found")
		}
		return domain.Metrics{ID: metricName, MType: metricType, Delta: nil, Value: &metricValue}, nil
	default:
		return domain.Metrics{}, errors.New("unknown metric type")
	}
}
