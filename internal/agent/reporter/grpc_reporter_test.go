package reporter

import (
	"context"
	"net"
	"testing"

	"github.com/frolmr/metrics/internal/agent/config"
	"github.com/frolmr/metrics/internal/agent/metrics"
	pb "github.com/frolmr/metrics/pkg/proto/metrics"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestGRPCReporter(t *testing.T) {
	t.Run("successful report", func(t *testing.T) {
		cfg := &config.Config{
			HTTPAddress: "localhost:8080",
			Key:         "test-key",
		}

		lis, err := net.Listen("tcp", cfg.HTTPAddress)
		require.NoError(t, err)
		defer lis.Close()

		s := grpc.NewServer()
		pb.RegisterMetricsServer(s, &mockMetricsServer{})

		serveErr := make(chan error, 1)
		go func() {
			serveErr <- s.Serve(lis)
		}()
		defer s.Stop()

		reporter, err := NewGRPCReporter(cfg)
		require.NoError(t, err)
		defer reporter.Close()

		ms := metrics.NewMetricsCollection()
		ms.GaugeMetrics["test"] = 1.23
		ms.CounterMetrics["count"] = 42

		reporter.ReportMetrics(*ms)

		select {
		case err := <-serveErr:
			require.NoError(t, err)
		default:
		}
	})

	t.Run("with retries", func(t *testing.T) {
		cfg := &config.Config{
			HTTPAddress: "invalid-address",
		}

		reporter, err := NewGRPCReporter(cfg)
		require.NoError(t, err)
		defer reporter.Close()

		ms := metrics.NewMetricsCollection()
		ms.GaugeMetrics["test"] = 1.23

		reporter.ReportMetrics(*ms)
	})
}

type mockMetricsServer struct {
	pb.UnimplementedMetricsServer
}

func (m *mockMetricsServer) UpdateMetricsBulk(ctx context.Context, req *pb.UpdateMetricsBulkRequest) (*pb.Ack, error) {
	return &pb.Ack{Received: true}, nil
}
