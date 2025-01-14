package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/frolmr/metrics.git/internal/server/config"
	"github.com/frolmr/metrics.git/internal/server/controller"
	"github.com/frolmr/metrics.git/internal/server/db/migrator"
	"github.com/frolmr/metrics.git/internal/server/logger"
	"github.com/frolmr/metrics.git/internal/server/storage"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Panic(err)
	}

	l, err := logger.NewLogger()
	if err != nil {
		log.Panic("error initializing logger")
	}

	ctrl := controller.NewController(l)

	var dbstor *storage.DBStorage
	var memstor *storage.MemStorage
	var server *http.Server

	if cfg.DatabaseDSN != "" {
		db, err := setupDB(cfg)
		if err != nil {
			log.Panic("could not setup to DB: ", err.Error())
		}

		defer db.Close()

		dbstor = storage.NewDBStorage(db)
		server = setupServer(ctrl, dbstor, cfg)
	} else {
		memstor = storage.NewMemStorage()
		server = setupServer(ctrl, memstor, cfg)
		setupSnapshots(cfg, memstor)
	}

	if err := server.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func setupDB(config *config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	migrator := migrator.NewMigrator(db)
	if err := migrator.RunMigrations(); err != nil {
		return nil, err
	}
	return db, err
}

func setupServer(c *controller.Controller, stor storage.Repository, config *config.Config) *http.Server {
	server := &http.Server{
		Addr:              config.HTTPAddress,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           c.SetupHandlers(stor),
	}
	return server
}

func setupSnapshots(config *config.Config, stor *storage.MemStorage) {
	fs := storage.NewFileSnapshot(stor, config.FileStoragePath)

	if config.Restore {
		if err := fs.RestoreData(); err != nil {
			log.Println("error restoring from snapshot: ", err.Error())
		}
	}

	if config.StoreInterval != 0 {
		makeSnapshots(fs, config.StoreInterval)
	}
}

func makeSnapshots(fs *storage.FileSnapshot, storeInterval time.Duration) {
	f := func() {
		if err := fs.SaveData(); err != nil {
			log.Println("error saving data to snapshot: ", err.Error())
		}
		makeSnapshots(fs, storeInterval)
	}
	time.AfterFunc(storeInterval, f)
}
