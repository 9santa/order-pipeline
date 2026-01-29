package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Consumer struct {
	reader      *kafka.Reader
	logger      *zap.Logger
	cfg         ConsumerConfig
	dlqProducer *Producer
}

type ConsumerConfig struct {
	Brokers     []string
	Topic       string
	GroupID     string
	DLQTopic    string
	MaxAttempts int
}

func NewConsumer(cfg ConsumerConfig, logger *zap.Logger) (*Consumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Brokers,
		Topic:          cfg.Topic,
		GroupID:        cfg.GroupID,
		MinBytes:       1e3,
		MaxBytes:       10e6,
		MaxWait:        500 * time.Millisecond,
		CommitInterval: 0, // manual commit
	})

	var dlq *Producer
	if cfg.DLQTopic != "" {
		dlq, _ = NewProducer(&ProducerConfig{
			Brokers: cfg.Brokers,
			Topic:   cfg.DLQTopic,
		}, logger)
	}

	return &Consumer{
		reader:      reader,
		logger:      logger.Named("kafka-consumer"),
		cfg:         cfg,
		dlqProducer: dlq,
	}, nil
}

func (c *Consumer) Close() error {
	if c.dlqProducer != nil {
		_ = c.dlqProducer.Close()
	}
	return c.reader.Close()
}

// Handler returns nil on success. Non-nil triggers retry/DLQ logic
type Handler func(ctx context.Context, msg kafka.Message) error

func (c *Consumer) Run(ctx context.Context, h Handler) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("fetch message: %w", err)
		}

		retryCount, _ := HeaderGetInt(msg.Headers, "retry-count")

		// Execute handler
		err = h(ctx, msg)
		if err == nil {
			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				return fmt.Errorf("commit: %w", err)
			}
			continue
		}

		c.logger.Warn("handler error",
			zap.String("topic", msg.Topic),
			zap.Int("partition", msg.Partition),
			zap.Int64("offset", msg.Offset),
			zap.Int("retry_count", retryCount),
			zap.Error(err),
		)

		// Retry by re-publishing
		if retryCount+1 < c.cfg.MaxAttempts {
			// Update header retry count
			newHeaders := HeaderSet(msg.Headers, "retry-count", []byte(fmt.Sprintf("%d", retryCount+1)))
			if c.dlqProducer == nil {
				// Use a temporary producer to same topic if DLQ producer is nil
				tmp, _ := NewProducer(&ProducerConfig{
					Brokers: c.cfg.Brokers,
					Topic:   c.cfg.Topic,
				}, c.logger)
				_ = tmp.Publish(ctx, string(msg.Key), msg.Value, newHeaders)
				_ = tmp.Close()
			} else {
				// TODO: dedicated retry topic
				tmp, _ := NewProducer(&ProducerConfig{
					Brokers: c.cfg.Brokers,
					Topic:   c.cfg.Topic,
				}, c.logger)
				_ = tmp.Publish(ctx, string(msg.Key), msg.Value, newHeaders)
				_ = tmp.Close()
			}

			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				return fmt.Errorf("commit after retry publish: %w", err)
			}
			continue
		}

		// DLQ after max attempts
		if c.dlqProducer != nil {
			dlqHeaders := HeaderSet(msg.Headers, "dlq-reason", []byte(err.Error()))
			_ = c.dlqProducer.Publish(ctx, string(msg.Key), msg.Value, dlqHeaders)
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			return fmt.Errorf("commit after dlq: %w", err)
		}
	}
}
