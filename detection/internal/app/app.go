package app

import (
	"log/slog"
	grpcapp "waf-detection/internal/app/grpc"
	"waf-detection/internal/kafka"
	kafkapkg "waf-detection/internal/kafka"
	"waf-detection/internal/services"
	redispkg "waf-detection/internal/storage/redis"
)

type App struct {
	GRPCServer *grpcapp.App
	Poller     *kafka.Poller
}

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
	}
}
