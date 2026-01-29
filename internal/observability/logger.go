// zap logger factory

package observability

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

func NewLogger(service string) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if v := strings.TrimSpace(os.Getenv("LOG_LEVEL")); v != "" {
		_ = level.Set(v) // if no env variable, keep Info level
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(level)
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return l.With(zap.String("service", service)), nil
}
