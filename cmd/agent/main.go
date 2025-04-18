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
	"github.com/frolmr/metrics/pkg/buildinfo"
	"github.com/go-resty/resty/v2"
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

	client := resty.New()
	mtrcs := metrics.NewMetricsCollection(client, cfg)

	jobsCh := make(chan metrics.MetricsCollection, runtime.GOMAXPROCS(0))

	var wg sync.WaitGroup
	for i := 0; i < cfg.RateLimit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobsCh {
				job.ReportMetrics()
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
