package wafgrpc

import (
	"context"
	gen "github.com/Goose47/wafpb/gen/go/waf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	gen.UnimplementedWAFServer
	waf WAF
}

type WAF interface {
	Analyze(ctx context.Context, request []byte, ip string) (float32, error)
}

func Register(gRPCServer *grpc.Server, waf WAF) {
	gen.RegisterWAFServer(gRPCServer, &serverAPI{waf: waf})
}

func (s *serverAPI) Analyze(
	ctx context.Context,
	in *gen.AnalyzeRequest,
) (*gen.AnalyzeResponse, error) {
	//todo add missing parameters
	if in.Ip == "" {
		return nil, status.Error(codes.InvalidArgument, "ip is required")
	}

	res, err := s.waf.Analyze(ctx, in.Request, in.Ip)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to analyze request")
	}

	return &gen.AnalyzeResponse{
		AttackProbability: res,
	}, nil
}
