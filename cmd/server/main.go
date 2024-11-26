package main

import (
	"net/http"
	"strconv"
)

const (
	gaugeType   = "gauge"
	counterType = "counter"
)

type Metricable interface {
	AddMetric(name string, value string)
}

type MemStorage struct {
	metrics map[string]string
}

func (ms *MemStorage) AddMetric(name string, value string) {
	ms.metrics[name] = value
}

func metricsUpdate(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

	if req.Header.Get("content-type") != "text/plain" {
		res.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	metricType := req.PathValue("type")

	if metricType != gaugeType && metricType != counterType {
		http.Error(res, "Wrong metric type", http.StatusBadRequest)
		return
	}

	metricValue := req.PathValue("value")

	if metricType == gaugeType {
		_, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(res, "Wrong metric value", http.StatusBadRequest)
			return
		}
	}

	if metricType == counterType {
		_, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(res, "Wrong metric value", http.StatusBadRequest)
			return
		}
	}

	res.Header().Set("content-type", "text/plain")
	res.Header().Set("charset", "utf-8")
	res.WriteHeader(http.StatusOK)
}

func main() {
	//ms := MemStorage{metrics: make(map[string]string)}

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/{type}/{name}/{value}`, metricsUpdate)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
