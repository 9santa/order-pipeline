package config

type OrderAPIConfig struct {
	HTTPAddr           string
	KafkaBrokers       []string
	OrdersCreatedTopic string
	Source             string
}

func LoadOrderAPI() OrderAPIConfig {
	return OrderAPIConfig{
		HTTPAddr:           EnvString("HTTP_ADDR", ":8080"),
		KafkaBrokers:       EnvStringsCSV("KAFKA_BROKERS", "localhost:9092"),
		OrdersCreatedTopic: EnvString("TOPIC_ORDERS_CREATED", "orders.created"),
		Source:             EnvString("SERVICE_NAME", "order-api"),
	}
}
