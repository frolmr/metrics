package metrics

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/frolmr/metrics.git/internal/common/constants"
)

const (
	serverURL = "http://localhost:8080"
)

func reportMetric(metricType, metricName, metricValue string) {
	url := serverURL + "/update/" + metricType + "/" + metricName + "/" + metricValue
	response, err := http.Post(url, constants.ContentType, nil)

	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer response.Body.Close()
}

func ReportCounterMetrics(counterMetrics *map[string]int64) {
	for key, value := range *counterMetrics {
		reportMetric(constants.CounterType, key, strconv.FormatInt(value, 10))
	}
}

func ReportGaugeMetrics(gaugeMetrics *map[string]float64) {
	for key, value := range *gaugeMetrics {
		reportMetric(constants.GaugeType, key, strconv.FormatFloat(value, 'g', -1, 64))
	}
}
