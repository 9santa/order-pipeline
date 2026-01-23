package events

type EventType string

const (
	OrderCreated   EventType = "order.created"
	OrderPaid      EventType = "order.paid"
	PaymentFailed  EventType = "payment.failed"
	OrderShipped   EventType = "order.shipped"
	OrderDelivered EventType = "order.delivered"
	OrderCancelled EventType = "order.cancelled"
)
