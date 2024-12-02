package routes

import (
	"github.com/frolmr/metrics.git/internal/server/handlers"
	"github.com/frolmr/metrics.git/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(repo storage.Repository) chi.Router {
	r := chi.NewRouter()
	rh := handlers.NewRequestHandler(repo)

	r.Get("/", rh.GetMetrics())

	r.Route("/update", func(r chi.Router) {
		r.Post("/{type}/{name}/{value}", rh.UpdateMetric())
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/{type}/{name}", rh.GetMetric())
	})

	return r
}
