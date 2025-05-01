package interceptors

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/frolmr/metrics/internal/domain"
	pb "github.com/frolmr/metrics/pkg/proto/metrics"
	"github.com/frolmr/metrics/pkg/signer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func NewSignatureInterceptor(signKey string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if signKey == "" {
			return handler(ctx, req)
		}

		if info.FullMethod != "/metrics.Metrics/UpdateMetricsBulk" {
			return handler(ctx, req)
		}

		updateReq, ok := req.(*pb.UpdateMetricsBulkRequest)
		if !ok {
			return nil, status.Errorf(codes.Internal, "invalid request type")
		}

		if err := validateSignature(ctx, signKey, updateReq); err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "signature validation failed: %v", err)
		}

		return handler(ctx, req)
	}
}

func validateSignature(ctx context.Context, signKey string, req *pb.UpdateMetricsBulkRequest) error {
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

	expectedSignature := signer.SignPayloadWithKey(jsonData, []byte(signKey))
	receivedSignature, err := hex.DecodeString(signatureHeaders[0])
	if err != nil {
		return fmt.Errorf("invalid signature format: %w", err)
	}

	if !bytes.Equal(expectedSignature, receivedSignature) {
		return errors.New("signature mismatch")
	}

	return nil
}
