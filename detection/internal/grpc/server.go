// Package detectiongrpc contains transport layer logic for detection service.
package detectiongrpc

import (
	"context"
	"fmt"
	gen "github.com/Goose47/wafpb/gen/go/detection"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type serverAPI struct {
	gen.UnimplementedDetectionServer
	detection detection
}

type detection interface {
	CheckIP(ctx context.Context, ip string, fingerprint time.Time) (bool, error)
	Fingerprints(ctx context.Context, ip string, after time.Time) ([]time.Time, error)
}

// Register associates gRPC server with service layer.
func Register(gRPCServer *grpc.Server, detection detection) {
	gen.RegisterDetectionServer(gRPCServer, &serverAPI{detection: detection})
}

func (s *serverAPI) CheckIP(
	ctx context.Context,
	in *gen.CheckIPRequest,
) (*gen.CheckIPResponse, error) {
	const op = "grpc.CheckIP"
	if in.Ip == "" {
		return nil, fmt.Errorf("%s: %w", op, status.Error(codes.InvalidArgument, "ip is required"))
	}

	isSuspicious, err := s.detection.CheckIP(ctx, in.Ip, in.Timestamp.AsTime())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, status.Error(codes.Internal, "failed to check ip"))
	}

	return &gen.CheckIPResponse{
		IsSuspicious: isSuspicious,
	}, nil
}

func (s *serverAPI) Fingerprints(
	ctx context.Context,
	in *gen.FingerprintsRequest,
) (*gen.FingerprintsResponse, error) {
	const op = "grpc.Fingerprints"
	if in.Ip == "" {
		return nil, fmt.Errorf("%s: %w", op, status.Error(codes.InvalidArgument, "ip is required"))
	}

	timestamps, err := s.detection.Fingerprints(ctx, in.Ip, in.After.AsTime())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, status.Error(codes.Internal, "failed to get fingerprints"))
	}

	pbTimestamps := make([]*timestamppb.Timestamp, len(timestamps))
	for i := range timestamps {
		pbTimestamps[i] = timestamppb.New(timestamps[i])
	}

	return &gen.FingerprintsResponse{
		Timestamps: pbTimestamps,
	}, nil
}
