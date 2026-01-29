package orderprocessor

import (
	"context"
	"fmt"

	"order_pipeline/internal/events"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
}

func NewHandler(logger *zap.Logger) *Handler {
	return &Handler{logger: logger.Named("orderprocessor-handler")}
}

func (h *Handler) Handle(ctx context.Context, msg kafka.Message) error {
	env, err := events.DecodeEnvelope(msg.Value)
	if err != nil {
		fmt.Errorf("decode envelope: %w", err)
	}

	switch env.Type {
	case events.OrderCreated:
		ev, err := events.DecodeData[events.OrderCreatedEvent](env)
		if err != nil {
			return fmt.Errorf("decode order.created data: %w", err)
		}
		h.logger.Info("processed order.created",
			zap.String("event_id", env.ID),
			zap.String("order_id", ev.OrderID),
			zap.String("customer_id", ev.CustomerID),
			zap.Float64("total_amount", ev.TotalAmount.InexactFloat64()),
			zap.String("currency", ev.Currency),
		)
		return nil
	default:
		// Fow now unknown event types are treated as handled and committed
		h.logger.Warn("unknown event type; committing anyway",
			zap.String("event_type", string(env.Type)),
			zap.String("event_id", env.ID),
		)
		return nil
	}
}
