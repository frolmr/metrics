package main

import (
	"log"
	"net/http"
	"time"

	"github.com/frolmr/metrics.git/internal/server/config"
	"github.com/frolmr/metrics.git/internal/server/controller"
	"github.com/frolmr/metrics.git/internal/server/logger"
	"github.com/frolmr/metrics.git/internal/server/storage"
)

func main() {
	var err error

	if err = config.GetConfig(); err != nil {
		log.Panic(err)
	}

	l, err := logger.NewLogger()
	if err != nil {
		log.Panic("error initializing logger")
	}

	ms := storage.NewMemStorage()
	fs := storage.NewFileSnapshot(ms, config.FileStoragePath)

	if config.Restore {
		if err := fs.RestoreData(); err != nil {
			log.Println("error restoring from snapshot: ", err.Error())
		}
	}

	if config.StoreInterval != 0 {
		makeSnapshots(fs)
	}

	c := controller.NewController(ms, l)

	server := &http.Server{
		Addr:              config.ServerAddress,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           c.SetupHandlers(),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func makeSnapshots(fs *storage.FileSnapshot) {
	f := func() {
		if err := fs.SaveData(); err != nil {
			log.Println("error saving data to snapshot: ", err.Error())
		}
		makeSnapshots(fs)
	}
	time.AfterFunc(config.StoreInterval, f)
}
