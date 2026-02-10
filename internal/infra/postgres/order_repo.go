package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type OrderRow struct {
	ID          string
	CustomerID  string
	TotalAmount decimal.Decimal
	Currency    string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) InsertTx(ctx context.Context, tx *sql.Tx, o OrderRow) error {
	_, err := tx.ExecContext(ctx, `
	INSERT INTO orders (id, customer_id, total_amount, currency, status, created_at, updated_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7)`, o.ID, o.CustomerID, o.TotalAmount, o.Currency, o.Status, o.CreatedAt, o.UpdatedAt)
	return err
}

func (r *OrderRepo) Get(ctx context.Context, id string) (OrderRow, bool, error) {
	var o OrderRow
	err := r.db.QueryRowContext(ctx, `SELECT id, customer_id, total_amount, currency, status, created_at, updated_at
	FROM orders WHERE id=$1`, id).Scan(&o.ID, &o.CustomerID, &o.TotalAmount, &o.Currency, &o.Status, &o.CreatedAt, &o.UpdatedAt)
	if err == sql.ErrNoRows {
		return OrderRow{}, false, nil
	}
	if err != nil {
		return OrderRow{}, false, err
	}
	return o, true, nil
}

func (r *OrderRepo) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE orders SET status=$2, updated_at=NOW() where id=$1`, id, status)
	return err
}
