package main

import (
	"log"
	"net/http"

	"github.com/frolmr/metrics.git/internal/server/config"
	"github.com/frolmr/metrics.git/internal/server/logger"
	"github.com/frolmr/metrics.git/internal/server/routes"
	"github.com/frolmr/metrics.git/internal/server/storage"
)

func main() {
	var err error

	if err = config.GetConfig(); err != nil {
		log.Panic(err)
	}

	logger, err := logger.NewLogger()
	if err != nil {
		log.Panic(err)
	}

	ms := storage.NewMemStorage()
	router := routes.NewRouter(ms, *logger)

	err = http.ListenAndServe(config.ServerAddress, router.SetupRoutes())
	if err != nil {
		log.Panic(err)
	}
}
