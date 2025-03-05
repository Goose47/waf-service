// Package redis provides API to interact with redis.
package redis

import (
	"context"
	"fmt"
	redislib "github.com/go-redis/redis/v8"
	"log/slog"
	"net"
	"strconv"
)

// Redis provides API to interact with redis.
type Redis struct {
	log    *slog.Logger
	client *redislib.Client
}

// New is a constructor for Redis.
func New(
	log *slog.Logger,
	host string,
	port int,
	pass string,
) *Redis {
	client := redislib.NewClient(&redislib.Options{
		Addr:     net.JoinHostPort(host, strconv.Itoa(port)),
		Password: pass,
		DB:       0,
	})
	return &Redis{
		log:    log,
		client: client,
	}
}

// Close closes redis connection.
func (r *Redis) Close() error {
	return fmt.Errorf("redis.Close: %w", r.client.Close())
}

// Fingerprint indicates whether given ip is present in redis.
func (r *Redis) Fingerprint(ctx context.Context, ip string) (bool, error) {
	const op = "storage.redis.Fingerprint"

	log := r.log.With(slog.String("op", op), slog.String("ip", ip))

	log.Info("checking ip existence")

	res, err := r.client.Exists(ctx, ip).Result()

	if err != nil {
		log.Error("failed to check ip existence", slog.Any("error", err.Error()))

		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("ip existence checked", slog.Int64("res", res))

	return res > 0, nil
}

// Save saves ip in redis.
func (r *Redis) Save(ctx context.Context, ip string) error {
	const op = "storage.redis.SetFingerprint"

	log := r.log.With(slog.String("op", op), slog.String("ip", ip))

	log.Info("trying to save ip")

	err := r.client.Set(ctx, ip, 1, 0).Err()

	if err != nil {
		log.Error("failed to save ip", slog.Any("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("ip saved successfully")

	return nil
}
