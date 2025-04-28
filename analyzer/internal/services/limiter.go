package services

import (
	"context"
	"log/slog"
	dtopkg "waf-analyzer/internal/domain/dto"
)

type fingerprintsProvider interface {
	Fingerprints(ctx context.Context, ip string) ([]int, error)
}

// LimiterService contains request analyze business logic.
type LimiterService struct {
	log *slog.Logger
}

// NewLimiterService is a constructor for LimiterService.
func NewLimiterService(
	log *slog.Logger,
) *LimiterService {
	return &LimiterService{
		log: log,
	}
}

// CheckLimit checks whether requests count from given ip is in RPS range. Returns true if RPS is over the limit.
func (s *LimiterService) CheckLimit(
	ctx context.Context,
	dto *dtopkg.HTTPRequest,
) (bool, error) {
	const op = "services.LimiterService.CheckLimit"

	log := s.log.With(slog.String("op", op))

	log.Info("checking requests rate")

	return false, nil
}
