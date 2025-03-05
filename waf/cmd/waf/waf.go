// Package main runs the application.
package main

import (
	"log/slog"
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
		cfg.Detection.Host,
		cfg.Detection.Port,
		cfg.Analyzer.Host,
		cfg.Analyzer.Port,
		cfg.Kafka.Host,
		cfg.Kafka.Port,
		cfg.Kafka.AnalyzerTopic,
	)

	go func() {
		application.GRPCServer.MustRun()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	application.GRPCServer.Stop()
	if err := application.Producer.Close(); err != nil {
		log.Error("failed to close sarama producer", slog.Any("error", err))
	}

	log.Info("gracefully stopped")
}
