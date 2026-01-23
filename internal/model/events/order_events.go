package events

import (
	"github.com/shopspring/decimal"
	"time"
)

type OrderCreatedEvent struct {
	OrderID         string      `json:"order_id"`
	CustomerID      string      `json:"customer_id"`
	TotalAmount     float64     `json:"total_amount"`
	Currency        string      `json:"currency"`
	Items           []OrderItem `json:"items"`
	ShippingAddress Address     `json:"shipping_address"`
	CreatedAt       time.Time   `json:"created_at"`
}

type OrderItem struct {
	ProductID string          `json:"product_id"`
	Name      string          `json:"name"`
	Quantity  int             `json:"quantity"`
	UnitPrice decimal.Decimal `json:"unit_price"`
}

type Address struct {
	Street   string `json:"street"`
	City     string `json:"city"`
	Postcode string `json:"postcode"`
	Country  string `json:"country"`
}

type OrderPaidEvent struct {
	OrderID       string          `json:"order_id"`
	PaymentID     string          `json:"payment_id"`
	Amount        decimal.Decimal `json:"amount"`
	PaidAt        time.Time       `json:"paid_at"`
	TransactionID string          `json:"transaction_id"`
}
