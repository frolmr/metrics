package application

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	_ "net/http/pprof" //nolint:gosec //need for the task

	pb "github.com/frolmr/metrics/pkg/proto/metrics"
	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/frolmr/metrics/internal/server/config"
	"github.com/frolmr/metrics/internal/server/controller"
	"github.com/frolmr/metrics/internal/server/db/migrator"
	"github.com/frolmr/metrics/internal/server/interceptors"
	"github.com/frolmr/metrics/internal/server/logger"
	"github.com/frolmr/metrics/internal/server/storage"
)

type Application struct {
	config         *config.Config
	logger         *logger.Logger
	httpServer     *http.Server
	pprofServer    *http.Server
	snapshotCancel context.CancelFunc
	wg             sync.WaitGroup
}

func NewApplication(cfg *config.Config, lgr *logger.Logger) *Application {
	return &Application{
		config: cfg,
		logger: lgr,
	}
}

func (app *Application) RunServer() error {
	storage, storageErr := app.setupStorage()
	if storageErr != nil {
		return fmt.Errorf("error while storage setup: %w", storageErr)
	}

	switch app.config.Scheme {
	case "http", "https":
		return app.runHTTPServer(storage)
	case "grpc":
		return app.runGRPCServer(storage)
	default:
		return errors.New("unknown protocol")
	}
}

func (app *Application) RunProfServer() {
	app.pprofServer = &http.Server{
		Addr:         "localhost:6060",
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	log.Println("Starting pprof server on :6060...")
	if err := app.pprofServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("pprof server error: %v", err)
	}
}

func (app *Application) runHTTPServer(stor storage.Repository) error {
	ctrl := controller.NewController(app.logger, app.config)

	app.httpServer = &http.Server{
		Addr:              app.config.HTTPAddress,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           ctrl.SetupHandlers(stor),
	}

	app.logger.SugaredLogger.Infof("Starting HTTP server on %s", app.config.HTTPAddress)
	return app.httpServer.ListenAndServe()
}

func (app *Application) runGRPCServer(stor storage.Repository) error {
	listen, err := net.Listen("tcp", app.config.HTTPAddress)
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption

	if app.config.Key != "" {
		opts = append(opts, grpc.UnaryInterceptor(interceptors.NewSignatureInterceptor(app.config.Key)))
	}

	if app.config.CryptoKey != nil {
		creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
		if err != nil {
			return fmt.Errorf("failed to create TLS credentials: %w", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)
	pb.RegisterMetricsServer(s, NewMetricsServer(stor))
	return s.Serve(listen)
}

func (app *Application) setupStorage() (storage.Repository, error) {
	if app.config.DatabaseDSN != "" {
		db, err := app.setupDB()
		if err != nil {
			return nil, fmt.Errorf("could not setup DB: %w", err)
		}
		retriableStor := storage.NewRetriableStorage(storage.NewDBStorage(db))
		return retriableStor, nil
	}

	memstor := storage.NewMemStorage()
	app.setupSnapshots(memstor)
	return memstor, nil
}

func (app *Application) setupDB() (*sql.DB, error) {
	db, err := sql.Open("pgx", app.config.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	m := migrator.NewMigrator(db)
	if migrationErr := m.RunMigrations(); migrationErr != nil {
		return nil, migrationErr
	}
	return db, err
}

func (app *Application) setupSnapshots(stor *storage.MemStorage) {
	fs := storage.NewFileSnapshot(stor, app.config.FileStoragePath)

	if app.config.Restore {
		if err := fs.RestoreData(); err != nil {
			app.logger.SugaredLogger.Error("error restoring from snapshot: ", err.Error())
		}
	}

	if app.config.StoreInterval != 0 {
		ctx, cancel := context.WithCancel(context.Background())
		app.snapshotCancel = cancel

		app.wg.Add(1)
		go app.runSnapshotSaver(ctx, fs)
	}
}

func (app *Application) runSnapshotSaver(ctx context.Context, fs *storage.FileSnapshot) {
	defer app.wg.Done()

	ticker := time.NewTicker(app.config.StoreInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := fs.SaveData(); err != nil {
				app.logger.SugaredLogger.Error("error saving data to snapshot: ", err.Error())
			}
		case <-ctx.Done():
			app.logger.SugaredLogger.Info("stopping snapshot saver...")
			if err := fs.SaveData(); err != nil {
				app.logger.SugaredLogger.Error("final snapshot save error: ", err.Error())
			}
			return
		}
	}
}

func (app *Application) Shutdown(ctx context.Context) error {
	if app.snapshotCancel != nil {
		app.snapshotCancel()
	}

	var httpErr error
	if app.httpServer != nil {
		httpErr = app.httpServer.Shutdown(ctx)
	}

	var pprofErr error
	if app.pprofServer != nil {
		pprofErr = app.pprofServer.Shutdown(ctx)
	}

	done := make(chan struct{})
	go func() {
		app.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
	}

	if httpErr != nil {
		return httpErr
	}
	return pprofErr
}
