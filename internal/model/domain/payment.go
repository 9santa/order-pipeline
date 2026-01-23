package domain

import (
	"github.com/shopspring/decimal"
	"time"
)

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusCaptured PaymentStatus = "captured"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

type Payment struct {
	ID              string          `db:"id"`
	OrderID         string          `db:"order_id"`
	CustomerID      string          `db:"customer_id"`
	Amount          decimal.Decimal `db:"amount"`
	Currency        string          `db:"currency"`
	Status          PaymentStatus   `db:"status"`
	PaymentMethod   PaymentMethod   `db:"payment_method"`
	TransactionID   string          `db:"transaction_id"`
	GatewayResponse string          `db:"gateway:response"`
	CreatedAt       time.Time       `db:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at"`
}

type PaymentMethod struct {
	Type     string `db:"type"`
	LastFour string `db:"last_four"`
	Expiry   string `db:"expiry"`
	Token    string `db:"token"`
}
