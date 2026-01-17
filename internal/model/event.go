package model

import "time"

type Event struct {
	EventID   string    `json:"event_id"`
	Type      string    `json:"type"`
	OrderID   string    `json:"order_id"`
	CreatedAt time.Time `json:"created_at"`

	// Payload - actual content of the message
	Payload map[string]any `json:"payload"`
}
