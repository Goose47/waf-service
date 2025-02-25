package app

import (
	"fmt"
	"github.com/IBM/sarama"
	"log/slog"
	grpcapp "waf-waf/internal/app/grpc"
	"waf-waf/internal/services"
)

type App struct {
	GRPCServer *grpcapp.App
	Producer   sarama.SyncProducer
}

func New(
	log *slog.Logger,
	grpcPort int,
	detectionHost string,
	detectionPort int,
	analyzerHost string,
	analyzerPort int,
	kafkaHost string,
	kafkaPort int,
	analyzerTopic string,
) *App {
	detectionService := services.MustCreateDetectionService(log, detectionHost, detectionPort)

	producer := mustCreateProducer(kafkaHost, kafkaPort)
	analyzerService := services.MustCreateAnalyzerService(log, analyzerHost, analyzerPort, producer, analyzerTopic)

	wafService := services.NewWAFService(log, detectionService, analyzerService, analyzerService)

	grpcApp := grpcapp.New(log, wafService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
		Producer:   producer,
	}
}

func mustCreateProducer(
	host string,
	port int,
) sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true

	brokers := []string{fmt.Sprintf("%s:%d", host, port)}
	producer, err := sarama.NewSyncProducer(brokers, config)

	if err != nil {
		panic(fmt.Sprintf("Failed to create producer: %v", err))
	}

	return producer
}
