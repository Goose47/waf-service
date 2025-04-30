package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	dtopkg "waf-waf/internal/domain/dto"
)

// WAFService contains main application business logic.
type WAFService struct {
	log       *slog.Logger
	detection isSuspiciousProvider
	analyzer  analyzer
	publisher publisher
}

type isSuspiciousProvider interface {
	IsSuspicious(ctx context.Context, ip string, timestamp time.Time) (bool, error)
}

type publisher interface {
	Publish(ctx context.Context, dto *dtopkg.HTTPRequest) error
}

type analyzer interface {
	Analyze(ctx context.Context, dto *dtopkg.HTTPRequest) (bool, error)
}

// NewWAFService is a constructor for WAFService.
func NewWAFService(
	log *slog.Logger,
	detection isSuspiciousProvider,
	analyzer analyzer,
	publisher publisher,
) *WAFService {
	return &WAFService{
		log:       log,
		detection: detection,
		analyzer:  analyzer,
		publisher: publisher,
	}
}

// Analyze contains main app logic. Check ip for fingerprints and calls analyzer inline or sends request to be analyzed
// in background based on fingerprints result.
func (s *WAFService) Analyze(ctx context.Context, dto *dtopkg.HTTPRequest) (float32, error) {
	const op = "services.waf.Analyze"
	log := s.log.With(slog.String("op", op))

	log.Info("checking ip", slog.String("ip", dto.ClientIP))

	isSuspicious, err := s.detection.IsSuspicious(ctx, dto.ClientIP, dto.Timestamp)

	if err != nil {
		log.Info("failed to check ip", slog.String("ip", dto.ClientIP))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("ip checked successfully", slog.String("ip", dto.ClientIP), slog.Bool("is_suspicious", isSuspicious))

	if isSuspicious {
		// Rate limiter is supposed to be there. But it does not.
		res, err := s.analyzer.Analyze(ctx, dto)
		if err != nil {
			log.Error("failed to analyze http request inline", slog.Any("error", err))
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		var probability float32
		if res {
			probability = 1
		}

		return probability, nil
	}

	log.Info("publishing http request")

	if err := s.publisher.Publish(ctx, dto); err != nil {
		log.Error("failed to publish request", slog.Any("error", err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("http request published")

	return 0, nil
}
