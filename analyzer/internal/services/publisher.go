package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log/slog"
)

type PublisherService struct {
	log      *slog.Logger
	producer sarama.SyncProducer
	topic    string
}

func NewPublisherService(
	log *slog.Logger,
	producer sarama.SyncProducer,
	topic string,
) *PublisherService {
	return &PublisherService{
		log:      log,
		producer: producer,
		topic:    topic,
	}
}

// Close closes sarama producer
func (s *PublisherService) Close() error {
	return s.producer.Close()
}

// Publish publishes attacker's ip to kafka.
func (s *PublisherService) Publish(
	ctx context.Context,
	ip string,
) error {
	const op = "services.PublisherService.Publish"

	log := s.log.With(slog.String("op", op))

	log.Info("trying to publish attacker's ip")

	type message struct {
		Ip string `json:"ip"`
	}

	marshalled, err := json.Marshal(message{Ip: ip})
	if err != nil {
		log.Error("failed to marshal message", slog.Any("error", err))
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

	log.Info("sent message to partition %d at offset %d\n", partition, offset)
	log.Info("ip published successfully")

	return nil
}
