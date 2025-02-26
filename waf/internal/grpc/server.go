package wafgrpc

import (
	"context"
	gen "github.com/Goose47/wafpb/gen/go/waf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	dtopkg "waf-waf/internal/domain/dto"
)

type serverAPI struct {
	gen.UnimplementedWAFServer
	waf waf
}

type waf interface {
	Analyze(ctx context.Context, dto *dtopkg.HTTPRequest) (float32, error)
}

func Register(gRPCServer *grpc.Server, waf waf) {
	gen.RegisterWAFServer(gRPCServer, &serverAPI{waf: waf})
}

func (s *serverAPI) Analyze(
	ctx context.Context,
	in *gen.AnalyzeRequest,
) (*gen.AnalyzeResponse, error) {
	if in.ClientIp == "" {
		return nil, status.Error(codes.InvalidArgument, "client ip is required")
	}
	if in.ClientPort == "" {
		return nil, status.Error(codes.InvalidArgument, "client port is required")
	}
	if in.ServerIp == "" {
		return nil, status.Error(codes.InvalidArgument, "server ip is required")
	}
	if in.ServerPort == "" {
		return nil, status.Error(codes.InvalidArgument, "server port is required")
	}
	if in.Uri == "" {
		return nil, status.Error(codes.InvalidArgument, "uri is required")
	}
	if in.Method == "" {
		return nil, status.Error(codes.InvalidArgument, "http method is required")
	}
	if in.Proto == "" {
		return nil, status.Error(codes.InvalidArgument, "http protocol message is required")
	}

	dto := dtopkg.NewHTTPRequest(in)

	res, err := s.waf.Analyze(ctx, dto)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to analyze request")
	}

	return &gen.AnalyzeResponse{
		AttackProbability: res,
	}, nil
}
