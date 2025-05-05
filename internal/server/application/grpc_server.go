package application

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/frolmr/metrics/internal/domain"
	"github.com/frolmr/metrics/internal/server/storage"
	pb "github.com/frolmr/metrics/pkg/proto/metrics"
	"github.com/frolmr/metrics/pkg/signer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServer
	stor    storage.Repository
	signKey string
}

func NewMetricsServer(stor storage.Repository, signKey string) *MetricsServer {
	return &MetricsServer{
		stor:    stor,
		signKey: signKey,
	}
}

func (s *MetricsServer) UpdateMetricsBulk(ctx context.Context, in *pb.UpdateMetricsBulkRequest) (*pb.Ack, error) {
	if s.signKey != "" {
		if err := s.validateSignature(ctx, in); err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "signature validation failed: %v", err)
		}
	}

	var metrics []domain.Metrics

	for _, v := range in.GetMetrics() {
		switch v.GetType() {
		case pb.Metric_MTYPE_COUNTER:
			delta := v.GetDelta()
			metrics = append(metrics, domain.Metrics{ID: v.Key, MType: domain.CounterType, Delta: &delta})
		case pb.Metric_MTYPE_GAUGE:
			value := v.GetValue()
			metrics = append(metrics, domain.Metrics{ID: v.Key, MType: domain.GaugeType, Value: &value})
		}
	}

	err := s.stor.UpdateMetrics(metrics)
	if err != nil {
		return nil, err
	}

	response := pb.Ack{
		Received: true,
	}

	return &response, nil
}

func (s *MetricsServer) validateSignature(ctx context.Context, req *pb.UpdateMetricsBulkRequest) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("no metadata found in request")
	}

	signatureHeaders := md.Get(domain.SignatureHeader)
	if len(signatureHeaders) == 0 {
		return errors.New("no signature header found")
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	expectedSignature := signer.SignPayloadWithKey(jsonData, []byte(s.signKey))
	receivedSignature, err := hex.DecodeString(signatureHeaders[0])
	if err != nil {
		return fmt.Errorf("invalid signature format: %w", err)
	}

	if !bytes.Equal(expectedSignature, receivedSignature) {
		return errors.New("signature mismatch")
	}

	return nil
}
