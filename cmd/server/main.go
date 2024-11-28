package main

import (
	"net/http"

	"github.com/frolmr/metrics.git/internal/server/handlers"
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

func main() {
	//ms := MemStorage{metrics: make(map[string]string)}

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/{type}/{name}/{value}`, handlers.MetricsUpdateHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
