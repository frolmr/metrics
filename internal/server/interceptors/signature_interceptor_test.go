package interceptors

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/frolmr/metrics/internal/domain"
	pb "github.com/frolmr/metrics/pkg/proto/metrics"
	"github.com/frolmr/metrics/pkg/signer"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestSignatureInterceptor(t *testing.T) {
	signKey := "test-key"
	interceptor := NewSignatureInterceptor(signKey)

	mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return &pb.Ack{Received: true}, nil
	}

	t.Run("valid signature", func(t *testing.T) {
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

		signature := signer.SignPayloadWithKey(jsonData, []byte(signKey))

		md := metadata.New(map[string]string{
			domain.SignatureHeader: hex.EncodeToString(signature),
		})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		info := &grpc.UnaryServerInfo{
			FullMethod: "/metrics.Metrics/UpdateMetricsBulk",
		}

		resp, err := interceptor(ctx, req, info, mockHandler)
		require.NoError(t, err)
		require.True(t, resp.(*pb.Ack).Received)
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

		info := &grpc.UnaryServerInfo{
			FullMethod: "/metrics.Metrics/UpdateMetricsBulk",
		}

		_, err := interceptor(ctx, req, info, mockHandler)
		require.Error(t, err)
		require.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("no signature when required", func(t *testing.T) {
		req := &pb.UpdateMetricsBulkRequest{
			Metrics: []*pb.Metric{
				{
					Key:    "test",
					Type:   pb.Metric_MTYPE_GAUGE,
					MValue: &pb.Metric_Value{Value: 1.23},
				},
			},
		}

		ctx := context.Background()
		info := &grpc.UnaryServerInfo{
			FullMethod: "/metrics.Metrics/UpdateMetricsBulk",
		}

		_, err := interceptor(ctx, req, info, mockHandler)
		require.Error(t, err)
		require.Equal(t, codes.Unauthenticated, status.Code(err))
	})

	t.Run("skip validation for other methods", func(t *testing.T) {
		req := &pb.UpdateMetricsBulkRequest{
			Metrics: []*pb.Metric{
				{
					Key:    "test",
					Type:   pb.Metric_MTYPE_GAUGE,
					MValue: &pb.Metric_Value{Value: 1.23},
				},
			},
		}

		ctx := context.Background()
		info := &grpc.UnaryServerInfo{
			FullMethod: "/metrics.Metrics/OtherMethod",
		}

		resp, err := interceptor(ctx, req, info, mockHandler)
		require.NoError(t, err)
		require.True(t, resp.(*pb.Ack).Received)
	})

	t.Run("no validation when no key", func(t *testing.T) {
		interceptor := NewSignatureInterceptor("")

		req := &pb.UpdateMetricsBulkRequest{
			Metrics: []*pb.Metric{
				{
					Key:    "test",
					Type:   pb.Metric_MTYPE_GAUGE,
					MValue: &pb.Metric_Value{Value: 1.23},
				},
			},
		}

		ctx := context.Background()
		info := &grpc.UnaryServerInfo{
			FullMethod: "/metrics.Metrics/UpdateMetricsBulk",
		}

		resp, err := interceptor(ctx, req, info, mockHandler)
		require.NoError(t, err)
		require.True(t, resp.(*pb.Ack).Received)
	})
}
