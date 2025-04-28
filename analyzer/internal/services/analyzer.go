// Package services contains application service layer logic.
package services

import (
	"context"
	"fmt"
	"log/slog"
	dtopkg "waf-analyzer/internal/domain/dto"
)

type analyzer interface {
	Analyze(ctx context.Context, dto *dtopkg.HTTPRequest) (bool, error)
}
type publisher interface {
	Publish(ctx context.Context, ip string) error
}
type limiter interface {
	CheckLimit(ctx context.Context, ip string) (bool, error)
}

// AnalyzerService contains request analyze business logic.
type AnalyzerService struct {
	log       *slog.Logger
	waf       analyzer
	publisher publisher
	limiter   limiter
}

// NewAnalyzerService is a constructor for AnalyzerService.
func NewAnalyzerService(
	log *slog.Logger,
	waf analyzer,
	publisher publisher,
	limiter limiter,
) *AnalyzerService {
	return &AnalyzerService{
		log:       log,
		waf:       waf,
		publisher: publisher,
		limiter:   limiter,
	}
}

// Analyze analyzes whether request contains attacks.
func (s *AnalyzerService) Analyze(
	ctx context.Context,
	dto *dtopkg.HTTPRequest,
) (bool, error) {
	const op = "services.AnalyzerService.Analyze"

	log := s.log.With(slog.String("op", op))

	log.Info("checking requests rate")

	isAttack, err := s.limiter.CheckLimit(ctx, dto.ClientIP)

	if err != nil {
		log.Error("failed to check requests rate", slog.Any("error", err))

		return false, fmt.Errorf("%s: %w", op, err)
	}
	if isAttack {
		log.Warn("request rate is too high")

		return true, nil
	}

	log.Info("trying to check request")

	isAttack, err = s.waf.Analyze(ctx, dto)

	if err != nil {
		log.Error("failed to analyze request", slog.Any("error", err))

		return false, fmt.Errorf("%s: %w", op, err)
	}
	if isAttack {
		if err := s.publisher.Publish(ctx, dto.ClientIP); err != nil {
			log.Error("failed to publish message", slog.Any("error", err))
			return isAttack, fmt.Errorf("%s: %w", op, err)
		}
	}

	log.Info(fmt.Sprintf("request analyzed successfully: %t", isAttack))

	return isAttack, nil
}
