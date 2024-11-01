package services

import (
	"context"
	"fmt"
	"log/slog"
)

type FingerprintProvider interface {
	Fingerprint(ctx context.Context, ip string) (bool, error)
}
type FingerprintSaver interface {
	Save(ctx context.Context, ip string) error
}

type FingerprintService struct {
	log                 *slog.Logger
	fingerprintProvider FingerprintProvider
	fingerprintSaver    FingerprintSaver
}

func NewFingerprintService(
	log *slog.Logger,
	provider FingerprintProvider,
	saver FingerprintSaver,
) *FingerprintService {
	return &FingerprintService{
		log:                 log,
		fingerprintProvider: provider,
		fingerprintSaver:    saver,
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

// SaveIP saves ip in database
func (s *FingerprintService) SaveIP(ctx context.Context, ip string) error {
	const op = "services.fingerprint.SaveIP"

	log := s.log.With(slog.String("op", op), slog.String("ip", ip))

	log.Info("trying to save ip")

	err := s.fingerprintSaver.Save(ctx, ip)

	if err != nil {
		log.Error("failed to save ip", slog.Any("error", err))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info(fmt.Sprintf("ip saved successfully"))

	return nil
}
