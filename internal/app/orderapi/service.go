package orderapi

import (
	"context"
	"fmt"
	"sync"
	"time"

	"order_pipeline/internal/events"
	"order_pipeline/internal/infra/kafka"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/shopspring/decimal"
)

// Service will push orders to kafka Producer
type Service struct {
	source string
	topic  string
	prod   *kafka.Producer
	mu     sync.Mutex
	orders map[string]Order // in memory for now
}

type Order struct {
	ID          string          `json:"id"`
	CustomerID  string          `json:"customer_id"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	Currency    string          `json:"currency"`
	CreatedAt   time.Time       `json:"created_at"`
}

func NewService(source, topic string, prod *kafka.Producer) *Service {
	return &Service{
		source: source,
		topic:  topic,
		prod:   prod,
		orders: make(map[string]Order),
	}
}

func (s *Service) CreateOrder(ctx context.Context, customerID string, totalAmount decimal.Decimal, currency string) (Order, string, error) {
	now := time.Now().UTC()
	orderID := "ord_" + events.NewID()

	// Construct the order
	order := Order{
		ID:          orderID,
		CustomerID:  customerID,
		TotalAmount: totalAmount,
		Currency:    currency,
		CreatedAt:   now,
	}

	// Lock orders storage before writing
	s.mu.Lock()
	s.orders[orderID] = order
	// Unlock
	s.mu.Unlock()

	// Construct order event
	ev := events.OrderCreatedEvent{
		OrderID:     orderID,
		CustomerID:  customerID,
		TotalAmount: totalAmount,
		Currency:    currency,
		CreatedAt:   now,
	}

	// Wrap the order event into Envelope
	env, err := events.NewEnvelope(events.OrderCreated, orderID, 1, s.source, ev, "")
	if err != nil {
		return Order{}, "", err
	}

	payload, err := events.EncodeEnvelope(env)
	if err != nil {
		return Order{}, "", err
	}

	headers := []kafkago.Header{
		{Key: "event-type", Value: []byte(env.Type)},
		{Key: "event-id", Value: []byte(env.ID)},
	}

	if err := s.prod.Publish(ctx, orderID, payload, headers); err != nil {
		return Order{}, "", fmt.Errorf("publish order.created: %w", err)
	}

	return order, env.ID, nil
}

func (s *Service) GetOrder(_ context.Context, id string) (Order, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o, ok := s.orders[id]
	return o, ok
}
