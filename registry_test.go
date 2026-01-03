package monitoring

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/adityakw90/go-monitoring/internal/logger"
	"github.com/adityakw90/go-monitoring/internal/metric"
	"github.com/adityakw90/go-monitoring/internal/tracer"
)

func TestMonitoring_Registry_ParseOptions(t *testing.T) {
	tests := []struct {
		name     string
		opts     []Option
		validate func(*testing.T, *Options)
	}{
		{
			name: "default options",
			opts: nil,
			validate: func(t *testing.T, o *Options) {
				if o.Environment != "development" {
					t.Errorf("expected Environment = 'development', got %q", o.Environment)
				}
				if o.LoggerLevel != "info" {
					t.Errorf("expected LoggerLevel = 'info', got %q", o.LoggerLevel)
				}
				if o.TracerProvider != "stdout" {
					t.Errorf("expected TracerProvider = 'stdout', got %q", o.TracerProvider)
				}
				if o.TracerSampleRatio != 1.0 {
					t.Errorf("expected TracerSampleRatio = 1.0, got %f", o.TracerSampleRatio)
				}
				if o.TracerBatchTimeout != 5*time.Second {
					t.Errorf("expected TracerBatchTimeout = 5s, got %v", o.TracerBatchTimeout)
				}
				if o.MetricProvider != "stdout" {
					t.Errorf("expected MetricProvider = 'stdout', got %q", o.MetricProvider)
				}
				if o.MetricInterval != 60*time.Second {
					t.Errorf("expected MetricInterval = 60s, got %v", o.MetricInterval)
				}
			},
		},
		{
			name: "with service name",
			opts: []Option{
				WithServiceName("test-service"),
			},
			validate: func(t *testing.T, o *Options) {
				if o.ServiceName != "test-service" {
					t.Errorf("expected ServiceName = 'test-service', got %q", o.ServiceName)
				}
			},
		},
		{
			name: "with all options",
			opts: []Option{
				WithServiceName("my-service"),
				WithEnvironment("production"),
				WithInstance("instance-1", "localhost"),
				WithLoggerLevel("debug"),
				WithTracerProvider("otlp", "localhost", 4317),
				WithTracerSampleRatio(0.5),
				WithTracerBatchTimeout(10 * time.Second),
				WithTracerInsecure(true),
				WithMetricProvider("otlp", "localhost", 4318),
				WithMetricInterval(30 * time.Second),
				WithMetricInsecure(true),
			},
			validate: func(t *testing.T, o *Options) {
				if o.ServiceName != "my-service" {
					t.Errorf("expected ServiceName = 'my-service', got %q", o.ServiceName)
				}
				if o.Environment != "production" {
					t.Errorf("expected Environment = 'production', got %q", o.Environment)
				}
				if o.InstanceName != "instance-1" {
					t.Errorf("expected InstanceName = 'instance-1', got %q", o.InstanceName)
				}
				if o.InstanceHost != "localhost" {
					t.Errorf("expected InstanceHost = 'localhost', got %q", o.InstanceHost)
				}
				if o.LoggerLevel != "debug" {
					t.Errorf("expected LoggerLevel = 'debug', got %q", o.LoggerLevel)
				}
				if o.TracerProvider != "otlp" {
					t.Errorf("expected TracerProvider = 'otlp', got %q", o.TracerProvider)
				}
				if o.TracerProviderHost != "localhost" {
					t.Errorf("expected TracerProviderHost = 'localhost', got %q", o.TracerProviderHost)
				}
				if o.TracerProviderPort != 4317 {
					t.Errorf("expected TracerProviderPort = 4317, got %d", o.TracerProviderPort)
				}
				if o.TracerSampleRatio != 0.5 {
					t.Errorf("expected TracerSampleRatio = 0.5, got %f", o.TracerSampleRatio)
				}
				if o.TracerBatchTimeout != 10*time.Second {
					t.Errorf("expected TracerBatchTimeout = 10s, got %v", o.TracerBatchTimeout)
				}
				if !o.TracerInsecure {
					t.Errorf("expected TracerInsecure = true, got %v", o.TracerInsecure)
				}
				if o.MetricProvider != "otlp" {
					t.Errorf("expected MetricProvider = 'otlp', got %q", o.MetricProvider)
				}
				if o.MetricProviderHost != "localhost" {
					t.Errorf("expected MetricProviderHost = 'localhost', got %q", o.MetricProviderHost)
				}
				if o.MetricProviderPort != 4318 {
					t.Errorf("expected MetricProviderPort = 4318, got %d", o.MetricProviderPort)
				}
				if o.MetricInterval != 30*time.Second {
					t.Errorf("expected MetricInterval = 30s, got %v", o.MetricInterval)
				}
				if !o.MetricInsecure {
					t.Errorf("expected MetricInsecure = true, got %v", o.MetricInsecure)
				}
			},
		},
		{
			name: "multiple options override",
			opts: []Option{
				WithServiceName("first"),
				WithServiceName("second"),
				WithEnvironment("dev"),
				WithEnvironment("prod"),
			},
			validate: func(t *testing.T, o *Options) {
				if o.ServiceName != "second" {
					t.Errorf("expected ServiceName = 'second', got %q", o.ServiceName)
				}
				if o.Environment != "prod" {
					t.Errorf("expected Environment = 'prod', got %q", o.Environment)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := parseOptions(tt.opts...)
			tt.validate(t, options)
		})
	}
}

func TestMonitoring_Registry_NewLogger(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
		check   func(*testing.T, Logger)
	}{
		{
			name:    "success with default level",
			opts:    nil,
			wantErr: false,
			check: func(t *testing.T, l Logger) {
				if l == nil {
					t.Error("expected logger, got nil")
				}
			},
		},
		{
			name: "success with info level",
			opts: []Option{
				WithLoggerLevel("info"),
			},
			wantErr: false,
			check: func(t *testing.T, l Logger) {
				if l == nil {
					t.Error("expected logger, got nil")
				}
			},
		},
		{
			name: "success with debug level",
			opts: []Option{
				WithLoggerLevel("debug"),
			},
			wantErr: false,
			check: func(t *testing.T, l Logger) {
				if l == nil {
					t.Error("expected logger, got nil")
				}
			},
		},
		{
			name: "success with warn level",
			opts: []Option{
				WithLoggerLevel("warn"),
			},
			wantErr: false,
			check: func(t *testing.T, l Logger) {
				if l == nil {
					t.Error("expected logger, got nil")
				}
			},
		},
		{
			name: "success with error level",
			opts: []Option{
				WithLoggerLevel("error"),
			},
			wantErr: false,
			check: func(t *testing.T, l Logger) {
				if l == nil {
					t.Error("expected logger, got nil")
				}
			},
		},
		{
			name: "success with fatal level",
			opts: []Option{
				WithLoggerLevel("fatal"),
			},
			wantErr: false,
			check: func(t *testing.T, l Logger) {
				if l == nil {
					t.Error("expected logger, got nil")
				}
			},
		},
		{
			name: "invalid log level",
			opts: []Option{
				WithLoggerLevel("invalid"),
			},
			wantErr: true,
			check: func(t *testing.T, l Logger) {
				if l != nil {
					t.Error("expected nil logger on error")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loggerInstance, err := NewLogger(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if !errors.Is(err, ErrLoggerInvalidLogLevel) && tt.wantErr {
					// Check if it's wrapped properly
					if !errors.Is(err, logger.ErrInvalidLogLevel) {
						t.Errorf("expected ErrLoggerInvalidLogLevel or wrapped error, got %v", err)
					}
				}
			}
			if loggerInstance != nil {
				defer func() {
					_ = loggerInstance.Sync()
				}()
			}
			tt.check(t, loggerInstance)
		})
	}
}

func TestMonitoring_Registry_NewTracer(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
		check   func(*testing.T, Tracer)
	}{
		{
			name: "success with stdout provider",
			opts: []Option{
				WithServiceName("test-service"),
				WithTracerProvider("stdout", "", 0),
			},
			wantErr: false,
			check: func(t *testing.T, tr Tracer) {
				if tr == nil {
					t.Error("expected tracer, got nil")
				}
			},
		},
		{
			name: "success with all options",
			opts: []Option{
				WithServiceName("test-service"),
				WithEnvironment("production"),
				WithInstance("instance-1", "localhost"),
				WithTracerProvider("stdout", "", 0),
				WithTracerSampleRatio(0.5),
				WithTracerBatchTimeout(10 * time.Second),
				WithTracerInsecure(true),
			},
			wantErr: false,
			check: func(t *testing.T, tr Tracer) {
				if tr == nil {
					t.Error("expected tracer, got nil")
				}
			},
		},
		{
			name: "invalid provider",
			opts: []Option{
				WithServiceName("test-service"),
				WithTracerProvider("invalid", "", 0),
			},
			wantErr: true,
			check: func(t *testing.T, tr Tracer) {
				if tr != nil {
					t.Error("expected nil tracer on error")
				}
			},
		},
		{
			name: "missing service name",
			opts: []Option{
				WithTracerProvider("stdout", "", 0),
			},
			wantErr: false, // Service name is not required for tracer alone
			check: func(t *testing.T, tr Tracer) {
				if tr == nil {
					t.Error("expected tracer, got nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracerInstance, err := NewTracer(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTracer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if !errors.Is(err, ErrTracerInvalidProvider) && tt.wantErr {
					// Check if it's wrapped properly
					if !errors.Is(err, tracer.ErrInvalidProvider) {
						t.Errorf("expected ErrTracerInvalidProvider or wrapped error, got %v", err)
					}
				}
			}
			if tracerInstance != nil {
				defer func() {
					_ = tracerInstance.Shutdown(context.Background())
				}()
			}
			tt.check(t, tracerInstance)
		})
	}
}

func TestMonitoring_Registry_NewMetric(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
		check   func(*testing.T, Metric)
	}{
		{
			name: "success with stdout provider",
			opts: []Option{
				WithServiceName("test-service"),
				WithMetricProvider("stdout", "", 0),
			},
			wantErr: false,
			check: func(t *testing.T, m Metric) {
				if m == nil {
					t.Error("expected metric, got nil")
				}
			},
		},
		{
			name: "success with all options",
			opts: []Option{
				WithServiceName("test-service"),
				WithEnvironment("production"),
				WithInstance("instance-1", "localhost"),
				WithMetricProvider("stdout", "", 0),
				WithMetricInterval(30 * time.Second),
				WithMetricInsecure(true),
			},
			wantErr: false,
			check: func(t *testing.T, m Metric) {
				if m == nil {
					t.Error("expected metric, got nil")
				}
			},
		},
		{
			name: "invalid provider",
			opts: []Option{
				WithServiceName("test-service"),
				WithMetricProvider("invalid", "", 0),
			},
			wantErr: true,
			check: func(t *testing.T, m Metric) {
				if m != nil {
					t.Error("expected nil metric on error")
				}
			},
		},
		{
			name: "missing service name",
			opts: []Option{
				WithMetricProvider("stdout", "", 0),
			},
			wantErr: false, // Service name is not required for metric alone
			check: func(t *testing.T, m Metric) {
				if m == nil {
					t.Error("expected metric, got nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricInstance, err := NewMetric(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if !errors.Is(err, ErrMetricInvalidProvider) && tt.wantErr {
					// Check if it's wrapped properly
					if !errors.Is(err, metric.ErrInvalidProvider) {
						t.Errorf("expected ErrMetricInvalidProvider or wrapped error, got %v", err)
					}
				}
			}
			if metricInstance != nil {
				defer func() {
					_ = metricInstance.Shutdown(context.Background())
				}()
			}
			tt.check(t, metricInstance)
		})
	}
}

func TestMonitoring_Registry_NewMonitoring(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
		check   func(*testing.T, *Monitoring)
	}{
		{
			name:    "missing service name",
			opts:    nil,
			wantErr: true,
			check: func(t *testing.T, m *Monitoring) {
				if m != nil {
					t.Error("expected nil monitoring on error")
				}
			},
		},
		{
			name: "with service name",
			opts: []Option{
				WithServiceName("test-service"),
			},
			wantErr: false,
			check: func(t *testing.T, m *Monitoring) {
				if m == nil {
					t.Error("expected monitoring, got nil")
					return
				}
				if m.Logger == nil {
					t.Error("expected Logger, got nil")
				}
				if m.Tracer == nil {
					t.Error("expected Tracer, got nil")
				}
				if m.Metric == nil {
					t.Error("expected Metric, got nil")
				}
			},
		},
		{
			name: "with all options",
			opts: []Option{
				WithServiceName("test-service"),
				WithEnvironment("test"),
				WithInstance("instance-1", "localhost"),
				WithLoggerLevel("debug"),
				WithTracerProvider("stdout", "", 0),
				WithTracerSampleRatio(0.5),
				WithTracerBatchTimeout(10 * time.Second),
				WithTracerInsecure(true),
				WithMetricProvider("stdout", "", 0),
				WithMetricInterval(30 * time.Second),
				WithMetricInsecure(true),
			},
			wantErr: false,
			check: func(t *testing.T, m *Monitoring) {
				if m == nil {
					t.Error("expected monitoring, got nil")
					return
				}
				if m.Logger == nil {
					t.Error("expected Logger, got nil")
				}
				if m.Tracer == nil {
					t.Error("expected Tracer, got nil")
				}
				if m.Metric == nil {
					t.Error("expected Metric, got nil")
				}
			},
		},
		{
			name: "invalid logger level",
			opts: []Option{
				WithServiceName("test-service"),
				WithLoggerLevel("invalid"),
			},
			wantErr: true,
			check: func(t *testing.T, m *Monitoring) {
				if m != nil {
					t.Error("expected nil monitoring on error")
				}
			},
		},
		{
			name: "invalid tracer provider",
			opts: []Option{
				WithServiceName("test-service"),
				WithTracerProvider("invalid", "", 0),
			},
			wantErr: true,
			check: func(t *testing.T, m *Monitoring) {
				if m != nil {
					t.Error("expected nil monitoring on error")
				}
			},
		},
		{
			name: "invalid metric provider",
			opts: []Option{
				WithServiceName("test-service"),
				WithMetricProvider("invalid", "", 0),
			},
			wantErr: true,
			check: func(t *testing.T, m *Monitoring) {
				if m != nil {
					t.Error("expected nil monitoring on error")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monitoring, err := NewMonitoring(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMonitoring() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if monitoring != nil {
				defer func() {
					if monitoring.Tracer != nil {
						_ = monitoring.Tracer.Shutdown(context.Background())
					}
					if monitoring.Metric != nil {
						_ = monitoring.Metric.Shutdown(context.Background())
					}
					if monitoring.Logger != nil {
						_ = monitoring.Logger.Sync()
					}
				}()
			}
			tt.check(t, monitoring)
			if err != nil {
				// Verify error types are properly wrapped
				if tt.name == "invalid logger level" {
					if !errors.Is(err, ErrLoggerInvalidLogLevel) {
						t.Errorf("expected ErrLoggerInvalidLogLevel, got %v", err)
					}
				}
				if tt.name == "invalid tracer provider" {
					if !errors.Is(err, ErrTracerInvalidProvider) {
						t.Errorf("expected ErrTracerInvalidProvider, got %v", err)
					}
				}
				if tt.name == "invalid metric provider" {
					if !errors.Is(err, ErrMetricInvalidProvider) {
						t.Errorf("expected ErrMetricInvalidProvider, got %v", err)
					}
				}
				if tt.name == "missing service name" {
					if !errors.Is(err, ErrServiceNameRequired) {
						t.Errorf("expected ErrServiceNameRequired, got %v", err)
					}
				}
			}
		})
	}
}

func TestMonitoring_Registry_NewMonitoring_CleanupOnFailure(t *testing.T) {
	// This test verifies that when initialization fails, previously initialized components are cleaned up
	// We can't easily test this without mocking, but we can verify the error handling path exists

	t.Run("tracer failure cleans up logger", func(t *testing.T) {
		// Use invalid tracer provider to cause tracer initialization to fail
		// Logger should be initialized first, then tracer fails
		monitoring, err := NewMonitoring(
			WithServiceName("test-service"),
			WithTracerProvider("invalid", "", 0),
		)
		if err == nil {
			t.Error("expected error, got nil")
			if monitoring != nil {
				// Cleanup if somehow created
				defer func() {
					if monitoring.Tracer != nil {
						_ = monitoring.Tracer.Shutdown(context.Background())
					}
					if monitoring.Metric != nil {
						_ = monitoring.Metric.Shutdown(context.Background())
					}
					if monitoring.Logger != nil {
						_ = monitoring.Logger.Sync()
					}
				}()
			}
			return
		}
		if monitoring != nil {
			t.Error("expected nil monitoring on error")
		}
		if !errors.Is(err, ErrTracerInvalidProvider) {
			t.Errorf("expected ErrTracerInvalidProvider, got %v", err)
		}
	})

	t.Run("metric failure cleans up logger and tracer", func(t *testing.T) {
		// Use invalid metric provider to cause metric initialization to fail
		// Logger and tracer should be initialized first, then metric fails
		monitoring, err := NewMonitoring(
			WithServiceName("test-service"),
			WithMetricProvider("invalid", "", 0),
		)
		if err == nil {
			t.Error("expected error, got nil")
			if monitoring != nil {
				// Cleanup if somehow created
				defer func() {
					if monitoring.Tracer != nil {
						_ = monitoring.Tracer.Shutdown(context.Background())
					}
					if monitoring.Metric != nil {
						_ = monitoring.Metric.Shutdown(context.Background())
					}
					if monitoring.Logger != nil {
						_ = monitoring.Logger.Sync()
					}
				}()
			}
			return
		}
		if monitoring != nil {
			t.Error("expected nil monitoring on error")
		}
		if !errors.Is(err, ErrMetricInvalidProvider) {
			t.Errorf("expected ErrMetricInvalidProvider, got %v", err)
		}
	})
}
