// Package app provides functions to create application instance and run it.
package app

import (
	"fmt"
	"github.com/IBM/sarama"
	"log/slog"
	"time"
	grpcapp "waf-analyzer/internal/app/grpc"
	kafkapkg "waf-analyzer/internal/kafka"
	"waf-analyzer/internal/services"
	wafpkg "waf-analyzer/internal/services/waf"
)

// App represents application and contains all its dependencies.
type App struct {
	GRPCServer *grpcapp.App
	Poller     *kafkapkg.Poller
	Producer   sarama.SyncProducer
}

// New is constructor for APP. Creates all dependencies and returns app instance.
func New(
	log *slog.Logger,
	grpcPort int,
	kafkaHost string,
	kafkaPort int,
	kafkaTopic string,
	kafkaDetectionTopic string,
	detectionHost string,
	detectionPort int,
	maxRequests int,
	per time.Duration,
) *App {
	producer := mustCreateProducer(kafkaHost, kafkaPort)

	publisherService := services.NewPublisherService(log, producer, kafkaDetectionTopic)
	limiterService := services.MustCreatLimiterService(log, detectionHost, detectionPort, maxRequests, per)

	waf := wafpkg.MustCreate(log)
	analyzerService := services.NewAnalyzerService(log, waf, publisherService, limiterService)

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
