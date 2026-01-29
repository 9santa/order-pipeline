package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// Producer struct
type Producer struct {
	writer *kafka.Writer
	logger *zap.Logger
	topic  string
}

type ProducerConfig struct {
	Brokers      []string
	Topic        string
	Balancer     kafka.Hash
	RequiredAcks kafka.RequiredAcks
	BatchSize    int
	BatchTimeout time.Duration
	MaxAttempts  int
	WriteTimeout time.Duration
	Async        bool
}

// Returns new Producer
func NewProducer(cfg *ProducerConfig, logger *zap.Logger) (*Producer, error) {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        cfg.Topic,
		Balancer:     &cfg.Balancer,
		RequiredAcks: cfg.RequiredAcks,
		BatchSize:    cfg.BatchSize,
		BatchTimeout: cfg.BatchTimeout,
		MaxAttempts:  cfg.MaxAttempts,
		WriteTimeout: cfg.WriteTimeout,
		Async:        cfg.Async,
		Transport:    &kafka.Transport{},
	}

	return &Producer{
		writer: writer,
		logger: logger.Named("kafka-producer"),
		topic:  cfg.Topic,
	}, nil
}

// Closes the producer
func (p *Producer) Close() error {
	p.logger.Info("closing kafka producer")
	return p.writer.Close()
}

// Publish a Message to kafka from a Producer p
func (p *Producer) Publish(ctx context.Context, key string, payload []byte, headers []kafka.Header) error {
	start := time.Now().UTC()

	msg := kafka.Message{
		Key:     []byte(key),
		Value:   payload,
		Time:    time.Now().UTC(),
		Headers: headers,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("kafka write: %w", err)
	}

	// Logging
	duration := time.Since(start)

	p.logger.Debug("message published",
		zap.String("key", key),
		zap.Duration("duration", duration),
	)

	return nil
}
