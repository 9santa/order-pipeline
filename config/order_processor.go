package config

type OrderProcessorConfig struct {
	HealthAddr string

	KafkaBrokers []string
	Topic        string
	RetryTopic   string
	DLQTopic     string
	GroupID      string
	MaxAttempts  int

	RedisAddr   string
	PostgresDSN string
}

func LoadOrderProcessor() (OrderProcessorConfig, error) {
	maxAttempts, err := EnvInt("MAX_ATTEMPTS", 3)
	if err != nil {
		return OrderProcessorConfig{}, err
	}

	topic := EnvString("TOPIC_ORDERS_CREATED", "orders.created")
	return OrderProcessorConfig{
		HealthAddr: EnvString("HEALTH_ADDR", ":8081"),

		KafkaBrokers: EnvStringsCSV("KAFKA_BROKERS", "localhost:9092"),
		Topic:        topic,
		RetryTopic:   EnvString("TOPIC_ORDERS_CREATED_RETRY", topic+".retry"),
		DLQTopic:     EnvString("TOPIC_ORDERS_CREATED_DLQ", topic+".dlq"),
		GroupID:      EnvString("KAFKA_GROUP_ID", "order-processor-v1"),
		MaxAttempts:  maxAttempts,

		RedisAddr:   EnvString("REDIS_ADDR", "localhost:6379"),
		PostgresDSN: EnvString("PG_DSN", "postgres://postgres:postgres@localhost:5432/order_pipeline?sslmode=disable"),
	}, nil
}
