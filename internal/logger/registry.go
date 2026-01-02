package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(opts ...Option) (Logger, error) {
	options := &Options{
		Level: "info",
	}

	for _, opt := range opts {
		opt(options)
	}

	atomicLevel := zap.NewAtomicLevel()

	// Parse log level
	logLevel, err := zapcore.ParseLevel(options.Level)
	if err != nil {
		return nil, ErrInvalidLogLevel
	}
	atomicLevel.SetLevel(logLevel)

	config := zap.NewProductionConfig()
	config.Level = atomicLevel
	config.Encoding = "json"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	loggerInstance, err := config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return &logger{
		logger: loggerInstance,
		level:  &atomicLevel,
	}, nil
}
