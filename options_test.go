package monitoring

import (
	"testing"
	"time"
)

func TestDefaultOptions(t *testing.T) {
	opts := defaultOptions()

	if opts.Environment != "development" {
		t.Errorf("defaultOptions() Environment = %v, want development", opts.Environment)
	}
	if opts.LoggerLevel != "info" {
		t.Errorf("defaultOptions() LoggerLevel = %v, want info", opts.LoggerLevel)
	}
	if opts.TracerProvider != "stdout" {
		t.Errorf("defaultOptions() TracerProvider = %v, want stdout", opts.TracerProvider)
	}
	if opts.TracerSampleRatio != 1.0 {
		t.Errorf("defaultOptions() TracerSampleRatio = %v, want 1.0", opts.TracerSampleRatio)
	}
	if opts.TracerBatchTimeout != 5*time.Second {
		t.Errorf("defaultOptions() TracerBatchTimeout = %v, want 5s", opts.TracerBatchTimeout)
	}
	if opts.MetricProvider != "stdout" {
		t.Errorf("defaultOptions() MetricProvider = %v, want stdout", opts.MetricProvider)
	}
	if opts.MetricInterval != 60*time.Second {
		t.Errorf("defaultOptions() MetricInterval = %v, want 60s", opts.MetricInterval)
	}
	if opts.TracerInsecure != false {
		t.Errorf("defaultOptions() TracerInsecure = %v, want false", opts.TracerInsecure)
	}
}

func TestOptions(t *testing.T) {
	opts := defaultOptions()

	WithServiceName("test-service")(opts)
	if opts.ServiceName != "test-service" {
		t.Errorf("WithServiceName() ServiceName = %v, want test-service", opts.ServiceName)
	}

	WithEnvironment("production")(opts)
	if opts.Environment != "production" {
		t.Errorf("WithEnvironment() Environment = %v, want production", opts.Environment)
	}

	WithInstance("instance-1", "localhost")(opts)
	if opts.InstanceName != "instance-1" {
		t.Errorf("WithInstance() InstanceName = %v, want instance-1", opts.InstanceName)
	}
	if opts.InstanceHost != "localhost" {
		t.Errorf("WithInstance() InstanceHost = %v, want localhost", opts.InstanceHost)
	}

	WithLoggerLevel("debug")(opts)
	if opts.LoggerLevel != "debug" {
		t.Errorf("WithLoggerLevel() LoggerLevel = %v, want debug", opts.LoggerLevel)
	}

	WithTracerProvider("otlp", "localhost", 4317)(opts)
	if opts.TracerProvider != "otlp" {
		t.Errorf("WithTracerProvider() TracerProvider = %v, want otlp", opts.TracerProvider)
	}
	if opts.TracerProviderHost != "localhost" {
		t.Errorf("WithTracerProvider() TracerProviderHost = %v, want localhost", opts.TracerProviderHost)
	}
	if opts.TracerProviderPort != 4317 {
		t.Errorf("WithTracerProvider() TracerProviderPort = %v, want 4317", opts.TracerProviderPort)
	}

	WithTracerSampleRatio(0.5)(opts)
	if opts.TracerSampleRatio != 0.5 {
		t.Errorf("WithTracerSampleRatio() TracerSampleRatio = %v, want 0.5", opts.TracerSampleRatio)
	}

	WithTracerBatchTimeout(10 * time.Second)(opts)
	if opts.TracerBatchTimeout != 10*time.Second {
		t.Errorf("WithTracerBatchTimeout() TracerBatchTimeout = %v, want 10s", opts.TracerBatchTimeout)
	}

	WithTracerInsecure(true)(opts)
	if opts.TracerInsecure != true {
		t.Errorf("WithTracerInsecure() TracerInsecure = %v, want true", opts.TracerInsecure)
	}

	WithMetricProvider("otlp", "localhost", 4318)(opts)
	if opts.MetricProvider != "otlp" {
		t.Errorf("WithMetricProvider() MetricProvider = %v, want otlp", opts.MetricProvider)
	}

	WithMetricInterval(30 * time.Second)(opts)
	if opts.MetricInterval != 30*time.Second {
		t.Errorf("WithMetricInterval() MetricInterval = %v, want 30s", opts.MetricInterval)
	}
}
