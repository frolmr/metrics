package main

import (
	"log"
	"net/http"
	"time"

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

	lgr, err := logger.NewLogger()
	if err != nil {
		log.Panic(err)
	}

	ms := storage.NewMemStorage()
	router := routes.NewRouter(ms, *lgr)

	server := &http.Server{
		Addr:              config.ServerAddress,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           router.SetupRoutes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
