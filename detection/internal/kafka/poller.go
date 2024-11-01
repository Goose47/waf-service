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
	client sarama.ConsumerGroup
	topic  string
}

func MustCreate(
	log *slog.Logger,
	host string,
	port int,
	topic string,
) *Poller {
	poller, err := New(log, host, port, topic)
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
) (*Poller, error) {
	client, err := newConsumerGroup(host, port)

	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &Poller{
		log:    log,
		client: client,
		topic:  topic,
	}, nil
}

// Poll consumes messages and saves
func (p *Poller) Poll(
	ctx context.Context,
	saver FingerprintSaver,
) {
	const op = "kafka.Poll"
	log := p.log.With(slog.String("op", op))

	handler := newConsumerHandler(ctx, log, saver)

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
