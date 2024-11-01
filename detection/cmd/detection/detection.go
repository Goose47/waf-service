package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"waf-detection/internal/config"
	kafkapkg "waf-detection/internal/kafka"
	"waf-detection/internal/logger"
	serverpkg "waf-detection/internal/server"
	"waf-detection/internal/server/controllers"
	"waf-detection/internal/services"
	redispkg "waf-detection/internal/storage/redis"
)

func main() {
	cfg := config.MustLoadPath("./config/local.yaml")

	log := logger.New(cfg.Env)

	redis := redispkg.New(log, cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Pass)

	fingerprintService := services.NewFingerprintService(log, redis)
	fingerprintController := controllers.NewFingerprintController(fingerprintService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kafka := kafkapkg.MustCreate(log, cfg.Kafka.Host, cfg.Kafka.Port, cfg.Kafka.Topic)
	go func() {
		kafka.Poll(ctx, redis)
	}()

	server := serverpkg.New(fingerprintController)
	go func() {
		err := server.Run(fmt.Sprintf(":%d", cfg.Port))
		log.Error("server stopped", slog.Any("error", err))
		cancel()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigs:
	case <-ctx.Done():
	}

	cancel()
	kafka.Close()
	redis.Close()

	log.Info("gracefully stopped")
}
