// ===== Redis integration testing via 'miniredis' =====

package orderprocessor

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestIdemptonecyStore(t *testing.T) {
	// in-memory redis
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis: %v", err)
	}
	defer s.Close()

	rdb := redis.NewClient(&redis.Options{Addr: s.Addr()})
	defer rdb.Close()

	store := NewIdempotencyStore(rdb, "test:processed", 10*time.Second)

	ctx := context.Background()

	ok, err := store.IsProcessed(ctx, "evt1")
	if err != nil {
		t.Fatalf("IsProcessed: %v", err)
	}
	if ok {
		t.Fatalf("expected not processed")
	}

	if err := store.MarkProcessed(ctx, "evt1"); err != nil {
		t.Fatalf("MarkProcessed: %v", err)
	}

	ok, err = store.IsProcessed(ctx, "evt1")
	if err != nil {
		t.Fatalf("IsProcessed: %v", err)
	}
	if !ok {
		t.Fatalf("expected processed")
	}

}
