package kafka

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type RetryStrategy string

const (
	RetryStrategyFixed       RetryStrategy = "fixed"
	RetryStrategyExponential RetryStrategy = "exponential"
	RetryStrategyBackoff     RetryStrategy = "backoff"
)

type RetryConfig struct {
	MaxAttempts int
	Backoff     time.Duration
	MaxBackoff  time.Duration
	Strategy    RetryStrategy
	Jitter      bool
}

type RetryManager struct {
	config RetryConfig
	logger *zap.Logger
}

func NewRetryManager(cfg RetryConfig, logger *zap.Logger) *RetryManager {
	return &RetryManager{
		config: cfg,
		logger: logger.Named("retry-manager"),
	}
}

func (rm *RetryManager) ExecuteWithRetry(ctx context.Context, msg *kafka.Message,
	fn func(ctx context.Context) error) error {
	var lastErr error

	for attempt := 1; attempt <= rm.config.MaxAttempts; attempt++ {
		err := fn(ctx)
		if err == nil {
			return nil
		}

		lastErr = err

		// Don't retry if context is cancelled
		if ctx.Err() != nil {
			return fmt.Errorf("context canceled: %w", ctx.Err())
		}

		// Log retry attempt
		rm.logger.Warn("message processing failed, retrying",
			zap.String("key", string(msg.Key)),
			zap.Int("attempt", attempt),
			zap.Int("max_attempts", rm.config.MaxAttempts),
			zap.Error(err),
		)

		// Calculate backoff
		backoff := rm.calculateBackoff(attempt)

		// Wait on this thread before retry
		select {
		case <-time.After(backoff):
			continue
		case <-ctx.Done():
			return fmt.Errorf("retry cnacelled: %w", ctx.Err())
		}
	}

	return fmt.Errorf("max retry attempts (%d) exceeded: %w", rm.config.MaxAttempts, lastErr)
}

func (rm *RetryManager) calculateBackoff(attempt int) time.Duration {
	var backoff time.Duration

	switch rm.config.Strategy {
	case RetryStrategyExponential:
		backoff = rm.config.Backoff * time.Duration(1<<uint(attempt-1))
	case RetryStrategyBackoff:
		backoff = rm.config.Backoff * time.Duration(attempt)
	default: // fixed
		backoff = rm.config.Backoff
	}

	// Apply jitter to prevent multiple failed clients from bombarding the server at the exact same moment
	if rm.config.Jitter {
		jitter := time.Duration(rand.Int63n(int64(backoff / 2)))
		if rand.Intn(2) == 0 {
			backoff -= jitter
		} else {
			backoff += jitter
		}
	}

	// Cap at max backoff
	if backoff > rm.config.MaxBackoff {
		backoff = rm.config.MaxBackoff
	}

	return backoff
}
