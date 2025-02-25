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

type AnalyzerService struct {
	log       *slog.Logger
	analyzer  analyzer
	publisher publisher
}

func NewAnalyzerService(
	log *slog.Logger,
	analyzer analyzer,
	publisher publisher,
) *AnalyzerService {
	return &AnalyzerService{
		log:       log,
		analyzer:  analyzer,
		publisher: publisher,
	}
}

// Analyze analyzes whether request contains attacks.
func (s *AnalyzerService) Analyze(
	ctx context.Context,
	dto *dtopkg.HTTPRequest,
) (bool, error) {
	const op = "services.AnalyzerService.Analyze"

	log := s.log.With(slog.String("op", op))

	log.Info("trying to check request")

	isAttack, err := s.analyzer.Analyze(ctx, dto)

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
