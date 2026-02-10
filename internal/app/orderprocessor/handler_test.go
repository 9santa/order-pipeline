// ===== Handler testing =====

package orderprocessor

import (
	"context"
	"testing"
	"time"

	"order_pipeline/internal/events"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	kafkago "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestHandler_IsIdempotent(t *testing.T) {
	// log capture
	core, logs := observer.New(zap.InfoLevel)
	logger := zap.New(core)

	// in-memory redis
	s, _ := miniredis.Run()
	defer s.Close()
	rdb := redis.NewClient(&redis.Options{Addr: s.Addr()})
	defer rdb.Close()

	idem := NewIdempotencyStore(rdb, "order-pipeline:processed:", 1*time.Hour)
	h := NewHandler(logger, idem)

	// test duplicates handling with constant envelope ID
	ev := events.OrderCreatedEvent{
		OrderID:    "ord_123",
		CustomerID: "cust_1",
		CreatedAt:  time.Now().UTC(),
	}
	env, _ := events.NewEnvelope(events.OrderCreated, "ord_123", 1, "test", ev, "")
	env.ID = "evt_fixed" // force constant ID for the test
	b, _ := events.EncodeEnvelope(env)

	msg := kafkago.Message{Value: b}

	// first time => processed
	if err := h.Handle(context.Background(), msg); err != nil {
		t.Fatalf("first handle: %v", err)
	}
	// second time => duplicate
	if err := h.Handle(context.Background(), msg); err != nil {
		t.Fatalf("second handle: %v", err)
	}

	entries := logs.All()
	var processed, dup int
	for _, e := range entries {
		if e.Message == "processed order.created" {
			processed++
		}
		if e.Message == "duplicate event; skipping" {
			dup++
		}
	}
	if processed != 1 || dup != 1 {
		t.Fatalf("expected processed=1 dup=1, got processed=%d dup=%d", processed, dup)
	}
}
