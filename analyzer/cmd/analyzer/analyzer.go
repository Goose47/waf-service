package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"waf-analyzer/internal/app"
	"waf-analyzer/internal/config"
	"waf-analyzer/internal/logger"
)

func main() {
	cfg := config.MustLoadPath("./config/local.yaml")
	log := logger.New(cfg.Env)

	application := app.New(
		log,
		cfg.GRPC.Port,
		cfg.Kafka.Host,
		cfg.Kafka.Port,
		cfg.Kafka.Topic,
		cfg.Kafka.DetectionTopic,
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
	application.Poller.Close()
	_ = application.Producer.Close()

	log.Info("gracefully stopped")
}
