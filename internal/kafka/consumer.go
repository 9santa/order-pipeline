package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"order_pipeline/internal/model"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	r *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       1e3,  // 1KB
		MaxBytes:       10e6, // 10MB
		MaxWait:        500 * time.Millisecond,
		CommitInterval: 0, // 0 -> will commit manually after successful processing
	})
	return &Consumer{r: r}
}

func (c *Consumer) Close() error { return c.r.Close() }

func (c *Consumer) Run(ctx context.Context, handler func(context.Context, model.Event) error) error {
	for {
		msg, err := c.r.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			return err
		}

		var ev model.Event
		if err := json.Unmarshal(msg.Value, &ev); err != nil {
			_ = c.r.CommitMessages(ctx, msg)
			continue
		}

		// Process the message with handler
		if err := handler(ctx, ev); err != nil {
			time.Sleep(250 * time.Millisecond)
			continue
		}

		// Commit after success
		if err := c.r.CommitMessages(ctx, msg); err != nil {
			return err
		}
	}
}
