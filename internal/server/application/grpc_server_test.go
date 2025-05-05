package application

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/frolmr/metrics/internal/domain"
	"github.com/frolmr/metrics/internal/server/storage"
	pb "github.com/frolmr/metrics/pkg/proto/metrics"
	"github.com/frolmr/metrics/pkg/signer"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestMetricsServer(t *testing.T) {
	mockStorage := &storage.MemStorage{
		GaugeMetrics:   make(map[string]float64),
		CounterMetrics: make(map[string]int64),
	}

	server := NewMetricsServer(mockStorage, "test-key")

	t.Run("valid request", func(t *testing.T) {
		req := &pb.UpdateMetricsBulkRequest{
			Metrics: []*pb.Metric{
				{
					Key:    "test",
					Type:   pb.Metric_MTYPE_GAUGE,
					MValue: &pb.Metric_Value{Value: 1.23},
				},
			},
		}

		jsonData, err := json.Marshal(req)
		require.NoError(t, err)

		signature := signer.SignPayloadWithKey(jsonData, []byte("test-key"))

		md := metadata.New(map[string]string{
			domain.SignatureHeader: hex.EncodeToString(signature),
		})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		resp, err := server.UpdateMetricsBulk(ctx, req)
		require.NoError(t, err)
		require.True(t, resp.Received)

		val, exists := mockStorage.GaugeMetrics["test"]
		require.True(t, exists)
		require.Equal(t, 1.23, val)
	})

	t.Run("invalid signature", func(t *testing.T) {
		req := &pb.UpdateMetricsBulkRequest{
			Metrics: []*pb.Metric{
				{
					Key:    "test",
					Type:   pb.Metric_MTYPE_GAUGE,
					MValue: &pb.Metric_Value{Value: 1.23},
				},
			},
		}

		md := metadata.New(map[string]string{
			domain.SignatureHeader: "invalid_signature",
		})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		_, err := server.UpdateMetricsBulk(ctx, req)
		require.Error(t, err)
		require.Equal(t, codes.Unauthenticated, status.Code(err))
	})
}
