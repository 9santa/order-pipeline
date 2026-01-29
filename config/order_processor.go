package config

type OrderProcessorConfig struct {
	HealthAddr   string
	KafkaBrokers []string
	Topic        string
	GroupID      string
	DLQTopic     string
	MaxAttempts  int
}

func LoadOrderProcessor() (OrderProcessorConfig, error) {
	maxAttempts, err := EnvInt("MAX_ATTEMPTS", 3)
	if err != nil {
		return OrderProcessorConfig{}, err
	}

	topic := EnvString("TOPIC_ORDERS_CREATED", "orders.created")
	return OrderProcessorConfig{
		HealthAddr:   EnvString("HEALTH_ADDR", ":8081"),
		KafkaBrokers: EnvStringsCSV("KAFKA_BROKERS", "localhost:9092"),
		Topic:        topic,
		GroupID:      EnvString("KAFKA_GROUP_ID", "order-processor-v1"),
		DLQTopic:     EnvString("DLQ_TOPIC", topic+".dlq"),
		MaxAttempts:  maxAttempts,
	}, nil
}
