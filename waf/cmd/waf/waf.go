package main

import (
	"os"
	"os/signal"
	"syscall"
	"waf-waf/internal/app"
	"waf-waf/internal/config"
	"waf-waf/internal/logger"
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
	)

	go func() {
		application.GRPCServer.MustRun()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	application.GRPCServer.Stop()

	log.Info("gracefully stopped")
}
