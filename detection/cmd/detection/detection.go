// Package main runs the application.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"waf-detection/internal/app"
	"waf-detection/internal/config"
	"waf-detection/internal/logger"
)

func main() {
	cfg := config.MustLoadPath("./config/local.yaml")
	log := logger.New(cfg.Env)

	application := app.New(
		log,
		cfg.GRPC.Port,
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Pass,
		cfg.Kafka.Host,
		cfg.Kafka.Port,
		cfg.Kafka.Topic,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		application.GRPCServer.MustRun()
	}()
	go func() {
		application.Poller.MustPoll(ctx)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	cancel()
	application.GRPCServer.Stop()
	if err := application.Poller.Close(); err != nil {
		log.Error("failed to close kafka poller", slog.Any("error", err))
	}
	if err := application.Redis.Close(); err != nil {
		log.Error("failed to close redis connection", slog.Any("error", err))
	}

	log.Info("gracefully stopped")
}
