// Package app provides functions to create application instance and run it.
package app

import (
	"io"
	"log/slog"
	grpcapp "waf-detection/internal/app/grpc"
	kafkapkg "waf-detection/internal/kafka"
	"waf-detection/internal/services"
	redispkg "waf-detection/internal/storage/redis"
)

// App represents application and contains all its dependencies.
type App struct {
	GRPCServer *grpcapp.App
	Poller     *kafkapkg.Poller
	Redis      io.Closer
}

// New is constructor for APP. Creates all dependencies and returns app instance.
func New(
	log *slog.Logger,
	grpcPort int,
	redisHost string,
	redisPort int,
	redisPass string,
	kafkaHost string,
	kafkaPort int,
	kafkaTopic string,
) *App {
	redis := redispkg.New(log, redisHost, redisPort, redisPass)
	fingerprintService := services.NewFingerprintService(log, redis, redis)

	grpcApp := grpcapp.New(log, fingerprintService, grpcPort)
	kafka := kafkapkg.MustCreate(log, kafkaHost, kafkaPort, kafkaTopic, fingerprintService)

	return &App{
		GRPCServer: grpcApp,
		Poller:     kafka,
		Redis:      redis,
	}
}
