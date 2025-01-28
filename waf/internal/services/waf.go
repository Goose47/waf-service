package services

import (
	"context"
	"fmt"
	"log/slog"
)

type WAFService struct {
	log       *slog.Logger
	detection IsSuspiciousProvider
}

type IsSuspiciousProvider interface {
	IsSuspicious(ctx context.Context, ip string) (bool, error)
}

func NewWAFService(
	log *slog.Logger,
	detection IsSuspiciousProvider,
) *WAFService {
	return &WAFService{
		log:       log,
		detection: detection,
	}
}

func (s *WAFService) Analyze(ctx context.Context, request []byte, ip string) (float32, error) {
	const op = "services.waf.Analyze"
	log := s.log.With(slog.String("op", op))

	log.Info("checking ip", slog.String("ip", ip))

	isSuspicious, err := s.detection.IsSuspicious(ctx, ip)

	if err != nil {
		log.Info("failed to check ip", slog.String("ip", ip))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("ip checked successfully", slog.String("ip", ip), slog.Bool("is_suspicious", isSuspicious))

	if isSuspicious {
		// analyze inline
		//modsec
		//rate limit

		return 0.5, nil
	}

	//send kafka message to analyze request
	//return 0
	return 0, nil
}
