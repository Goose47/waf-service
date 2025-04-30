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
)

// DetectionService contains business logic related to detection service.
type DetectionService struct {
	log    *slog.Logger
	client gen.DetectionClient
}

// MustCreateDetectionService is a constructor for DetectionService.
func MustCreateDetectionService(
	log *slog.Logger,
	host string,
	port int,
) *DetectionService {
	gRPCAddress := net.JoinHostPort(host, strconv.Itoa(port))
	cc, err := grpc.NewClient(gRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Errorf("failed to connect to grpc server: %w", err))
	}

	client := gen.NewDetectionClient(cc)

	return &DetectionService{
		log:    log,
		client: client,
	}
}

// IsSuspicious calls detection service to check whether given ip is suspicious.
func (d *DetectionService) IsSuspicious(ctx context.Context, ip string, timestamp time.Time) (bool, error) {
	const op = "services.detection.IsSuspicious"
	log := d.log.With(slog.String("op", op), slog.String("ip", ip))

	log.Info("checking ip")

	res, err := d.client.CheckIP(ctx, &gen.CheckIPRequest{
		Ip:        ip,
		Timestamp: timestamppb.New(timestamp),
	})

	if err != nil {
		log.Error("failed to check ip", slog.Any("error", err))

		return true, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("ip checked successfully", slog.Bool("is_suspicious", res.IsSuspicious))

	return res.IsSuspicious, nil
}
