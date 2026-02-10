package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

const schema = `
CREATE TABLE IF NOT EXISTS orders (
	id TEXT PRIMARY KEY,
	customer_id TEXT NOT NULL,
	total_amount DECIMAL NOT NULL,
	currency TEXT NOT NULL,
	status TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS outbox_events (
	id TEXT PRIMARY KEY,
	topic TEXT NOT NULL,
	key TEXT NOT NULL,
	payload JSONB NOT NULL,
	created_at TIMESTAMPTZ NOT NULL,
	sent_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS outbox_unsent_idx
	ON outbox_events (created_at)
	WHERE sent_at IS NULL;
`

func Migrate(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}
