package app

import (
	"fmt"
	"github.com/IBM/sarama"
	"log/slog"
	grpcapp "waf-analyzer/internal/app/grpc"
	kafkapkg "waf-analyzer/internal/kafka"
	"waf-analyzer/internal/services"
	wafpkg "waf-analyzer/internal/services/waf"
)

type App struct {
	GRPCServer *grpcapp.App
	Poller     *kafkapkg.Poller
	Producer   sarama.SyncProducer
}

func New(
	log *slog.Logger,
	grpcPort int,
	kafkaHost string,
	kafkaPort int,
	kafkaTopic string,
	kafkaDetectionTopic string,
) *App {
	producer := mustCreateProducer(kafkaHost, kafkaPort)

	publisherService := services.NewPublisherService(log, producer, kafkaDetectionTopic)

	waf := wafpkg.MustCreate(log)
	analyzerService := services.NewAnalyzerService(log, waf, publisherService)

	grpcApp := grpcapp.New(log, analyzerService, grpcPort)
	kafka := kafkapkg.MustCreate(log, kafkaHost, kafkaPort, kafkaTopic, analyzerService)

	return &App{
		GRPCServer: grpcApp,
		Poller:     kafka,
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
