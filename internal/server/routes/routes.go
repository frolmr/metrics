package routes

import (
	"github.com/frolmr/metrics.git/internal/server/handlers"
	"github.com/frolmr/metrics.git/internal/server/logger"
	"github.com/frolmr/metrics.git/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	repo   storage.Repository
	logger logger.Logger
}

func NewRouter(repo storage.Repository, logger logger.Logger) *Router {
	return &Router{
		repo:   repo,
		logger: logger,
	}
}

func (router *Router) SetupRoutes() chi.Router {
	r := chi.NewRouter()
	rh := handlers.NewRequestHandler(router.repo)

	r.Get("/", router.logger.WithLogging(rh.GetMetrics()))

	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", router.logger.WithLogging(rh.UpdateMetric()))
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", router.logger.WithLogging(rh.GetMetric()))
	})

	return r
}
