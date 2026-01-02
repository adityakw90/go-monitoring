package logger

import (
	"context"
	"testing"
	"time"

	"github.com/adityakw90/go-monitoring/internal/tracer"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap/zapcore"
)

func TestLogger_Logger_SetLogLevel(t *testing.T) {
	tests := []struct {
		name          string
		level         string
		expectedLevel zapcore.Level
	}{
		{
			name:          "valid debug level",
			level:         "debug",
			expectedLevel: zapcore.DebugLevel,
		},
		{
			name:          "valid info level",
			level:         "info",
			expectedLevel: zapcore.InfoLevel,
		},
		{
			name:          "valid warn level",
			level:         "warn",
			expectedLevel: zapcore.WarnLevel,
		},
		{
			name:          "valid error level",
			level:         "error",
			expectedLevel: zapcore.ErrorLevel,
		},
		{
			name:          "valid fatal level",
			level:         "fatal",
			expectedLevel: zapcore.FatalLevel,
		},
		{
			name:          "invalid level defaults to info",
			level:         "invalid",
			expectedLevel: zapcore.InfoLevel,
		},
		{
			name:          "empty level defaults to info",
			level:         "",
			expectedLevel: zapcore.InfoLevel,
		},
		{
			name:          "level is case insensitive",
			level:         "DEBUG",
			expectedLevel: zapcore.DebugLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loggerInstance, err := NewLogger()
			require.NoError(t, err)
			loggerInstance.SetLogLevel(tt.level)
			if loggerInstance.(*logger).level.Level() != tt.expectedLevel {
				t.Errorf("SetLogLevel() level = %v, want %v", loggerInstance.(*logger).level.Level(), tt.expectedLevel)
			}
		})
	}
}

func TestLogger_Logger_Debug(t *testing.T) {
	loggerInstance, err := NewLogger(WithLevel("debug"))
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	tests := []struct {
		name    string
		message string
		fields  map[string]interface{}
	}{
		{
			name:    "debug with fields",
			message: "debug message",
			fields: map[string]interface{}{
				"key":  "value",
				"num":  42,
				"bool": true,
			},
		},
		{
			name:    "debug without fields",
			message: "debug message",
			fields:  nil,
		},
		{
			name:    "debug with empty fields",
			message: "debug message",
			fields:  map[string]interface{}{},
		},
		{
			name:    "debug with empty message",
			message: "",
			fields: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name:    "debug with various field types",
			message: "debug message",
			fields: map[string]interface{}{
				"string": "value",
				"int":    42,
				"float":  3.14,
				"bool":   true,
				"slice":  []string{"a", "b"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loggerInstance.Debug(tt.message, tt.fields)
		})
	}
}

func TestLogger_Logger_Info(t *testing.T) {
	loggerInstance, err := NewLogger()
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	tests := []struct {
		name    string
		message string
		fields  map[string]interface{}
	}{
		{
			name:    "info with fields",
			message: "info message",
			fields: map[string]interface{}{
				"key":  "value",
				"num":  42,
				"bool": true,
			},
		},
		{
			name:    "info without fields",
			message: "info message",
			fields:  nil,
		},
		{
			name:    "info with empty fields",
			message: "info message",
			fields:  map[string]interface{}{},
		},
		{
			name:    "info with empty message",
			message: "",
			fields: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name:    "info with various field types",
			message: "info message",
			fields: map[string]interface{}{
				"string": "value",
				"int":    42,
				"float":  3.14,
				"bool":   true,
				"slice":  []string{"a", "b"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loggerInstance.Info(tt.message, tt.fields)
		})
	}
}

func TestLogger_Logger_Warn(t *testing.T) {
	loggerInstance, err := NewLogger()
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	tests := []struct {
		name    string
		message string
		fields  map[string]interface{}
	}{
		{
			name:    "warn with fields",
			message: "warn message",
			fields: map[string]interface{}{
				"key":  "value",
				"num":  42,
				"bool": true,
			},
		},
		{
			name:    "warn without fields",
			message: "warn message",
			fields:  nil,
		},
		{
			name:    "warn with empty fields",
			message: "warn message",
			fields:  map[string]interface{}{},
		},
		{
			name:    "warn with empty message",
			message: "",
			fields: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name:    "warn with various field types",
			message: "warn message",
			fields: map[string]interface{}{
				"string": "value",
				"int":    42,
				"float":  3.14,
				"bool":   true,
				"slice":  []string{"a", "b"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loggerInstance.Warn(tt.message, tt.fields)
		})
	}
}

func TestLogger_Logger_Error(t *testing.T) {
	loggerInstance, err := NewLogger()
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	tests := []struct {
		name    string
		message string
		fields  map[string]interface{}
	}{
		{
			name:    "error with fields",
			message: "error message",
			fields: map[string]interface{}{
				"key":  "value",
				"num":  42,
				"bool": true,
			},
		},
		{
			name:    "error without fields",
			message: "error message",
			fields:  nil,
		},
		{
			name:    "error with empty fields",
			message: "error message",
			fields:  map[string]interface{}{},
		},
		{
			name:    "error with empty message",
			message: "",
			fields: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name:    "error with various field types",
			message: "error message",
			fields: map[string]interface{}{
				"string": "value",
				"int":    42,
				"float":  3.14,
				"bool":   true,
				"slice":  []string{"a", "b"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loggerInstance.Error(tt.message, tt.fields)
		})
	}
}

func TestLogger_Logger_Fatal(t *testing.T) {
	// Note: Fatal calls os.Exit(1), so we can't test it normally
	// This test verifies the method exists and can be called without compilation errors
	// In a real scenario, you might use a subprocess or mock for testing
	loggerInstance, err := NewLogger()
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	tests := []struct {
		name    string
		message string
		fields  map[string]interface{}
		skip    bool
	}{
		{
			name:    "fatal with fields",
			message: "fatal message",
			fields: map[string]interface{}{
				"key": "value",
			},
			skip: true, // Skip actual execution as it would exit
		},
		{
			name:    "fatal without fields",
			message: "fatal message",
			fields:  nil,
			skip:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip("Skipping fatal test as it would exit the process")
			}
			loggerInstance.Fatal(tt.message, tt.fields)
		})
	}
}

func TestLogger_Logger_WithSpanContext(t *testing.T) {
	loggerInstance, err := NewLogger()
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	tracerInstance, err := tracer.NewTracer(tracer.WithServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewTracer() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = tracerInstance.Shutdown(ctx)
	}()

	ctx := context.Background()
	_, span := tracerInstance.StartSpan(ctx, "test-operation")
	defer span.End()

	spanContext := span.SpanContext()

	tests := []struct {
		name           string
		spanContext    trace.SpanContext
		checkFunc      func(t *testing.T, logger Logger)
		wantNil        bool
		wantSameLogger bool
	}{
		{
			name:        "with valid span context",
			spanContext: spanContext,
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("WithSpanContext() returned nil logger")
				}
				logger.Info("test message", map[string]interface{}{
					"test": "value",
				})
			},
			wantNil:        false,
			wantSameLogger: false,
		},
		{
			name:        "with invalid span context",
			spanContext: trace.SpanContext{},
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("WithSpanContext() returned nil logger")
				}
				logger.Info("test message", nil)
			},
			wantNil:        false,
			wantSameLogger: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := loggerInstance.WithSpanContext(tt.spanContext)
			if (got == nil) != tt.wantNil {
				t.Errorf("WithSpanContext() returned nil = %v, want nil = %v", got == nil, tt.wantNil)
				return
			}
			if got == nil {
				return
			}

			if (loggerInstance == got) != tt.wantSameLogger {
				t.Errorf("WithSpanContext() returned same logger = %v, want same = %v", loggerInstance == got, tt.wantSameLogger)
			}

			if tt.checkFunc != nil {
				tt.checkFunc(t, got)
			}
		})
	}
}

func TestLogger_Logger_Sync(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T) Logger
		wantErr bool
	}{
		{
			name: "sync valid logger",
			setup: func(t *testing.T) Logger {
				loggerInstance, err := NewLogger()
				if err != nil {
					t.Fatalf("NewLogger() error = %v", err)
				}
				return loggerInstance
			},
			wantErr: false,
		},
		{
			name: "sync nil logger",
			setup: func(t *testing.T) Logger {
				return (*logger)(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loggerInstance := tt.setup(t)
			err := loggerInstance.Sync()
			// Sync() may return errors on some systems (e.g., when stderr is redirected)
			// This is expected behavior, so we only check that nil logger doesn't error
			if loggerInstance == nil && err != nil {
				t.Errorf("Sync() on nil logger error = %v, want nil", err)
			}
		})
	}
}

func TestLogger_Logger_AllLogLevels(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error", "fatal"}

	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			loggerInstance, err := NewLogger(WithLevel(level))
			if err != nil {
				t.Fatalf("NewLogger() with level %s error = %v", level, err)
			}

			fields := map[string]interface{}{
				"level": level,
			}

			loggerInstance.Debug("debug message", fields)
			loggerInstance.Info("info message", fields)
			loggerInstance.Warn("warn message", fields)
			loggerInstance.Error("error message", fields)
		})
	}
}
