package metrics

import (
	"math/rand"
	"runtime"
)

func CollectMetrics(counterMetrics map[string]int64, gaugeMetrics map[string]float64) (map[string]int64, map[string]float64) {
	pollCount++

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	counterMetrics["PollCount"] = pollCount
	gaugeMetrics["RandomValue"] = rand.Float64()

	gaugeMetrics["Alloc"] = float64(ms.Alloc)
	gaugeMetrics["BuckHashSys"] = float64(ms.BuckHashSys)
	gaugeMetrics["Frees"] = float64(ms.Frees)
	gaugeMetrics["GCCPUFraction"] = float64(ms.GCCPUFraction)
	gaugeMetrics["GCSys"] = float64(ms.GCSys)
	gaugeMetrics["HeapAlloc"] = float64(ms.HeapAlloc)
	gaugeMetrics["HeapIdle"] = float64(ms.HeapIdle)
	gaugeMetrics["HeapInuse"] = float64(ms.HeapInuse)
	gaugeMetrics["HeapObjects"] = float64(ms.HeapObjects)
	gaugeMetrics["HeapReleased"] = float64(ms.HeapReleased)
	gaugeMetrics["HeapSys"] = float64(ms.HeapSys)
	gaugeMetrics["LastGC"] = float64(ms.LastGC)
	gaugeMetrics["Lookups"] = float64(ms.Lookups)
	gaugeMetrics["MCacheInuse"] = float64(ms.MCacheInuse)
	gaugeMetrics["MCacheSys"] = float64(ms.MCacheSys)
	gaugeMetrics["MSpanInuse"] = float64(ms.MSpanInuse)
	gaugeMetrics["MSpanSys"] = float64(ms.MSpanSys)
	gaugeMetrics["Mallocs"] = float64(ms.Mallocs)
	gaugeMetrics["NextGC"] = float64(ms.NextGC)
	gaugeMetrics["NumForcedGC"] = float64(ms.NumForcedGC)
	gaugeMetrics["NumGC"] = float64(ms.NumGC)
	gaugeMetrics["OtherSys"] = float64(ms.OtherSys)
	gaugeMetrics["PauseTotalNs"] = float64(ms.PauseTotalNs)
	gaugeMetrics["StackInuse"] = float64(ms.StackInuse)
	gaugeMetrics["StackSys"] = float64(ms.StackSys)
	gaugeMetrics["Sys"] = float64(ms.Sys)
	gaugeMetrics["TotalAlloc"] = float64(ms.TotalAlloc)

	return counterMetrics, gaugeMetrics
}
