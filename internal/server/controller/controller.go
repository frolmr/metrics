package controller

import (
	"github.com/frolmr/metrics.git/internal/server/handlers"
	"github.com/frolmr/metrics.git/internal/server/logger"
	"github.com/frolmr/metrics.git/internal/server/middleware"
	"github.com/frolmr/metrics.git/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

type Controller struct {
	logger *logger.Logger
}

func NewController(lgr *logger.Logger) *Controller {
	return &Controller{
		logger: lgr,
	}
}

func (c *Controller) SetupHandlers(storage storage.Repository) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Compressor)
	r.Use(middleware.WithLog(c.logger))

	rh := handlers.NewRequestHandler(storage)

	r.Get("/", rh.GetMetrics())

	r.Route("/update", func(r chi.Router) {
		r.Post("/", rh.UpdateMetricJSON())
		r.Post("/{type}/{name}/{value}", rh.UpdateMetric())
	})

	r.Route("/value", func(r chi.Router) {
		r.Post("/", rh.GetMetricJSON())
		r.Get("/{type}/{name}", rh.GetMetric())
	})

	r.Get("/ping", rh.Ping())

	return r
}
