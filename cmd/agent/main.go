package main

import (
	"log"
	"os"
	"time"

	"github.com/frolmr/metrics.git/internal/agent/config"
	"github.com/frolmr/metrics.git/internal/agent/metrics"
	"github.com/go-resty/resty/v2"
)

func main() {
	if err := config.GetConfig(); err != nil {
		log.Panic(err)
		os.Exit(1) // NOTE: не знаю на сколько это правильное/удачное решение
	}

	metrics := metrics.NewMetricsCollection()
	client := resty.New()

	go func() {
		for {
			metrics.CollectMetrics()
			time.Sleep(config.PollInterval)
		}
	}()

	for {
		metrics.ReportGaugeMetrics(client)
		metrics.ReportCounterMetrics(client)
		time.Sleep(config.ReportInterval)
	}
}
