package monitoring

import (
	"context"
	"testing"
	"time"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name    string
		opts    []LoggerOption
		wantErr bool
	}{
		{
			name:    "default logger",
			opts:    nil,
			wantErr: false,
		},
		{
			name:    "with valid log level",
			opts:    []LoggerOption{withLoggerLevel("debug")},
			wantErr: false,
		},
		{
			name:    "with invalid log level",
			opts:    []LoggerOption{withLoggerLevel("invalid")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewLogger(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && logger == nil {
				t.Errorf("NewLogger() returned nil logger")
			}
		})
	}
}

func TestLogger_SetLogLevel(t *testing.T) {
	logger, err := NewLogger()
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
			err := logger.SetLogLevel(tt.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetLogLevel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLogger_LogMethods(t *testing.T) {
	logger, err := NewLogger()
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	fields := map[string]interface{}{
		"key": "value",
	}

	// Test that methods don't panic
	logger.Debug("debug message", fields)
	logger.Info("info message", fields)
	logger.Warn("warn message", fields)
	logger.Error("error message", fields)
	logger.Info("message without fields", nil)
}

func TestLogger_WithSpanContext(t *testing.T) {
	logger, err := NewLogger()
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	// Create a mock span context
	ctx := context.Background()
	tracer, err := NewTracer(withTracerServiceName("test-service"))
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
	loggerWithSpan := logger.WithSpanContext(span.SpanContext())
	if loggerWithSpan == nil {
		t.Errorf("WithSpanContext() returned nil")
	}

	// Test that the logger with span context can log
	loggerWithSpan.Info("message with span context", map[string]interface{}{
		"test": "value",
	})

	// Verify it's a different instance
	if logger == loggerWithSpan {
		t.Errorf("WithSpanContext() returned same logger instance")
	}
}

func TestLogger_AllLogLevels(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error"}

	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			logger, err := NewLogger(withLoggerLevel(level))
			if err != nil {
				t.Fatalf("NewLogger() with level %s error = %v", level, err)
			}

			fields := map[string]interface{}{
				"level": level,
			}

			// Test all log methods
			logger.Debug("debug message", fields)
			logger.Info("info message", fields)
			logger.Warn("warn message", fields)
			logger.Error("error message", fields)
		})
	}
}

func TestLogger_ConvertFields(t *testing.T) {
	logger, err := NewLogger()
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	// Test with nil fields
	logger.Info("message with nil fields", nil)

	// Test with empty fields
	logger.Info("message with empty fields", map[string]interface{}{})

	// Test with various field types
	fields := map[string]interface{}{
		"string": "value",
		"int":    42,
		"float":  3.14,
		"bool":   true,
		"nil":    nil,
		"slice":  []string{"a", "b"},
		"map":    map[string]int{"key": 1},
	}
	logger.Info("message with various field types", fields)
}
