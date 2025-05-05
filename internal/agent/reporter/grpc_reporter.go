package reporter

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/frolmr/metrics/internal/agent/config"
	"github.com/frolmr/metrics/internal/agent/metrics"
	"github.com/frolmr/metrics/internal/domain"
	pb "github.com/frolmr/metrics/pkg/proto/metrics"
	"github.com/frolmr/metrics/pkg/signer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	requestTimeout = 5 * time.Second
)

type GRPCReporter struct {
	config *config.Config
	client pb.MetricsClient
	conn   *grpc.ClientConn
}

func NewGRPCReporter(cfg *config.Config) (*GRPCReporter, error) {
	var opts []grpc.DialOption

	if cfg.CryptoKey == nil {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		creds := credentials.NewClientTLSFromCert(nil, "")
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	conn, err := grpc.NewClient(cfg.HTTPAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	client := pb.NewMetricsClient(conn)

	return &GRPCReporter{
		config: cfg,
		client: client,
		conn:   conn,
	}, nil
}

func (r *GRPCReporter) Close() error {
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

func (r *GRPCReporter) ReportMetrics(ms metrics.MetricsCollection) {
	metrics := make([]*pb.Metric, 0, len(ms.GaugeMetrics)+len(ms.CounterMetrics))

	for key, value := range ms.GaugeMetrics {
		metric := &pb.Metric{
			Key:  key,
			Type: pb.Metric_MTYPE_GAUGE,
			MValue: &pb.Metric_Value{
				Value: value,
			},
		}
		metrics = append(metrics, metric)
	}

	for key, value := range ms.CounterMetrics {
		metric := &pb.Metric{
			Key:  key,
			Type: pb.Metric_MTYPE_COUNTER,
			MValue: &pb.Metric_Delta{
				Delta: value,
			},
		}
		metrics = append(metrics, metric)
	}

	if len(metrics) == 0 {
		return
	}

	req := &pb.UpdateMetricsBulkRequest{
		Metrics: metrics,
	}

	retryIntervals := []time.Duration{time.Second, time.Second * 2, time.Second * 5}

	for _, interval := range retryIntervals {
		err := r.sendMetrics(req)
		if err != nil && (errors.Is(err, context.DeadlineExceeded) || isConnectionRefused(err)) {
			time.Sleep(interval)
			continue
		}
		break
	}
}

func (r *GRPCReporter) sendMetrics(req *pb.UpdateMetricsBulkRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	if r.config.Key != "" {
		jsonData, err := json.Marshal(req)
		if err != nil {
			return fmt.Errorf("failed to marshal metrics for signing: %w", err)
		}

		signature := signer.SignPayloadWithKey(jsonData, []byte(r.config.Key))
		ctx = metadata.AppendToOutgoingContext(ctx, domain.SignatureHeader, hex.EncodeToString(signature))
	}

	resp, err := r.client.UpdateMetricsBulk(ctx, req)
	if err != nil {
		return fmt.Errorf("gRPC call failed: %w", err)
	}

	if !resp.Received {
		if resp.Error != nil {
			return errors.New(*resp.Error)
		}
		return errors.New("server did not acknowledge receipt of metrics")
	}

	log.Printf("Successfully sent metrics via gRPC, server response: %v", resp)
	return nil
}
