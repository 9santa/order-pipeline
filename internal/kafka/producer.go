package kafka

import (
	"context"
	"encoding/json"
	"time"

	"order_pipeline/internal/model"

	"github.com/segmentio/kafka-go"
)

// Producer struct
type Producer struct {
	w *kafka.Writer
}

// Returns new Producer
func NewProducer(brokers []string, topic string) *Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.Hash{}, // for partitions within a topic: key -> partition (for correct order)
		RequiredAcks: kafka.RequireAll,
		BatchTimeout: 10 * time.Millisecond,
	}
	return &Producer{w: w}
}

// Closes the producer
func (p *Producer) Close() error { return p.w.Close() }

// Publish a Message to kafka from a Producer p
func (p *Producer) Publish(ctx context.Context, key string, ev model.Event) error {
	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: b,
	}

	return p.w.WriteMessages(ctx, msg)
}
