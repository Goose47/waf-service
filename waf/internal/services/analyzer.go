package services

import (
	"context"
	"encoding/json"
	"fmt"
	gen "github.com/Goose47/wafpb/gen/go/analyzer"
	"github.com/IBM/sarama"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net"
	"strconv"
	dtopkg "waf-waf/internal/domain/dto"
)

type AnalyzerService struct {
	log      *slog.Logger
	client   gen.AnalyzerClient
	producer sarama.SyncProducer
	topic    string
}

func MustCreateAnalyzerService(
	log *slog.Logger,
	host string,
	port int,
	producer sarama.SyncProducer,
	topic string,
) *AnalyzerService {
	gRPCAddress := net.JoinHostPort(host, strconv.Itoa(port))
	cc, err := grpc.NewClient(gRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Errorf("failed to connect to grpc server: %v", err))
	}

	client := gen.NewAnalyzerClient(cc)

	return &AnalyzerService{
		log:      log,
		client:   client,
		producer: producer,
		topic:    topic,
	}
}

// Analyze sends request to analyzer service to check whether given http request contains an attack.
func (s *AnalyzerService) Analyze(ctx context.Context, dto *dtopkg.HTTPRequest) (bool, error) {
	const op = "services.AnalyzerService.Analyze"
	log := s.log.With(slog.String("op", op))

	log.Info("analyzing http request")

	res, err := s.client.Analyze(ctx, dto.ToAnalyzeRequest())

	if err != nil {
		log.Error("failed to analyze request", slog.Any("error", err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("request analyzed successfully", slog.Bool("is_attack", res.IsAttack))

	return res.IsAttack, nil
}

// Publish publishes given http request to be analyzed in background.
func (s *AnalyzerService) Publish(ctx context.Context, dto *dtopkg.HTTPRequest) error {
	const op = "services.AnalyzerService.Publish"
	log := s.log.With(slog.String("op", op))

	log.Info("publishing http request")

	marshalled, err := json.Marshal(dto)
	if err != nil {
		log.Error("failed to marshal http request", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	msg := &sarama.ProducerMessage{
		Topic: s.topic,
		Value: sarama.StringEncoder(marshalled),
	}

	partition, offset, err := s.producer.SendMessage(msg)
	if err != nil {
		log.Error("failed to send message", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info(fmt.Sprintf("sent message to partition %d at offset %d\n", partition, offset))
	log.Info("http request published successfully")

	return nil
}
