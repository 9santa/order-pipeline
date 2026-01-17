package producer

import (
	"context"
	"log"
	"order_pipeline/internal/kafka"
	"order_pipeline/internal/model"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	brokers := strings.Split(env("KAFKA_BROKERS", "localhost:9092"), ",")
	topic := env("KAFKA_TOPIC", "orders.created")

	p := kafka.NewProducer(brokers, topic)
	defer p.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	i := 1
	for {
		select {
		case <-ctx.Done():
			log.Println("producer: shutdown")
			return
		case <-ticker.C:
			orderID := "order-" + itoa(i)
			ev := model.Event{
				EventID:   "evt-" + itoa(i),
				Type:      "order_created",
				OrderID:   orderID,
				CreatedAt: time.Now().UTC(),
				Payload: map[string]any{
					"amount": 100 + i,
				},
			}

			if err := p.Publish(ctx, orderID, ev); err != nil {
				log.Printf("publish error: %v\n", err)
				continue
			}
			log.Printf("pubslihed: %s\n", orderID)
			i++
		}
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	sign := ""
	if n < 0 {
		sign = "-"
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return sign + string(b[i:])
}
