package config

type RetryDispatcherConfig struct {
	KafkaBrokers []string
	RetryTopic   string
	GroupID      string
}

func LoadRetryDispatcher() RetryDispatcherConfig {
	return RetryDispatcherConfig{
		KafkaBrokers: EnvStringsCSV("KAFKA_BROKERS", "localhost:9092"),
		RetryTopic:   EnvString("TOPIC_ORDERS_CREATED_RETRY", "orders.created.retry"),
		GroupID:      EnvString("KAFKA_GROUP_ID", "retry-dispatcher-v1"),
	}
}
