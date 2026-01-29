// codec.go: Creates, Encodes, Decodes Envelopes

package events

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// Random ID generator for Envelopes
func NewID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

// Envelope constructor
func NewEnvelope(typ EventType, aggregateID string, version int, source string, data any, correlationID string) (EventEnvelope, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return EventEnvelope{}, fmt.Errorf("error, marshal data: %w", err)
	}

	now := time.Now().UTC()
	return EventEnvelope{
		ID:          NewID(),
		Type:        typ,
		AggregateID: aggregateID,
		Version:     version,
		Metadata: Metadata{
			Source:      source,
			RetryCount:  0,
			PublishedAt: now,
		},
		Data:          raw,
		CorrelationID: correlationID,
	}, nil
}

// TODO: Encode Envelope (right now doesn't do anything)
func EncodeEnvelope(env EventEnvelope) ([]byte, error) {
	b, err := json.Marshal(env)
	if err != nil {
		return nil, fmt.Errorf("error, marshal data: %w", err)
	}
	return b, nil
}

// TODO: Decode Envelope
func DecodeEnvelope(b []byte) (EventEnvelope, error) {
	var env EventEnvelope
	if err := json.Unmarshal(b, &env); err != nil {
		return EventEnvelope{}, fmt.Errorf("unmarshal envelope: %w", err)
	}
	return env, nil
}

func DecodeData[T any](env EventEnvelope) (T, error) {
	var out T
	if err := json.Unmarshal(env.Data, &out); err != nil {
		return out, fmt.Errorf("unmarshal data: %w", err)
	}
	return out, nil
}
