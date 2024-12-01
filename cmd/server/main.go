package main

import (
	"net/http"

	"github.com/frolmr/metrics.git/internal/server/handlers"
	"github.com/frolmr/metrics.git/internal/server/storage"
	"github.com/go-chi/chi/v5"
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
	ms := storage.MemStorage{
		CounterMetrics: make(map[string]int64),
		GaugeMetrics:   make(map[string]float64),
	}

	r := chi.NewRouter()

	// r.Use(middleware.ContentCharset("UTF-8"))
	// r.Use(middleware.AllowContentType("text/plain"))

	r.Get("/", handlers.GetMetricsHandler(ms))

	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", handlers.UpdateMetricHandler(ms))
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", handlers.GetMetricHandler(ms))
	})

	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(err)
	}
}
