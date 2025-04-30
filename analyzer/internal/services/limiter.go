package services

import (
	"context"
	"fmt"
	gen "github.com/Goose47/wafpb/gen/go/detection"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
	"net"
	"strconv"
	"time"
	dtopkg "waf-analyzer/internal/domain/dto"
)

// LimiterService contains request analyze business logic.
type LimiterService struct {
	log         *slog.Logger
	client      gen.DetectionClient
	maxRequests int
	per         time.Duration
}

// MustCreatLimiterService is a constructor for LimiterService. Panics on error.
func MustCreatLimiterService(
	log *slog.Logger,
	host string,
	port int,
	maxRequests int,
	per time.Duration,
) *LimiterService {
	gRPCAddress := net.JoinHostPort(host, strconv.Itoa(port))
	cc, err := grpc.NewClient(gRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Errorf("failed to connect to grpc server: %w", err))
	}

	client := gen.NewDetectionClient(cc)

	return &LimiterService{
		log:         log,
		client:      client,
		maxRequests: maxRequests,
		per:         per,
	}
}

// CheckLimit checks whether requests count from given ip is in RPS range. Returns true if RPS is over the limit.
func (s *LimiterService) CheckLimit(
	ctx context.Context,
	dto *dtopkg.HTTPRequest,
) (bool, error) {
	const op = "services.LimiterService.CheckLimit"

	log := s.log.With(slog.String("op", op))

	log.Info("getting fingerprints")

	after := time.Now().Add(-s.per)

	fingerprints, err := s.client.Fingerprints(ctx, &gen.FingerprintsRequest{
		Ip:    dto.ClientIP,
		After: timestamppb.New(after),
	})
	if err != nil {
		log.Error("failed to get fingerprints", slog.Any("error", err))
		return false, fmt.Errorf("%s: %w", op, err)
	}
	requestCount := 0
	for i := len(fingerprints.Timestamps) - 1; i >= 0; i-- {
		timestamp := fingerprints.Timestamps[i].AsTime()
		if timestamp.Before(after) {
			break
		}
		requestCount++
	}
	if requestCount > s.maxRequests {
		return true, nil
	}

	return false, nil
}
