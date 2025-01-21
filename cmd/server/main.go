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

	ctrl := controller.NewController(l, cfg)

	var server *http.Server

	if cfg.DatabaseDSN != "" {
		db, err := setupDB(cfg)
		if err != nil {
			log.Panic("could not setup to DB: ", err.Error())
		}

		defer db.Close()

		retriableStor := storage.NewRetriableStorage(storage.NewDBStorage(db))
		server = setupServer(ctrl, retriableStor, cfg)
	} else {
		memstor := storage.NewMemStorage()
		server = setupServer(ctrl, memstor, cfg)
		setupSnapshots(cfg, memstor)
	}

	if err := server.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func setupDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	m := migrator.NewMigrator(db)
	if err := m.RunMigrations(); err != nil {
		return nil, err
	}
	return db, err
}

func setupServer(c *controller.Controller, stor storage.Repository, cfg *config.Config) *http.Server {
	return &http.Server{
		Addr:              cfg.HTTPAddress,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           c.SetupHandlers(stor),
	}
}

func setupSnapshots(cfg *config.Config, stor *storage.MemStorage) {
	fs := storage.NewFileSnapshot(stor, cfg.FileStoragePath)

	if cfg.Restore {
		if err := fs.RestoreData(); err != nil {
			log.Println("error restoring from snapshot: ", err.Error())
		}
	}

	if cfg.StoreInterval != 0 {
		makeSnapshots(fs, cfg.StoreInterval)
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
