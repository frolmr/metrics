// Package for routing and middleware.
package controller

import (
	"github.com/frolmr/metrics/internal/server/config"
	"github.com/frolmr/metrics/internal/server/handlers"
	"github.com/frolmr/metrics/internal/server/logger"
	"github.com/frolmr/metrics/internal/server/middleware"
	"github.com/frolmr/metrics/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

type Controller struct {
	logger *logger.Logger
	config *config.Config
}

// NewController function is constructor for controller object.
func NewController(lgr *logger.Logger, cfg *config.Config) *Controller {
	return &Controller{
		logger: lgr,
		config: cfg,
	}
}

// SetupHandlers functions is resonsible for app routing
func (c *Controller) SetupHandlers(stor storage.Repository) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Compressor)
	r.Use(middleware.WithLog(c.logger))
	r.Use(middleware.WithSignature(c.config.Key))

	rh := handlers.NewRequestHandler(stor)

	r.Get("/", rh.GetMetrics())

	r.Route("/update/", func(r chi.Router) {
		r.Post("/", rh.UpdateMetricJSON())
		r.Post("/{type}/{name}/{value}", rh.UpdateMetric())
	})

	r.Route("/value/", func(r chi.Router) {
		r.Post("/", rh.GetMetricJSON())
		r.Get("/{type}/{name}", rh.GetMetric())
	})

	r.Get("/ping", rh.Ping())
	r.Post("/updates/", rh.BulkUpdateMetricJSON())

	return r
}
