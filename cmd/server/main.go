package main

import (
	"log"
	"net/http"
	"os"

	"github.com/frolmr/metrics.git/internal/server/config"
	"github.com/frolmr/metrics.git/internal/server/routes"
	"github.com/frolmr/metrics.git/internal/server/storage"
)

func main() {
	var err error

	if err = config.GetConfig(); err != nil {
		log.Panic(err)
		os.Exit(1) // NOTE: не знаю на сколько это правильное/удачное решениe
	}

	ms := storage.NewMemStorage()

	r := routes.SetupRoutes(ms)

	err = http.ListenAndServe(config.ServerAddress, r)
	if err != nil {
		log.Panic(err)
		os.Exit(1) // NOTE: не знаю на сколько это правильное/удачное решениe
	}
}
