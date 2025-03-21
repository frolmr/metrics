// Agent to collect and send metrics to server.
package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/frolmr/metrics.git/internal/agent/config"
	"github.com/frolmr/metrics.git/internal/agent/metrics"
	"github.com/go-resty/resty/v2"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Panic(err)
	}

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
	go func() {
		for {
			select {
			case <-pollTicker.C:
				go mtrcs.CollectMetrics()
				go mtrcs.CollectAdditionalMetrics()
			case <-doneCh:
				return
			}
		}
	}()

	reportTicker := time.NewTicker(cfg.ReportInterval)
	go func() {
		for {
			select {
			case <-reportTicker.C:
				jobsCh <- *mtrcs
			case <-doneCh:
				return
			}
		}
	}()

	termCh := make(chan os.Signal, 1)
	signal.Notify(termCh, syscall.SIGINT)
	<-termCh
	pollTicker.Stop()
	reportTicker.Stop()
	doneCh <- struct{}{}
	close(jobsCh)
	wg.Wait()
}
