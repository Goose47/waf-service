package kafka

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"log/slog"
	dtopkg "waf-analyzer/internal/domain/dto"
)

type Analyzer interface {
	Analyze(ctx context.Context, dto *dtopkg.HTTPRequest) (bool, error)
}

type consumerHandler struct {
	ctx      context.Context
	log      *slog.Logger
	analyzer Analyzer
}

func newConsumerHandler(
	ctx context.Context,
	log *slog.Logger,
	analyzer Analyzer,
) *consumerHandler {
	return &consumerHandler{
		ctx:      ctx,
		log:      log,
		analyzer: analyzer,
	}
}

func (h *consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}
func (h *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
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

		var parsedMessage dtopkg.HTTPRequest
		err := json.Unmarshal(msg.Value, &parsedMessage)
		if err != nil {
			log.Error("failed to unmarshal message", slog.Any("error", err))
			continue
		}

		log.Info("unmarshalled message")
		log.Info("analyzing incoming http request")

		_, err = h.analyzer.Analyze(h.ctx, &parsedMessage)
		if err != nil {
			log.Error("failed to analyze request", slog.Any("error", err))
			continue
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
