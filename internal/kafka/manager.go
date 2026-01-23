package kafka

import (
	"order_pipeline/internal/events"
	"time"
)

type TopicConfig struct {
	Name              string
	Partitions        int
	ReplicationFactor int
	Retention         time.Duration
	Compression       string
}

var Topics = map[events.EventType]TopicConfig{
	events.OrderCreated: {
		Name:              "orders.created",
		Partitions:        6,
		ReplicationFactor: 3,
		Retention:         7 * 24 * time.Hour, // 1 Week
		Compression:       "snappe",
	},
	events.OrderPaid: {
		Name:              "orders.paid",
		Partitions:        4,
		ReplicationFactor: 3,
	},
	// TODO: other topics
}

type ConsumerGroup struct {
	Name        string
	RetryPolicy RetryPolicy
	DLQ         bool
}

var ConsumerGroups = map[string]ConsumerGroup{
	"order-processor": {
		Name: "order-processor-v1",
		RetryPolicy: RetryPolicy{
			MaxAttempts: 3,
			Backoff:     2 * time.Second,
		},
		DLQ: true,
	},
}

type RetryPolicy struct {
	MaxAttempts int
	Backoff     time.Duration
}
