package monitoring

import (
	"errors"
	"testing"

	"github.com/adityakw90/go-monitoring/internal/logger"
	"github.com/adityakw90/go-monitoring/internal/metric"
	"github.com/adityakw90/go-monitoring/internal/tracer"
)

func TestMonitoring_Errors_ParseError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		message  string
		validate func(*testing.T, error)
	}{
		{
			name:    "logger invalid log level",
			err:     logger.ErrInvalidLogLevel,
			message: "test message",
			validate: func(t *testing.T, got error) {
				if !errors.Is(got, ErrLoggerInvalidLogLevel) {
					t.Errorf("expected ErrLoggerInvalidLogLevel, got %v", got)
				}
				// parseError returns ErrLoggerInvalidLogLevel directly, not wrapped
				if got != ErrLoggerInvalidLogLevel {
					t.Errorf("expected direct ErrLoggerInvalidLogLevel, got %v", got)
				}
			},
		},
		{
			name:    "tracer invalid provider",
			err:     tracer.ErrInvalidProvider,
			message: "test message",
			validate: func(t *testing.T, got error) {
				if !errors.Is(got, ErrTracerInvalidProvider) {
					t.Errorf("expected ErrTracerInvalidProvider, got %v", got)
				}
				// parseError returns ErrTracerInvalidProvider directly, not wrapped
				if got != ErrTracerInvalidProvider {
					t.Errorf("expected direct ErrTracerInvalidProvider, got %v", got)
				}
			},
		},
		{
			name:    "metric invalid provider",
			err:     metric.ErrInvalidProvider,
			message: "test message",
			validate: func(t *testing.T, got error) {
				if !errors.Is(got, ErrMetricInvalidProvider) {
					t.Errorf("expected ErrMetricInvalidProvider, got %v", got)
				}
				// parseError returns ErrMetricInvalidProvider directly, not wrapped
				if got != ErrMetricInvalidProvider {
					t.Errorf("expected direct ErrMetricInvalidProvider, got %v", got)
				}
			},
		},
		{
			name:    "tracer provider host required",
			err:     tracer.ErrProviderHostRequired,
			message: "test message",
			validate: func(t *testing.T, got error) {
				if !errors.Is(got, ErrTracerProviderHostRequired) {
					t.Errorf("expected ErrTracerProviderHostRequired, got %v", got)
				}
				// parseError returns ErrTracerProviderHostRequired directly, not wrapped
				if got != ErrTracerProviderHostRequired {
					t.Errorf("expected direct ErrTracerProviderHostRequired, got %v", got)
				}
			},
		},
		{
			name:    "tracer provider port required",
			err:     tracer.ErrProviderPortRequired,
			message: "test message",
			validate: func(t *testing.T, got error) {
				if !errors.Is(got, ErrTracerProviderPortRequired) {
					t.Errorf("expected ErrTracerProviderPortRequired, got %v", got)
				}
				// parseError returns ErrTracerProviderPortRequired directly, not wrapped
				if got != ErrTracerProviderPortRequired {
					t.Errorf("expected direct ErrTracerProviderPortRequired, got %v", got)
				}
			},
		},
		{
			name:    "tracer provider port invalid",
			err:     tracer.ErrProviderPortInvalid,
			message: "test message",
			validate: func(t *testing.T, got error) {
				if !errors.Is(got, ErrTracerProviderPortInvalid) {
					t.Errorf("expected ErrTracerProviderPortInvalid, got %v", got)
				}
				// parseError returns ErrTracerProviderPortInvalid directly, not wrapped
				if got != ErrTracerProviderPortInvalid {
					t.Errorf("expected direct ErrTracerProviderPortInvalid, got %v", got)
				}
			},
		},
		{
			name:    "tracer batch timeout invalid",
			err:     tracer.ErrBatchTimeoutInvalid,
			message: "test message",
			validate: func(t *testing.T, got error) {
				if !errors.Is(got, ErrTracerBatchTimeoutInvalid) {
					t.Errorf("expected ErrTracerBatchTimeoutInvalid, got %v", got)
				}
				// parseError returns ErrTracerBatchTimeoutInvalid directly, not wrapped
				if got != ErrTracerBatchTimeoutInvalid {
					t.Errorf("expected direct ErrTracerBatchTimeoutInvalid, got %v", got)
				}
			},
		},
		{
			name:    "metric provider host required",
			err:     metric.ErrProviderHostRequired,
			message: "test message",
			validate: func(t *testing.T, got error) {
				if !errors.Is(got, ErrMetricProviderHostRequired) {
					t.Errorf("expected ErrMetricProviderHostRequired, got %v", got)
				}
				// parseError returns ErrMetricProviderHostRequired directly, not wrapped
				if got != ErrMetricProviderHostRequired {
					t.Errorf("expected direct ErrMetricProviderHostRequired, got %v", got)
				}
			},
		},
		{
			name:    "metric provider port required",
			err:     metric.ErrProviderPortRequired,
			message: "test message",
			validate: func(t *testing.T, got error) {
				if !errors.Is(got, ErrMetricProviderPortRequired) {
					t.Errorf("expected ErrMetricProviderPortRequired, got %v", got)
				}
				// parseError returns ErrMetricProviderPortRequired directly, not wrapped
				if got != ErrMetricProviderPortRequired {
					t.Errorf("expected direct ErrMetricProviderPortRequired, got %v", got)
				}
			},
		},
		{
			name:    "metric provider port invalid",
			err:     metric.ErrProviderPortInvalid,
			message: "test message",
			validate: func(t *testing.T, got error) {
				if !errors.Is(got, ErrMetricProviderPortInvalid) {
					t.Errorf("expected ErrMetricProviderPortInvalid, got %v", got)
				}
				// parseError returns ErrMetricProviderPortInvalid directly, not wrapped
				if got != ErrMetricProviderPortInvalid {
					t.Errorf("expected direct ErrMetricProviderPortInvalid, got %v", got)
				}
			},
		},
		{
			name:    "metric interval invalid",
			err:     metric.ErrIntervalInvalid,
			message: "test message",
			validate: func(t *testing.T, got error) {
				if !errors.Is(got, ErrMetricIntervalInvalid) {
					t.Errorf("expected ErrMetricIntervalInvalid, got %v", got)
				}
				// parseError returns ErrMetricIntervalInvalid directly, not wrapped
				if got != ErrMetricIntervalInvalid {
					t.Errorf("expected direct ErrMetricIntervalInvalid, got %v", got)
				}
			},
		},
		{
			name:    "generic error wrapped",
			err:     errors.New("some error"),
			message: "failed to initialize",
			validate: func(t *testing.T, got error) {
				if got == nil {
					t.Error("expected error, got nil")
					return
				}
				// Generic errors are wrapped with fmt.Errorf using %w
				// Verify the error message format
				errMsg := got.Error()
				if errMsg != "failed to initialize: some error" {
					t.Errorf("expected error message 'failed to initialize: some error', got %q", errMsg)
				}
				// Verify it's a wrapped error by checking it unwraps correctly
				unwrapped := errors.Unwrap(got)
				if unwrapped == nil {
					t.Error("expected wrapped error with Unwrap, got nil")
				} else if unwrapped.Error() != "some error" {
					t.Errorf("expected unwrapped error 'some error', got %q", unwrapped.Error())
				}
			},
		},
		{
			name:    "nil error",
			err:     nil,
			message: "test message",
			validate: func(t *testing.T, got error) {
				if got == nil {
					t.Error("expected error, got nil")
					return
				}
				// parseError returns fmt.Errorf("%s: unknown error", message) for nil errors
				errMsg := got.Error()
				if errMsg != "test message: unknown error" {
					t.Errorf("expected error message 'test message: unknown error', got %q", errMsg)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseError(tt.err, tt.message)
			tt.validate(t, got)
		})
	}
}
