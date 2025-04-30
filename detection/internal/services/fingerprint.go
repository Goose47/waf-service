// Package services contains application business logic.
package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"waf-detection/internal/domain/dto"
	"waf-detection/internal/storage"
)

type clientProvider interface {
	Client(ctx context.Context, ip string) (*dto.Client, error)
}
type clientSaver interface {
	Save(ctx context.Context, client *dto.Client) error
}

// FingerprintService contains fingerprints business logic.
type FingerprintService struct {
	log            *slog.Logger
	clientProvider clientProvider
	clientSaver    clientSaver
}

// NewFingerprintService is a constructor for FingerprintService.
func NewFingerprintService(
	log *slog.Logger,
	clientProvider clientProvider,
	clientSaver clientSaver,
) *FingerprintService {
	return &FingerprintService{
		log:            log,
		clientProvider: clientProvider,
		clientSaver:    clientSaver,
	}
}

// CheckIP checks if client is suspicious and saves another fingerprint on each call.
func (s *FingerprintService) CheckIP(
	ctx context.Context,
	ip string,
	fingerprint time.Time,
) (bool, error) {
	const op = "services.fingerprint.CheckIP"

	log := s.log.With(slog.String("op", op), slog.String("ip", ip))

	log.Info("trying to get client")

	client, err := s.clientProvider.Client(ctx, ip)

	if errors.Is(err, storage.ErrNotFound) {
		client = &dto.Client{
			IP: ip,
		}
	} else if err != nil {
		log.Error("failed to get client", slog.Any("error", err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	client.Fingerprints = append(client.Fingerprints, fingerprint)
	err = s.clientSaver.Save(ctx, client)

	if err != nil {
		log.Error("failed to save client", slog.Any("error", err))

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return client.IsSuspicious, nil
}

// MarkIPSuspicious marks client as suspicious.
func (s *FingerprintService) MarkIPSuspicious(ctx context.Context, ip string) error {
	const op = "services.fingerprint.MarkIPSuspicious"

	log := s.log.With(slog.String("op", op), slog.String("ip", ip))

	log.Info("trying to get client")
	client, err := s.clientProvider.Client(ctx, ip)

	if errors.Is(err, storage.ErrNotFound) {
		client = &dto.Client{
			IP: ip,
		}
	} else if err != nil {
		log.Error("failed to get client", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("trying to save client")
	client.IsSuspicious = true
	err = s.clientSaver.Save(ctx, client)

	if err != nil {
		log.Error("failed to save client", slog.Any("error", err))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("client saved successfully")

	return nil
}

// Fingerprints retrieves all ips fingerprints after given time. Deletes all fingerprints before given time.
func (s *FingerprintService) Fingerprints(
	ctx context.Context,
	ip string,
	after time.Time,
) ([]time.Time, error) {
	const op = "services.fingerprint.Fingerprints"

	log := s.log.With(slog.String("op", op), slog.String("ip", ip))

	log.Info("trying to get client")

	client, err := s.clientProvider.Client(ctx, ip)

	if errors.Is(err, storage.ErrNotFound) {
		client = &dto.Client{
			IP: ip,
		}
	} else if err != nil {
		log.Error("failed to get client", slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	offset := len(client.Fingerprints)
	for i, timestamp := range client.Fingerprints {
		if timestamp.After(after) {
			offset = i
			break
		}
	}

	client.Fingerprints = client.Fingerprints[offset:]
	err = s.clientSaver.Save(ctx, client)

	if err != nil {
		log.Error("failed to save client", slog.Any("error", err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return client.Fingerprints, nil
}
