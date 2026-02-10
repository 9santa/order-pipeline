package config

import "time"

type OutboxPublisherConfig struct {
	KafkaBrokers []string
	PostgresDSN  string

	PollInterval time.Duration
	BatchSize    int
}

func LoadOutboxPublisher() (OutboxPublisherConfig, error) {
	poll, err := EnvDuration("POLL_INTERVAL", 1*time.Second)
	if err != nil {
		return OutboxPublisherConfig{}, err
	}
	batch, err := EnvInt("BATCH_SIZE", 50)
	if err != nil {
		return OutboxPublisherConfig{}, err
	}

	return OutboxPublisherConfig{
		KafkaBrokers: EnvStringsCSV("KAFKA_BROKERS", "localhost:9092"),
		PostgresDSN:  EnvString("PG_DSN", "postgres://postgres:postgres@localhost:5432/order_pipeline?sslmode=disable"),
		PollInterval: poll,
		BatchSize:    batch,
	}, nil
}
