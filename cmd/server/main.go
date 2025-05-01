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
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/frolmr/metrics/internal/server/application"
	"github.com/frolmr/metrics/internal/server/config"
	"github.com/frolmr/metrics/internal/server/logger"
	"github.com/frolmr/metrics/pkg/buildinfo"
)

const (
	shutdownTimeout = 10 * time.Second
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	buildinfo.PrintBuildInfo(buildVersion, buildDate, buildCommit)

	cfg, configErr := config.NewConfig()
	if configErr != nil {
		log.Panic(configErr)
	}

	lgr, loggerErr := logger.NewLogger()
	if loggerErr != nil {
		log.Panic(loggerErr)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	app := application.NewApplication(cfg, lgr)

	if cfg.Profiling {
		go func() {
			app.RunProfServer()
		}()
	}

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- app.RunServer()
	}()

	select {
	case err := <-serverErr:
		log.Panic(err)
	case <-ctx.Done():
		lgr.SugaredLogger.Info("shutting down server gracefully...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := app.Shutdown(shutdownCtx); err != nil {
			lgr.SugaredLogger.Error("server shutdown failed", "error", err)
		}
		lgr.SugaredLogger.Info("server exited properly")
	}
}
