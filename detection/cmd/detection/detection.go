package main

import (
	"context"
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
		application.Poller.Poll(ctx)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	cancel()
	application.GRPCServer.Stop()
	application.Poller.Close()
	//todo close redis conn

	log.Info("gracefully stopped")
}
