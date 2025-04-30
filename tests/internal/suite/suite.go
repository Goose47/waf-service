package suite

import (
	"context"
	"github.com/Goose47/wafpb/gen/go/analyzer"
	"github.com/Goose47/wafpb/gen/go/detection"
	"github.com/Goose47/wafpb/gen/go/waf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"
	"tests/internal/config"
	"time"
)

type Suite struct {
	*testing.T
	Waf        waf.WAFClient
	Detection  detection.DetectionClient
	Analyzer   analyzer.AnalyzerClient
	LimiterCfg config.LimiterConfig
}

func New(
	t *testing.T,
) (*Suite, context.Context) {
	cfg := config.MustLoadPath("../config/config.yml")

	wafClient := waf.NewWAFClient(mustCreateClient(t, cfg.WAF))
	detectionClient := detection.NewDetectionClient(mustCreateClient(t, cfg.Detection))
	analyzerClient := analyzer.NewAnalyzerClient(mustCreateClient(t, cfg.Analyzer))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	return &Suite{
		t,
		wafClient,
		detectionClient,
		analyzerClient,
		cfg.Limiter,
	}, ctx
}

func mustCreateClient(t *testing.T, cfg config.GRPCConfig) *grpc.ClientConn {
	grpcAddress := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	cc, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc connection failed: %v", err)
	}
	return cc
}
