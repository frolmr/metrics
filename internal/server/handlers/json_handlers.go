package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/frolmr/metrics/internal/domain"
)

// UpdateMetricJSON updates a metric based on the provided JSON payload.
// @Summary Update a metric
// @Description Updates a metric with the provided JSON payload.
// @Tags metrics
// @Accept json
// @Produce json
// @Param metrics body domain.Metrics true "Metric data to update"
// @Success 200 {object} domain.Metrics "Updated metric"
// @Failure 400 {string} string "Invalid request payload or metric type"
// @Failure 500 {string} string "Internal server error"
// @Router /update [post]
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
			if updateErr := rh.repo.UpdateCounterMetric(metricsRequest.ID, *metricsRequest.Delta); updateErr != nil {
				http.Error(res, "error updating metric", http.StatusBadRequest)
				return
			}
		} else {
			if updateErr := rh.repo.UpdateGaugeMetric(metricsRequest.ID, *metricsRequest.Value); updateErr != nil {
				http.Error(res, "error updating metric", http.StatusBadRequest)
				return
			}
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

		if _, err := res.Write(resp); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// BulkUpdateMetricJSON updates multiple metrics based on the provided JSON payload.
// @Summary Update multiple metrics
// @Description Updates multiple metrics with the provided JSON payload.
// @Tags metrics
// @Accept json
// @Produce json
// @Param metrics body []domain.Metrics true "List of metrics to update"
// @Success 200 {string} string "Metrics updated successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Internal server error"
// @Router /updates [post]
func (rh *RequestHandler) BulkUpdateMetricJSON() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", domain.JSONContentType)

		metricsSlice, err := rh.readPayloadToMetricsSlice(req)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err := rh.repo.UpdateMetrics(metricsSlice); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		res.WriteHeader(http.StatusOK)
	}
}

// GetMetricJSON retrieves a metric based on the provided JSON payload.
// @Summary Get a metric
// @Description Retrieves a metric with the provided JSON payload.
// @Tags metrics
// @Accept json
// @Produce json
// @Param metrics body domain.Metrics true "Metric data to retrieve"
// @Success 200 {object} domain.Metrics "Requested metric"
// @Failure 400 {string} string "Invalid request payload or metric type"
// @Failure 404 {string} string "Metric not found"
// @Failure 500 {string} string "Internal server error"
// @Router /value [post]
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

		if _, err := res.Write(resp); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
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

func (rh *RequestHandler) readPayloadToMetricsSlice(req *http.Request) ([]domain.Metrics, error) {
	var metrics []domain.Metrics
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		return []domain.Metrics{}, err
	}
	if err := json.Unmarshal(buf.Bytes(), &metrics); err != nil {
		return []domain.Metrics{}, err
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
