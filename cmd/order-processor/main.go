package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"order_pipeline/config"
	"order_pipeline/internal/app/orderprocessor"
	"order_pipeline/internal/infra/httpserver"
	"order_pipeline/internal/infra/kafka"
	"order_pipeline/internal/observability"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadOrderProcessor()
	if err != nil {
		panic(err)
	}

	logger, err := observability.NewLogger("order-processor")
	if err != nil {
		panic(err)
	}
	defer func() { _ = logger.Sync() }()

	health := &httpserver.Health{}
	health.SetReady(false)

	hs := httpserver.NewServer(httpserver.Config{Addr: cfg.HealthAddr}, health)
	go func() {
		logger.Info("health server listening", zap.String("addr", cfg.HealthAddr))
		if err := hs.ListenAndServe(); err != nil {
			logger.Fatal("health server failed", zap.Error(err))
		}
	}()

	cons, _ := kafka.NewConsumer(kafka.ConsumerConfig{
		Brokers:     cfg.KafkaBrokers,
		Topic:       cfg.Topic,
		GroupID:     cfg.GroupID,
		DLQTopic:    cfg.DLQTopic,
		MaxAttempts: cfg.MaxAttempts,
	}, logger)
	defer func() { _ = cons.Close() }()

	handler := orderprocessor.NewHandler(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	health.SetReady(true)
	logger.Info("order-processor started",
		zap.String("topic", cfg.Topic),
		zap.String("group_id", cfg.GroupID),
		zap.String("dlq_topic", cfg.DLQTopic),
		zap.Int("max_attempts", cfg.MaxAttempts),
	)

	runErr := make(chan error, 1)
	go func() {
		runErr <- cons.Run(ctx, handler.Handle)
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case err := <-runErr:
		if err != nil {
			logger.Error("consumer stopped with error", zap.Error(err))
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = hs.Shutdown(shutdownCtx)

	logger.Info("order-processor shutdown complete")
	time.Sleep(200 * time.Millisecond)
}
