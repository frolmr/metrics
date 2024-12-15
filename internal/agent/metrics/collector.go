package metrics

import (
	"crypto/rand"
	"encoding/binary"
	"log"
	"runtime"
)

type MetricsCollector interface {
	CollectMetrics() (map[string]int64, map[string]float64)
}

func (mc *MetricsCollection) CollectMetrics() (counterMetrics map[string]int64, gaugeMetrics map[string]float64) {
	pollCount++

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	mc.CounterMetrics["PollCount"] = pollCount
	mc.GaugeMetrics["RandomValue"] = randomFloat64()

	mc.GaugeMetrics["Alloc"] = float64(ms.Alloc)
	mc.GaugeMetrics["BuckHashSys"] = float64(ms.BuckHashSys)
	mc.GaugeMetrics["Frees"] = float64(ms.Frees)
	mc.GaugeMetrics["GCCPUFraction"] = float64(ms.GCCPUFraction)
	mc.GaugeMetrics["GCSys"] = float64(ms.GCSys)
	mc.GaugeMetrics["HeapAlloc"] = float64(ms.HeapAlloc)
	mc.GaugeMetrics["HeapIdle"] = float64(ms.HeapIdle)
	mc.GaugeMetrics["HeapInuse"] = float64(ms.HeapInuse)
	mc.GaugeMetrics["HeapObjects"] = float64(ms.HeapObjects)
	mc.GaugeMetrics["HeapReleased"] = float64(ms.HeapReleased)
	mc.GaugeMetrics["HeapSys"] = float64(ms.HeapSys)
	mc.GaugeMetrics["LastGC"] = float64(ms.LastGC)
	mc.GaugeMetrics["Lookups"] = float64(ms.Lookups)
	mc.GaugeMetrics["MCacheInuse"] = float64(ms.MCacheInuse)
	mc.GaugeMetrics["MCacheSys"] = float64(ms.MCacheSys)
	mc.GaugeMetrics["MSpanInuse"] = float64(ms.MSpanInuse)
	mc.GaugeMetrics["MSpanSys"] = float64(ms.MSpanSys)
	mc.GaugeMetrics["Mallocs"] = float64(ms.Mallocs)
	mc.GaugeMetrics["NextGC"] = float64(ms.NextGC)
	mc.GaugeMetrics["NumForcedGC"] = float64(ms.NumForcedGC)
	mc.GaugeMetrics["NumGC"] = float64(ms.NumGC)
	mc.GaugeMetrics["OtherSys"] = float64(ms.OtherSys)
	mc.GaugeMetrics["PauseTotalNs"] = float64(ms.PauseTotalNs)
	mc.GaugeMetrics["StackInuse"] = float64(ms.StackInuse)
	mc.GaugeMetrics["StackSys"] = float64(ms.StackSys)
	mc.GaugeMetrics["Sys"] = float64(ms.Sys)
	mc.GaugeMetrics["TotalAlloc"] = float64(ms.TotalAlloc)

	return mc.CounterMetrics, mc.GaugeMetrics
}

func randomFloat64() float64 {
	var buf [8]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		log.Fatal(err)
	}

	var f float64
	binary.LittleEndian.PutUint64(buf[:], uint64(f))

	return f
}
