package postgres

import (
	"context"
	"database/sql"
	"time"
)

type OutboxRow struct {
	ID        string
	Topic     string
	Key       string
	Payload   []byte
	CreatedAt time.Time
}

type OutboxRepo struct {
	db *sql.DB
}

func NewOutboxRepo(db *sql.DB) *OutboxRepo {
	return &OutboxRepo{db: db}
}

func (r *OutboxRepo) InsertTx(ctx context.Context, tx *sql.Tx, id, topic, key string, payload []byte, createdAt time.Time) error {
	_, err := tx.ExecContext(ctx, `
	INSERT INTO outbox_events (id, topic, key, payload, created_at)
	VALUES ($1,$2,$3,$4,$5)`, id, topic, key, payload, createdAt)
	return err
}

// Locks rows, skips already locked rows, ensuring efficient work distribution and high concurrency
func (r *OutboxRepo) FetchUnsentTx(ctx context.Context, tx *sql.Tx, limit int) ([]OutboxRow, error) {
	rows, err := tx.QueryContext(ctx, `
	SELECT id, topic, key, payload, created_at
	FROM outbox_events
	WHERE sent_at IS NULL
	ORDER BY created_at
	FOR UPDATE SKIP LOCKED
	LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Append to output
	var out []OutboxRow
	for rows.Next() {
		var rrow OutboxRow
		if err := rows.Scan(&rrow.ID, &rrow.Topic, &rrow.Key, &rrow.Payload, &rrow.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, rrow)
	}
	return out, rows.Err()
}

func (r *OutboxRepo) MarkSentTx(ctx context.Context, tx *sql.Tx, id string) error {
	_, err := tx.ExecContext(ctx, `UPDATE outbox_events SET sent_at = NOW() WHERE id=$1`, id)
	return err
}
