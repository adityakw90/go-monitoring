package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates and configures a zap-backed Logger according to the provided options.
// It defaults the log level to "info", parses and applies the configured level (returning ErrInvalidLogLevel on parse failure),
// enforces JSON encoding and a fixed timestamp layout ("2006-01-02T15:04:05.000-0700"), and optionally directs output to a custom path.
// The built logger includes caller information and a caller-skip of 1; on build failure it returns a wrapped error.
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
	// Use TimeEncoderOfLayout to ensure consistent format with +0000 for UTC instead of Z
	// This ensures timestamps are always in offset format (e.g., +0000, +0700) regardless of timezone
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000-0700")

	if options.OutputPath != "" {
		config.OutputPaths = []string{options.OutputPath}
	}

	loggerInstance, err := config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return &logger{
		logger: loggerInstance,
		level:  &atomicLevel,
	}, nil
}