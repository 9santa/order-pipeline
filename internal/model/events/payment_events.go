package events

import (
	"github.com/shopspring/decimal"

	"order_pipeline/internal/model/domain"

	"time"
)

type PaymentAcceptedEvent struct {
	PaymentID     string               `json:"payment_id"`
	OrderID       string               `json:"order_id"`
	CustomerID    string               `json:"customer_id"`
	Amount        decimal.Decimal      `json:"amount"`
	Currency      string               `json:"currency"`
	AcceptedAt    time.Time            `json:"accepted_at"`
	PaymentMethod domain.PaymentMethod `json:"payment_method"`
}

type PaymentFailedEvent struct {
	PaymentID  string    `json:"payment_id"`
	OrderID    string    `json:"order_id"`
	CustomerID string    `json:"customer_id"`
	Reason     string    `json:"reason"`
	FailedAt   time.Time `json:"failed_at"`
	ErrorCode  string    `json:"error_code"`
}
