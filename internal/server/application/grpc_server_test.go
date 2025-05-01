package application

import (
	"context"
	"testing"

	"github.com/frolmr/metrics/internal/server/storage"
	pb "github.com/frolmr/metrics/pkg/proto/metrics"
	"github.com/stretchr/testify/require"
)

func TestMetricsServer(t *testing.T) {
	mockStorage := &storage.MemStorage{
		GaugeMetrics:   make(map[string]float64),
		CounterMetrics: make(map[string]int64),
	}

	server := NewMetricsServer(mockStorage)

	t.Run("valid gauge metric", func(t *testing.T) {
		req := &pb.UpdateMetricsBulkRequest{
			Metrics: []*pb.Metric{
				{
					Key:    "test",
					Type:   pb.Metric_MTYPE_GAUGE,
					MValue: &pb.Metric_Value{Value: 1.23},
				},
			},
		}

		resp, err := server.UpdateMetricsBulk(context.Background(), req)
		require.NoError(t, err)
		require.True(t, resp.Received)

		val, exists := mockStorage.GaugeMetrics["test"]
		require.True(t, exists)
		require.Equal(t, 1.23, val)
	})

	t.Run("valid counter metric", func(t *testing.T) {
		req := &pb.UpdateMetricsBulkRequest{
			Metrics: []*pb.Metric{
				{
					Key:    "count",
					Type:   pb.Metric_MTYPE_COUNTER,
					MValue: &pb.Metric_Delta{Delta: 42},
				},
			},
		}

		resp, err := server.UpdateMetricsBulk(context.Background(), req)
		require.NoError(t, err)
		require.True(t, resp.Received)

		val, exists := mockStorage.CounterMetrics["count"]
		require.True(t, exists)
		require.Equal(t, int64(42), val)
	})

	t.Run("multiple metrics", func(t *testing.T) {
		req := &pb.UpdateMetricsBulkRequest{
			Metrics: []*pb.Metric{
				{
					Key:    "gauge",
					Type:   pb.Metric_MTYPE_GAUGE,
					MValue: &pb.Metric_Value{Value: 1.23},
				},
				{
					Key:    "counter",
					Type:   pb.Metric_MTYPE_COUNTER,
					MValue: &pb.Metric_Delta{Delta: 42},
				},
			},
		}

		resp, err := server.UpdateMetricsBulk(context.Background(), req)
		require.NoError(t, err)
		require.True(t, resp.Received)

		gaugeVal, exists := mockStorage.GaugeMetrics["gauge"]
		require.True(t, exists)
		require.Equal(t, 1.23, gaugeVal)

		counterVal, exists := mockStorage.CounterMetrics["counter"]
		require.True(t, exists)
		require.Equal(t, int64(42), counterVal)
	})
}
