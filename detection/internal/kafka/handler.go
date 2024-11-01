package kafka

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"log/slog"
)

type FingerprintSaver interface {
	Save(ctx context.Context, ip string) error
}

type consumerHandler struct {
	ctx              context.Context
	log              *slog.Logger
	fingerprintSaver FingerprintSaver
}

func newConsumerHandler(
	ctx context.Context,
	log *slog.Logger,
	saver FingerprintSaver,
) *consumerHandler {
	return &consumerHandler{
		ctx:              ctx,
		log:              log,
		fingerprintSaver: saver,
	}
}

func (h *consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}
func (h *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

type message struct {
	Ip string `json:"ip"`
}

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

		log.Info("unmarshalled message", slog.String("ip", parsedMessage.Ip))

		log.Info("saving message")

		err = h.fingerprintSaver.Save(h.ctx, parsedMessage.Ip)
		if err != nil {
			log.Error("failed to save message", slog.Any("error", err))
			continue
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
