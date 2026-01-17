package consumer

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
	groupID := env("KAFKA_GROUP_ID", "orders-service")

	c := kafka.NewConsumer(brokers, topic, groupID)
	defer c.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Printf("consumer started: topic=%s group=%s\n", topic, groupID)

	handler := func(ctx context.Context, ev model.Event) error {
		// for now imitation of work
		time.Sleep(100 * time.Millisecond)
		log.Printf("got: type=%s order=%s event=%s\n", ev.Type, ev.OrderID, ev.EventID)
		return nil
	}

	if err := c.Run(ctx, handler); err != nil {
		log.Fatalf("consumer error: %v", err)
	}

}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
