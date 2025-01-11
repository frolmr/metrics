package controller

import (
	"database/sql"

	"github.com/frolmr/metrics.git/internal/server/handlers"
	"github.com/frolmr/metrics.git/internal/server/logger"
	"github.com/frolmr/metrics.git/internal/server/middleware"
	"github.com/frolmr/metrics.git/internal/server/storage"
	"github.com/go-chi/chi/v5"
)

type Controller struct {
	repo   storage.Repository
	logger *logger.Logger
	db     *sql.DB
}

func NewController(repo storage.Repository, lgr *logger.Logger, db *sql.DB) *Controller {
	return &Controller{
		repo:   repo,
		logger: lgr,
		db:     db,
	}
}

func (c *Controller) SetupHandlers() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Compressor)
	r.Use(middleware.WithLog(c.logger))

	rh := handlers.NewRequestHandler(c.repo, c.db)

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
