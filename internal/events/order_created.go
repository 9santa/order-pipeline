package events

import (
	"time"

	"github.com/shopspring/decimal"
)

// Payload for Order Created
type OrderCreatedEvent struct {
	OrderID     string          `json:"order_id"`
	CustomerID  string          `json:"customer_id"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	Currency    string          `json:"currency"`
	CreatedAt   time.Time       `json:"created_at"`
}
