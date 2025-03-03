// Package kafka contains functions to consume and process kafka messages.
package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"log/slog"
	"net"
	"strconv"
	"time"
	"waf-analyzer/internal/lib/random"
)

// Poller listens for new kafka messages and processed them.
type Poller struct {
	analyzer analyzer
	client   sarama.ConsumerGroup
	log      *slog.Logger
	host     string
	topic    string
	port     int
}

// MustCreate creates new Poller and panics on error.
func MustCreate(
	log *slog.Logger,
	host string,
	port int,
	topic string,
	analyzer analyzer,
) *Poller {
	poller, err := New(log, host, port, topic, analyzer)
	if err != nil {
		panic(err)
	}
	return poller
}

// New is a constructor for Poller.
func New(
	log *slog.Logger,
	host string,
	port int,
	topic string,
	analyzer analyzer,
) (*Poller, error) {
	return &Poller{
		log:      log,
		host:     host,
		port:     port,
		topic:    topic,
		analyzer: analyzer,
	}, nil
}

// MustPoll creates consumer group, panics on error then consumes messages and analyzes them.
func (p *Poller) MustPoll(
	ctx context.Context,
) {
	const op = "kafka.Poll"
	log := p.log.With(slog.String("op", op))

	client, err := newConsumerGroup(p.host, p.port)
	if err != nil {
		log.Error("failed to create consumer group", slog.Any("error", err))
	}
	p.client = client

	handler := newConsumerHandler(ctx, log, p.analyzer)

	for {
		err := p.client.Consume(ctx, []string{p.topic}, handler)
		if err != nil {
			log.Error("consume error", slog.Any("error", err))
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

// Close closes publisher.
func (p *Poller) Close() error {
	return fmt.Errorf("services.Publisher.Close: %w", p.client.Close())
}

func newConsumerGroup(
	host string,
	port int,
) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	brokers := []string{net.JoinHostPort(host, strconv.Itoa(port))}

	client, err := sarama.NewConsumerGroup(brokers, random.String(10), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return client, nil
}
