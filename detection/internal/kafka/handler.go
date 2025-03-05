// Package kafka contains sarama configuration logic.
package kafka

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"log/slog"
)

type fingerprintSaver interface {
	SaveIP(ctx context.Context, ip string) error
}

type consumerHandler struct {
	ctx              context.Context
	log              *slog.Logger
	fingerprintSaver fingerprintSaver
}

func newConsumerHandler(
	ctx context.Context,
	log *slog.Logger,
	saver fingerprintSaver,
) *consumerHandler {
	return &consumerHandler{
		ctx:              ctx,
		log:              log,
		fingerprintSaver: saver,
	}
}

// Setup method exists to comply with sarama ConsumerGroupHandler interface.
func (h *consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup method exists to comply with sarama ConsumerGroupHandler interface.
func (h *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

type message struct {
	IP string `json:"ip"`
}

// ConsumeClaim consumes incoming messages, retrieves IP's and passes to fingerprintSaver to save it in redis.
func (h *consumerHandler) ConsumeClaim(
	sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for msg := range claim.Messages() {
		log := h.log.With(
			slog.String("key", string(msg.Key)),
			slog.String("value", string(msg.Value)),
			slog.Int("partition", int(msg.Partition)),
			slog.Int64("offset", msg.Offset),
		)

		log.Info("received message")
		log.Info("unmarshalling message")

		var parsedMessage message
		err := json.Unmarshal(msg.Value, &parsedMessage)
		if err != nil {
			log.Error("failed to unmarshal message", slog.Any("error", err))
			continue
		}

		log.Info("unmarshalled message", slog.String("ip", parsedMessage.IP))

		log.Info("saving message")

		err = h.fingerprintSaver.SaveIP(h.ctx, parsedMessage.IP)
		if err != nil {
			log.Error("failed to save message", slog.Any("error", err))
			continue
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
