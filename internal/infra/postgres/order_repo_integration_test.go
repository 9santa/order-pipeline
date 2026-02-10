// ===== Test OrderRepo + Migrations =====

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestOrderRepo_CRUD(t *testing.T) {
	dsn := startPostgres(t)

	db, err := Open(dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	ctx := context.Background()

	if err := Ping(ctx, db); err != nil {
		t.Fatalf("db ping: %v", err)
	}

	if err := Migrate(ctx, db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	repo := NewOrderRepo(db)

	now := time.Now().UTC()
	order := OrderRow{
		ID:          "ord_1",
		CustomerID:  "cust_1",
		TotalAmount: decimal.NewFromInt(105),
		Currency:    "RUB",
		Status:      "created",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("begin tx: %v", err)
	}
	if err := repo.InsertTx(ctx, tx, order); err != nil {
		t.Fatalf("insert: %v", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit: %v", err)
	}

	got, ok, err := repo.Get(ctx, "ord_1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if !ok {
		t.Fatalf("expected order to exist")
	}
	if got.Status != "created" || got.CustomerID != "cust_1" {
		t.Fatalf("unexpected order: %+v", got)
	}

	if err := repo.UpdateStatus(ctx, "ord_1", "paid"); err != nil {
		t.Fatalf("update status: %v", err)
	}

	got2, ok, err := repo.Get(ctx, "ord_1")
	if err != nil {
		t.Fatalf("get2: %v", err)
	}
	if !ok || got2.Status != "paid" {
		t.Fatalf("expected status=paid, got: %+v", got2)
	}
}
