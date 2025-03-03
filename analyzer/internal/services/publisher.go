package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log/slog"
)

// PublisherService contains kafka publishing messages logic.
type PublisherService struct {
	log      *slog.Logger
	producer sarama.SyncProducer
	topic    string
}

// NewPublisherService is a constructor for PublisherService.
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

// Close closes sarama producer.
func (s *PublisherService) Close() error {
	return fmt.Errorf("services.Publisher.Close: %w", s.producer.Close())
}

// Publish publishes attacker's ip to kafka.
func (s *PublisherService) Publish(
	_ context.Context,
	ip string,
) error {
	const op = "services.PublisherService.Publish"

	log := s.log.With(slog.String("op", op))

	log.Info("trying to publish attacker's ip")

	type message struct {
		IP string `json:"ip"`
	}

	marshalled, err := json.Marshal(message{IP: ip})
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

	log.Info(fmt.Sprintf("sent message to partition %d at offset %d\n", partition, offset))
	log.Info("ip published successfully")

	return nil
}
