package domain

import (
	"github.com/shopspring/decimal"
	"time"
)

type OrderStatus string

const (
	OrderStatusCreated    OrderStatus = "created"
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusPaid       OrderStatus = "paid"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusFailed     OrderStatus = "failed"
)

type Order struct {
	ID              string          `db:"id"`
	CustomerID      string          `db:"customer_id"`
	Status          OrderStatus     `db:"status"`
	TotalAmount     decimal.Decimal `db:"total_amount"`
	Currency        string          `db:"currency"`
	ShippingAddress Address         `db:"shipping_address"`
	Items           []OrderItem     `db:"-"`
	CreatedAt       time.Time       `db:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at"`
	Version         int             `db:"version"`
}

type OrderItem struct {
	ID        string          `db:"id"`
	OrderID   string          `db:"order_id"`
	ProductID string          `db:"product_id"`
	Quantity  int             `db:"quantitiy"`
	UnitPrice decimal.Decimal `db:"unit_price"`
	Total     decimal.Decimal `db:"total"`
}

type Address struct {
	Street   string `db:"street"`
	City     string `db:"city"`
	Postcode string `db:"postcode"`
	Country  string `db:"country"`
}
