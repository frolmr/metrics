package handlers

import (
	"net/http"

	"github.com/frolmr/metrics.git/internal/common/constants"
	"github.com/frolmr/metrics.git/internal/common/utils"
	"github.com/frolmr/metrics.git/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

func UpdateMetricHandler(repo storage.Repository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("content-type", constants.ContentType)

		// Пришлось вернуть, т.к. в автотестах проверка обязательная, а middlware
		// требует обязательно дополнять ContentCharset, который в автотестах не передается
		// if req.Header.Get("content-type") != "text/plain" {
		// 	res.WriteHeader(http.StatusUnsupportedMediaType)
		// 	return
		// }

		metricType := chi.URLParam(req, "type")

		if metricType != constants.GaugeType && metricType != constants.CounterType {
			http.Error(res, "Wrong metric type", http.StatusBadRequest)
			return
		}

		metricName := chi.URLParam(req, "name")
		metricValue := chi.URLParam(req, "value")

		if metricType == constants.GaugeType {
			value, err := utils.StringToFloat(metricValue)
			if err != nil {
				http.Error(res, "Wrong metric value", http.StatusBadRequest)
				return
			}
			repo.UpdateGaugeMetric(metricName, value)
		}

		if metricType == constants.CounterType {
			value, err := utils.StringToInt(metricValue)
			if err != nil {
				http.Error(res, "Wrong metric value", http.StatusBadRequest)
				return
			}
			repo.UpdateCounterMetric(metricName, value)
		}

		res.WriteHeader(http.StatusOK)
		res.Write([]byte("Metric: " + metricName + " value: " + metricValue + " has added"))
	}
}

func GetMetricHandler(repo storage.Repository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("content-type", constants.ContentType)

		metricType := chi.URLParam(req, "type")
		metricName := chi.URLParam(req, "name")

		switch metricType {
		case constants.CounterType:
			if value, err := repo.GetCounterMetric(metricName); err != nil {
				http.Error(res, "Metric Not Found", http.StatusNotFound)
			} else {
				res.WriteHeader(http.StatusOK)
				res.Write([]byte(utils.IntToString(value)))
			}
		case constants.GaugeType:
			if value, err := repo.GetGaugeMetric(metricName); err != nil {
				http.Error(res, "Metric Not Found", http.StatusNotFound)
			} else {
				res.WriteHeader(http.StatusOK)
				res.Write([]byte(utils.FloatToString(value)))
			}
		default:
			http.Error(res, "Wrong metric type", http.StatusBadRequest)
		}
	}
}

func GetMetricsHandler(repo storage.Repository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("content-type", constants.ContentType)

		for name, value := range repo.GetCounterMetrics() {
			res.Write([]byte(name + " " + utils.IntToString(value) + "\n"))
		}

		for name, value := range repo.GetGaugeMetrics() {
			res.Write([]byte(name + " " + utils.FloatToString(value) + "\n"))
		}
	}
}
