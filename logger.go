package monitoring

import (
	"fmt"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap logger with OpenTelemetry integration.
type Logger struct {
	logger *zap.Logger
	level  *zap.AtomicLevel
}

// LoggerOptions contains logger configuration.
type LoggerOptions struct {
	Level string
}

// LoggerOption configures LoggerOptions.
type LoggerOption func(*LoggerOptions)

// withLoggerLevel sets the log level (internal use).
func withLoggerLevel(level string) LoggerOption {
	return func(o *LoggerOptions) {
		o.Level = level
	}
}

// NewLogger initializes a new zap logger with the given options.
func NewLogger(opts ...LoggerOption) (*Logger, error) {
	options := &LoggerOptions{
		Level: "info",
	}

	for _, opt := range opts {
		opt(options)
	}

	atomicLevel := zap.NewAtomicLevel()

	// Parse log level
	logLevel, err := zapcore.ParseLevel(options.Level)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidLogLevel, options.Level)
	}
	atomicLevel.SetLevel(logLevel)

	config := zap.NewProductionConfig()
	config.Level = atomicLevel
	config.Encoding = "json"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return &Logger{
		logger: logger,
		level:  &atomicLevel,
	}, nil
}

// SetLogLevel dynamically changes the log level.
func (l *Logger) SetLogLevel(level string) error {
	logLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		l.Info(fmt.Sprintf("Invalid log level: %s, defaulting to INFO", level), nil)
		logLevel = zapcore.InfoLevel
	}
	l.level.SetLevel(logLevel)
	return nil
}

// Debug logs a debug message with optional structured fields.
func (l *Logger) Debug(message string, fields map[string]interface{}) {
	zapFields := l.convertFields(fields)
	l.logger.Debug(message, zapFields...)
}

// Info logs an informational message with optional structured fields.
func (l *Logger) Info(message string, fields map[string]interface{}) {
	zapFields := l.convertFields(fields)
	l.logger.Info(message, zapFields...)
}

// Warn logs a warning message with optional structured fields.
func (l *Logger) Warn(message string, fields map[string]interface{}) {
	zapFields := l.convertFields(fields)
	l.logger.Warn(message, zapFields...)
}

// Error logs an error message with optional structured fields.
func (l *Logger) Error(message string, fields map[string]interface{}) {
	zapFields := l.convertFields(fields)
	l.logger.Error(message, zapFields...)
}

// Fatal logs a fatal message and exits.
func (l *Logger) Fatal(message string, fields map[string]interface{}) {
	zapFields := l.convertFields(fields)
	l.logger.Fatal(message, zapFields...)
}

// WithSpanContext adds trace and span IDs to the logger.
func (l *Logger) WithSpanContext(span trace.SpanContext) *Logger {
	return &Logger{
		logger: l.logger.With(
			zap.String("traceID", span.TraceID().String()),
			zap.String("spanID", span.SpanID().String()),
		),
		level: l.level,
	}
}

// convertFields converts map[string]interface{} to zap fields.
func (l *Logger) convertFields(fields map[string]interface{}) []zap.Field {
	if fields == nil {
		return nil
	}
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}
