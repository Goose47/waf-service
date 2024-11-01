package services

import (
	"context"
	"fmt"
	"log/slog"
)

type FingerprintProvider interface {
	Fingerprint(ctx context.Context, ip string) (bool, error)
}

type FingerprintService struct {
	log                 *slog.Logger
	fingerprintProvider FingerprintProvider
}

func NewFingerprintService(
	log *slog.Logger,
	provider FingerprintProvider,
) *FingerprintService {
	return &FingerprintService{
		log:                 log,
		fingerprintProvider: provider,
	}
}

// CheckIP checks if ip is present in database
func (s *FingerprintService) CheckIP(ctx context.Context, ip string) (bool, error) {
	const op = "services.fingerprint.CheckIP"

	log := s.log.With(slog.String("op", op), slog.String("ip", ip))

	log.Info("trying to check ip")

	// ip is considered suspicious if it is present in db
	res, err := s.fingerprintProvider.Fingerprint(ctx, ip)

	if err != nil {
		log.Error("failed to check ip", slog.Any("error", err))

		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info(fmt.Sprintf("ip checked successfully: %t", res))

	return res, nil
}
