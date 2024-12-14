package main

import (
	"log"
	"time"

	"github.com/frolmr/metrics.git/internal/agent/config"
	"github.com/frolmr/metrics.git/internal/agent/metrics"
	"github.com/go-resty/resty/v2"
)

func main() {
	if err := config.GetConfig(); err != nil {
		log.Panic(err)
	}

	client := resty.New()
	mtrcs := metrics.NewMetricsCollection(client)

	go func() {
		for {
			mtrcs.CollectMetrics()
			time.Sleep(config.PollInterval)
		}
	}()

	for {
		mtrcs.ReportGaugeMetrics()
		mtrcs.ReportCounterMetrics()
		time.Sleep(config.ReportInterval)
	}
}
