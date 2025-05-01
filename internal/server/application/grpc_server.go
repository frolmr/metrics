package application

import (
	"context"

	"github.com/frolmr/metrics/internal/domain"
	"github.com/frolmr/metrics/internal/server/storage"
	pb "github.com/frolmr/metrics/pkg/proto/metrics"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServer
	stor storage.Repository
}

func NewMetricsServer(stor storage.Repository) *MetricsServer {
	return &MetricsServer{
		stor: stor,
	}
}

func (s *MetricsServer) UpdateMetricsBulk(ctx context.Context, in *pb.UpdateMetricsBulkRequest) (*pb.Ack, error) {
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

	return &pb.Ack{Received: true}, nil
}
