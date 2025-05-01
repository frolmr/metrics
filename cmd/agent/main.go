// Agent to collect and send metrics to server.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/frolmr/metrics/internal/agent/config"
	"github.com/frolmr/metrics/internal/agent/metrics"
	"github.com/frolmr/metrics/internal/agent/reporter"
	"github.com/frolmr/metrics/pkg/buildinfo"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	buildinfo.PrintBuildInfo(buildVersion, buildDate, buildCommit)

	cfg, err := config.NewConfig()
	if err != nil {
		log.Panic(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	mtrcs := metrics.NewMetricsCollection()

	var metricsReporter metrics.MetricsReporter
	switch cfg.Scheme {
	case "http", "https":
		metricsReporter = reporter.NewHTTPReporter(cfg)
	case "grpc":
		var err error
		metricsReporter, err = reporter.NewGRPCReporter(cfg)
		if err != nil {
			log.Panic(err)
		}
	default:
		log.Panic("invalid protocol")
	}
	defer metricsReporter.Close()

	jobsCh := make(chan metrics.MetricsCollection, runtime.GOMAXPROCS(0))

	var wg sync.WaitGroup
	for i := 0; i < cfg.RateLimit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobsCh {
				metricsReporter.ReportMetrics(job)
			}
		}()
	}

	doneCh := make(chan struct{})
	defer close(doneCh)

	pollTicker := time.NewTicker(cfg.PollInterval)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(cfg.ReportInterval)
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			go mtrcs.CollectMetrics()
			go mtrcs.CollectAdditionalMetrics()
		case <-reportTicker.C:
			jobsCh <- *mtrcs
		case <-ctx.Done():
			close(jobsCh)
			wg.Wait()
			log.Println("Agent shutdown gracefully")
			return
		}
	}
}
