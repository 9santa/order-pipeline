package events

import (
	"encoding/json"
	"time"
)

type EventEnvelope struct {
	ID            string          `json:"id"`
	Type          EventType       `json:"type"`
	AggregateID   string          `json:"aggregate_id"` // e.g. "event:123"
	Timestamp     time.Time       `json:"timestamp"`
	Version       int             `json:"version"`
	Metadata      Metadata        `json:"metadata"`
	Data          json.RawMessage `json:"data"`
	CorrelationID string          `json:"correlation_id"`
	TraceID       string          `json:"trace_id"`
}

type Metadata struct {
	Source      string    `json:"source"`
	RetryCount  int       `json:"retry_count"`
	PublishedAt time.Time `json:"published_at"`
	// TODO: might want to add more metadata
}
