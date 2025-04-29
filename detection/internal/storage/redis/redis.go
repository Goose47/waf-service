// Package redis provides API to interact with redis.
package redis

import (
	"context"
	"encoding/json"
	"fmt"
	redislib "github.com/go-redis/redis/v8"
	"log/slog"
	"net"
	"strconv"
	"waf-detection/internal/domain/dto"
	"waf-detection/internal/storage"
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

// Client retrieves client from redis.
func (r *Redis) Client(ctx context.Context, ip string) (*dto.Client, error) {
	const op = "storage.redis.Client"
	log := r.log.With(slog.String("op", op), slog.String("ip", ip))

	log.Info("getting client")
	res, err := r.client.Get(ctx, ip).Result()

	if err != nil {
		if err == redislib.Nil {
			log.Warn("client not found")
			return &dto.Client{}, storage.ErrNotFound
		}
		log.Error("failed to get client", slog.Any("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("marshalling client")
	var client dto.Client
	err = json.Unmarshal([]byte(res), &client)

	if err != nil {
		log.Error("failed to unmarshal client", slog.Any("error", err.Error()))
	}

	return &client, nil
}

// Save saves client in redis.
func (r *Redis) Save(ctx context.Context, client *dto.Client) error {
	const op = "storage.redis.SaveClient"
	log := r.log.With(slog.String("op", op), slog.String("ip", client.IP))

	log.Info("marshalling client")
	marshaled, err := json.Marshal(client)

	if err != nil {
		log.Error("failed to marshal client", slog.Any("error", err.Error()))
	}

	log.Info("saving client")
	err = r.client.Set(ctx, client.IP, marshaled, 0).Err()

	if err != nil {
		log.Error("failed to save client", slog.Any("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
