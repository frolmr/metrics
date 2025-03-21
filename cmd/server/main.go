// Server to receive and store metrics.

// @Title Metrics API
// @Description Service for metrics storage
// @Version 1.0

// @BasePath /
// @Host localhost:8080

// @Tag.name Health
// @Tag.description "Requests to check api health"

// @Tag.name Metrics
// @Tag.description "Requests to manipulate metrics"
package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "net/http/pprof" //nolint:gosec //need for the task

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

	if cfg.Profiling {
		go func() {
			log.Println("Starting pprof server on :6060...")

			server := &http.Server{
				Addr:         "localhost:6060",
				ReadTimeout:  3 * time.Second,
				WriteTimeout: 3 * time.Second,
				IdleTimeout:  5 * time.Second,
			}

			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Failed to start pprof server: %v", err)
			}
		}()
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
