package detectiongrpc

import (
	"context"
	gen "github.com/Goose47/wafpb/gen/go/detection"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	gen.UnimplementedDetectionServer
	detection Detection
}

type Detection interface {
	CheckIP(ctx context.Context, ip string) (bool, error)
}

func Register(gRPCServer *grpc.Server, detection Detection) {
	gen.RegisterDetectionServer(gRPCServer, &serverAPI{detection: detection})
}

func (s *serverAPI) CheckIP(
	ctx context.Context,
	in *gen.CheckIPRequest,
) (*gen.CheckIPResponse, error) {
	if in.Ip == "" {
		return nil, status.Error(codes.InvalidArgument, "ip is required")
	}

	isSuspicious, err := s.detection.CheckIP(ctx, in.Ip)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to check ip")
	}

	return &gen.CheckIPResponse{
		IsSuspicious: isSuspicious,
	}, nil
}
