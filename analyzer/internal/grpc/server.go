// Package analyzergrpc contains transport layer logic for analyzer service.
package analyzergrpc

import (
	"context"
	"fmt"
	gen "github.com/Goose47/wafpb/gen/go/analyzer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	dtopkg "waf-analyzer/internal/domain/dto"
)

type serverAPI struct {
	gen.UnimplementedAnalyzerServer
	analyzer analyzer
}

type analyzer interface {
	Analyze(
		ctx context.Context,
		dto *dtopkg.HTTPRequest,
	) (bool, error)
}

// Register associates gRPC server with service layer.
func Register(gRPCServer *grpc.Server, analyzer analyzer) {
	gen.RegisterAnalyzerServer(gRPCServer, &serverAPI{analyzer: analyzer})
}

func (s *serverAPI) Analyze(
	ctx context.Context,
	in *gen.AnalyzeRequest,
) (*gen.AnalyzeResponse, error) {
	const op = "grpc.Analyze"

	if in.ClientIp == "" {
		return nil, fmt.Errorf("%s: %w", op, status.Error(codes.InvalidArgument, "client ip is required"))
	}
	if in.ClientPort == "" {
		return nil, fmt.Errorf("%s: %w", op, status.Error(codes.InvalidArgument, "client port is required"))
	}
	if in.ServerIp == "" {
		return nil, fmt.Errorf("%s: %w", op, status.Error(codes.InvalidArgument, "server ip is required"))
	}
	if in.ServerPort == "" {
		return nil, fmt.Errorf("%s: %w", op, status.Error(codes.InvalidArgument, "server port is required"))
	}
	if in.Uri == "" {
		return nil, fmt.Errorf("%s: %w", op, status.Error(codes.InvalidArgument, "uri is required"))
	}
	if in.Method == "" {
		return nil, fmt.Errorf("%s: %w", op, status.Error(codes.InvalidArgument, "http method is required"))
	}
	if in.Proto == "" {
		return nil, fmt.Errorf("%s: %w", op, status.Error(codes.InvalidArgument, "http protocol message is required"))
	}

	dto := dtopkg.NewHTTPRequest(in)

	isAttack, err := s.analyzer.Analyze(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, status.Error(codes.Internal, "failed to check ip"))
	}

	return &gen.AnalyzeResponse{
		IsAttack: isAttack,
	}, nil
}
