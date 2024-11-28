package handlers

import (
	"net/http"
	"strconv"

	"github.com/frolmr/metrics.git/internal/common/constants"
)

func MetricsUpdateHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", constants.ContentType)

	if req.Method != http.MethodPost {
		http.Error(res, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

	if req.Header.Get("content-type") != constants.ContentType {
		res.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	metricType := req.PathValue("type")

	if metricType != constants.GaugeType && metricType != constants.CounterType {
		http.Error(res, "Wrong metric type", http.StatusBadRequest)
		return
	}

	metricValue := req.PathValue("value")

	if metricType == constants.GaugeType {
		_, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(res, "Wrong metric value", http.StatusBadRequest)
			return
		}
	}

	if metricType == constants.CounterType {
		_, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(res, "Wrong metric value", http.StatusBadRequest)
			return
		}
	}

	res.WriteHeader(http.StatusOK)
}
