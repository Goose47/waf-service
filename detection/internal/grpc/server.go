// Package detectiongrpc contains transport layer logic for detection service.
package detectiongrpc

import (
	"context"
	"fmt"
	gen "github.com/Goose47/wafpb/gen/go/detection"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	gen.UnimplementedDetectionServer
	detection detection
}

type detection interface {
	CheckIP(ctx context.Context, ip string) (bool, error)
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

	isSuspicious, err := s.detection.CheckIP(ctx, in.Ip)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, status.Error(codes.Internal, "failed to check ip"))
	}

	return &gen.CheckIPResponse{
		IsSuspicious: isSuspicious,
	}, nil
}
