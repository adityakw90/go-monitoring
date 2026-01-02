package logger

import (
	"context"
	"testing"
	"time"

	"github.com/adityakw90/go-monitoring/internal/tracer"
)

func TestLogger_Logger_SetLogLevel(t *testing.T) {
	loggerInstance, err := NewLogger()
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	tests := []struct {
		name    string
		level   string
		wantErr bool
	}{
		{
			name:    "valid level",
			level:   "debug",
			wantErr: false,
		},
		{
			name:    "invalid level",
			level:   "invalid",
			wantErr: false, // SetLogLevel doesn't return error, just defaults
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := loggerInstance.SetLogLevel(tt.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetLogLevel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLogger_Logger_LogMethods(t *testing.T) {
	loggerInstance, err := NewLogger()
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	fields := map[string]interface{}{
		"key": "value",
	}

	// Test that methods don't panic
	loggerInstance.Debug("debug message", fields)
	loggerInstance.Info("info message", fields)
	loggerInstance.Warn("warn message", fields)
	loggerInstance.Error("error message", fields)
	loggerInstance.Info("message without fields", nil)
}

func TestLogger_Logger_WithSpanContext(t *testing.T) {
	loggerInstance, err := NewLogger()
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	// Create a mock span context
	ctx := context.Background()
	tracer, err := tracer.NewTracer(tracer.WithServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewTracer() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = tracer.Shutdown(ctx)
	}()

	_, span := tracer.StartSpan(ctx, "test-operation")
	defer span.End()

	// Test WithSpanContext
	loggerWithSpan := loggerInstance.WithSpanContext(span.SpanContext())
	if loggerWithSpan == nil {
		t.Errorf("WithSpanContext() returned nil")
	}

	// Test that the logger with span context can log
	loggerWithSpan.Info("message with span context", map[string]interface{}{
		"test": "value",
	})

	// Verify it's a different instance
	if loggerInstance == loggerWithSpan {
		t.Errorf("WithSpanContext() returned same logger instance")
	}
}

func TestLogger_Logger_AllLogLevels(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error"}

	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			loggerInstance, err := NewLogger(WithLevel(level))
			if err != nil {
				t.Fatalf("NewLogger() with level %s error = %v", level, err)
			}

			fields := map[string]interface{}{
				"level": level,
			}

			// Test all log methods
			loggerInstance.Debug("debug message", fields)
			loggerInstance.Info("info message", fields)
			loggerInstance.Warn("warn message", fields)
			loggerInstance.Error("error message", fields)
		})
	}
}
