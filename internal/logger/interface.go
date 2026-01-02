package logger

import "go.opentelemetry.io/otel/trace"

// Logger defines the contract for logging operations.
type Logger interface {
	SetLogLevel(level string) error
	Debug(message string, fields map[string]interface{})
	Info(message string, fields map[string]interface{})
	Warn(message string, fields map[string]interface{})
	Error(message string, fields map[string]interface{})
	Fatal(message string, fields map[string]interface{})
	WithSpanContext(span trace.SpanContext) Logger
	Sync() error
}
