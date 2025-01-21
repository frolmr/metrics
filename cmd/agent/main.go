package main

import (
	"log"
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

	go func() {
		for {
			mtrcs.CollectMetrics()
			time.Sleep(cfg.PollInterval)
		}
	}()

	for {
		mtrcs.ReportMetrics()
		time.Sleep(cfg.ReportInterval)
	}
}
