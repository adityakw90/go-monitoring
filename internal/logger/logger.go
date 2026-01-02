package logger

import (
	"fmt"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	logger *zap.Logger
	level  *zap.AtomicLevel
}

// SetLogLevel dynamically changes the log level at runtime.
// This allows adjusting log verbosity without restarting the application.
//
// Parameters:
//   - level: The new log level ("debug", "info", "warn", "error", "fatal")
//
// Returns an error if the log level is invalid (defaults to INFO in that case).
//
// Example:
//
//	if err := logger.SetLogLevel("debug"); err != nil {
//	    log.Printf("Failed to set log level: %v", err)
//	}
func (l *logger) SetLogLevel(level string) error {
	logLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		l.Info(fmt.Sprintf("Invalid log level: %s, defaulting to INFO", level), nil)
		logLevel = zapcore.InfoLevel
	}
	l.level.SetLevel(logLevel)
	return nil
}

// Debug logs a debug-level message with optional structured fields.
// Debug logs are typically used for detailed diagnostic information.
//
// Parameters:
//   - message: The log message
//   - fields: Optional key-value pairs for structured logging (can be nil)
//
// Example:
//
//	logger.Debug("Processing request", map[string]interface{}{
//	    "request_id": "123",
//	    "user_id":    456,
//	})
func (l *logger) Debug(message string, fields map[string]interface{}) {
	zapFields := convertFields(fields)
	l.logger.Debug(message, zapFields...)
}

// Info logs an informational message with optional structured fields.
// Info logs are used for general operational information.
//
// Parameters:
//   - message: The log message
//   - fields: Optional key-value pairs for structured logging (can be nil)
//
// Example:
//
//	logger.Info("Request completed", map[string]interface{}{
//	    "status_code": 200,
//	    "duration_ms": 150,
//	})
func (l *logger) Info(message string, fields map[string]interface{}) {
	zapFields := convertFields(fields)
	l.logger.Info(message, zapFields...)
}

// Warn logs a warning message with optional structured fields.
// Warning logs indicate potentially harmful situations that don't stop execution.
//
// Parameters:
//   - message: The log message
//   - fields: Optional key-value pairs for structured logging (can be nil)
//
// Example:
//
//	logger.Warn("Rate limit approaching", map[string]interface{}{
//	    "current_rate": 90,
//	    "limit":        100,
//	})
func (l *logger) Warn(message string, fields map[string]interface{}) {
	zapFields := convertFields(fields)
	l.logger.Warn(message, zapFields...)
}

// Error logs an error message with optional structured fields.
// Error logs indicate error events that might still allow the application to continue.
//
// Parameters:
//   - message: The log message
//   - fields: Optional key-value pairs for structured logging (can be nil)
//
// Example:
//
//	logger.Error("Failed to process payment", map[string]interface{}{
//	    "payment_id": "pay_123",
//	    "error":      err.Error(),
//	})
func (l *logger) Error(message string, fields map[string]interface{}) {
	zapFields := convertFields(fields)
	l.logger.Error(message, zapFields...)
}

// Fatal logs a fatal message and exits the application.
// Fatal logs indicate severe errors that cause the application to abort.
// This function calls os.Exit(1) after logging.
//
// Parameters:
//   - message: The log message
//   - fields: Optional key-value pairs for structured logging (can be nil)
//
// Example:
//
//	logger.Fatal("Failed to initialize database", map[string]interface{}{
//	    "error": err.Error(),
//	})
//	// Application exits here
func (l *logger) Fatal(message string, fields map[string]interface{}) {
	zapFields := convertFields(fields)
	l.logger.Fatal(message, zapFields...)
}

// WithSpanContext creates a new logger instance with trace and span IDs added to all log entries.
// This enables correlation between logs and traces in distributed systems.
//
// Parameters:
//   - span: The span context containing trace and span IDs
//
// Returns:
//   - A new Logger instance with trace context embedded
//
// Example:
//
//	ctx, span := tracer.StartSpan(ctx, "operation")
//	defer tracer.EndSpan(span)
//
//	logger := logger.WithSpanContext(span.SpanContext())
//	logger.Info("Operation started", nil)
//	// Logs will include traceID and spanID fields
func (l *logger) WithSpanContext(span trace.SpanContext) Logger {
	return &logger{
		logger: l.logger.With(
			zap.String("traceID", span.TraceID().String()),
			zap.String("spanID", span.SpanID().String()),
		),
		level: l.level,
	}
}

// Sync flushes any buffered log entries.
// This should be called before application shutdown to ensure all logs are written.
// It is safe to call on a nil logger.
//
// Returns an error if flushing fails.
//
// Example:
//
//	defer func() {
//	    if err := logger.Sync(); err != nil {
//	        log.Printf("Failed to sync logger: %v", err)
//	    }
//	}()
func (l *logger) Sync() error {
	if l == nil || l.logger == nil {
		return nil
	}
	return l.logger.Sync()
}
