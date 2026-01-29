package main

import (
	"context"
	"net"
	"net/http"
	"order_pipeline/config"
	"order_pipeline/internal/app/orderapi"
	"order_pipeline/internal/infra/kafka"
	"order_pipeline/internal/observability"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadOrderAPI()

	logger, err := observability.NewLogger(cfg.Source)
	if err != nil {
		panic(err)
	}
	defer func() { _ = logger.Sync() }()

	prod, _ := kafka.NewProducer(&kafka.ProducerConfig{
		Brokers: cfg.KafkaBrokers,
		Topic:   cfg.OrdersCreatedTopic,
	}, logger)
	defer func() { _ = prod.Close() }()

	svc := orderapi.NewService(cfg.Source, cfg.OrdersCreatedTopic, prod)
	handlers := orderapi.NewHandlers(svc)

	mux := http.NewServeMux()
	handlers.Register(mux)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	})

	serv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("order-api listening", zap.String("addr", cfg.HTTPAddr))
		ln, err := net.Listen("tcp", cfg.HTTPAddr)
		if err != nil {
			logger.Fatal("listen failed", zap.Error(err))
		}
		if err := serv.Serve(ln); err != nil && err != http.ErrServerClosed {
			logger.Fatal("http serve failed", zap.Error(err))
		}
	}()

	<-ctx.Done()

	shutDownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = serv.Shutdown(shutDownCtx)

	logger.Info("order-api shutdown complete")
}
