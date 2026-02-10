package config

type OrderAPIConfig struct {
	HTTPAddr           string
	PostgresDSN        string
	OrdersCreatedTopic string
	Source             string
}

func LoadOrderAPI() OrderAPIConfig {
	return OrderAPIConfig{
		HTTPAddr:           EnvString("HTTP_ADDR", ":8080"),
		PostgresDSN:        EnvString("PS_DSN", "postgres://postgres:postgres@localhost:5432/order_pipeline?sslmode=disable"),
		OrdersCreatedTopic: EnvString("TOPIC_ORDERS_CREATED", "orders.created"),
		Source:             EnvString("SERVICE_NAME", "order-api"),
	}
}
