package controller

import (
	"github.com/frolmr/metrics.git/internal/server/handlers"
	"github.com/frolmr/metrics.git/internal/server/logger"
	"github.com/frolmr/metrics.git/internal/server/middleware"
	"github.com/frolmr/metrics.git/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

type Controller struct {
	repo   storage.Repository
	logger *logger.Logger
}

func NewController(repo storage.Repository, lgr *logger.Logger) *Controller {
	return &Controller{
		repo:   repo,
		logger: lgr,
	}
}

func (c *Controller) SetupHandlers() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Compressor)
	r.Use(middleware.WithLog(c.logger))

	rh := handlers.NewRequestHandler(c.repo)

	r.Get("/", rh.GetMetrics())

	r.Route("/update", func(r chi.Router) {
		r.Post("/", rh.UpdateMetricJSON())
		r.Post("/{type}/{name}/{value}", rh.UpdateMetric())
	})

	r.Route("/value", func(r chi.Router) {
		r.Post("/", rh.GetMetricJSON())
		r.Get("/{type}/{name}", rh.GetMetric())
	})

	return r
}
