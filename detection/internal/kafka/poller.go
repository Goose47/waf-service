package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"log/slog"
	"net"
	"strconv"
	"time"
	"waf-detection/internal/lib/random"
)

type Poller struct {
	log    *slog.Logger
	host   string
	port   int
	topic  string
	saver  FingerprintSaver
	client sarama.ConsumerGroup
}

func MustCreate(
	log *slog.Logger,
	host string,
	port int,
	topic string,
	saver FingerprintSaver,
) *Poller {
	poller, err := New(log, host, port, topic, saver)
	if err != nil {
		panic(err)
	}
	return poller
}

func New(
	log *slog.Logger,
	host string,
	port int,
	topic string,
	saver FingerprintSaver,
) (*Poller, error) {
	return &Poller{
		log:   log,
		host:  host,
		port:  port,
		topic: topic,
		saver: saver,
	}, nil
}

// MustPoll creates consumer group, panics on error then consumes messages and saves fingerprints
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

	handler := newConsumerHandler(ctx, log, p.saver)

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

func (p *Poller) Close() {
	p.client.Close()
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
