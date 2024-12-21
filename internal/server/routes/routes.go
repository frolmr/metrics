package routes

import (
	"github.com/frolmr/metrics.git/internal/server/handlers"
	"github.com/frolmr/metrics.git/internal/server/middleware"
	"github.com/frolmr/metrics.git/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	repo storage.Repository
}

func NewRouter(repo storage.Repository) *Router {
	return &Router{
		repo: repo,
	}
}

func (router *Router) SetupRoutes() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Compressor)
	r.Use(middleware.Logger)

	rh := handlers.NewRequestHandler(router.repo)

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
