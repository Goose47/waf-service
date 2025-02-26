package analyzergrpc

import (
	"context"
	gen "github.com/Goose47/wafpb/gen/go/analyzer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	dtopkg "waf-analyzer/internal/domain/dto"
)

type serverAPI struct {
	gen.UnimplementedAnalyzerServer
	analyzer Analyzer
}

type Analyzer interface {
	Analyze(
		ctx context.Context,
		dto *dtopkg.HTTPRequest,
	) (bool, error)
}

func Register(gRPCServer *grpc.Server, analyzer Analyzer) {
	gen.RegisterAnalyzerServer(gRPCServer, &serverAPI{analyzer: analyzer})
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

	dto := dtopkg.NewHTTPRequest(
		in.ClientIp,
		in.ClientPort,
		in.ServerIp,
		in.ServerPort,
		in.Uri,
		in.Method,
		in.Proto,
		in.Headers,
		in.QueryParams,
		in.BodyParams,
	)

	isAttack, err := s.analyzer.Analyze(
		ctx,
		dto,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to check ip")
	}

	return &gen.AnalyzeResponse{
		IsAttack: isAttack,
	}, nil
}
