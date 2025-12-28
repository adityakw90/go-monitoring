package monitoring

import "time"

// Options contains all configuration for monitoring components.
type Options struct {
	// Service information
	ServiceName  string
	Environment  string
	InstanceName string
	InstanceHost string

	// Logger configuration
	LoggerLevel string // debug, info, warn, error, fatal

	// Tracer configuration
	TracerProvider     string // "stdout", "otlp"
	TracerProviderHost string
	TracerProviderPort int
	TracerSampleRatio  float64 // 0.0 to 1.0
	TracerBatchTimeout time.Duration

	// Metric configuration
	MetricProvider     string // "stdout", "otlp"
	MetricProviderHost string
	MetricProviderPort int
	MetricInterval     time.Duration
}

// Option is a function that configures Options.
type Option func(*Options)

// WithServiceName sets the service name.
func WithServiceName(name string) Option {
	return func(o *Options) {
		o.ServiceName = name
	}
}

// WithEnvironment sets the environment (e.g., "development", "production").
func WithEnvironment(env string) Option {
	return func(o *Options) {
		o.Environment = env
	}
}

// WithInstance sets the instance name and host.
func WithInstance(name, host string) Option {
	return func(o *Options) {
		o.InstanceName = name
		o.InstanceHost = host
	}
}

// WithLoggerLevel sets the logger level (debug, info, warn, error, fatal).
func WithLoggerLevel(level string) Option {
	return func(o *Options) {
		o.LoggerLevel = level
	}
}

// WithTracerProvider sets the tracer provider configuration.
func WithTracerProvider(provider, host string, port int) Option {
	return func(o *Options) {
		o.TracerProvider = provider
		o.TracerProviderHost = host
		o.TracerProviderPort = port
	}
}

// WithTracerSampleRatio sets the tracer sampling ratio (0.0 to 1.0).
func WithTracerSampleRatio(ratio float64) Option {
	return func(o *Options) {
		o.TracerSampleRatio = ratio
	}
}

// WithTracerBatchTimeout sets the tracer batch timeout.
func WithTracerBatchTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.TracerBatchTimeout = timeout
	}
}

// WithMetricProvider sets the metric provider configuration.
func WithMetricProvider(provider, host string, port int) Option {
	return func(o *Options) {
		o.MetricProvider = provider
		o.MetricProviderHost = host
		o.MetricProviderPort = port
	}
}

// WithMetricInterval sets the metric export interval.
func WithMetricInterval(interval time.Duration) Option {
	return func(o *Options) {
		o.MetricInterval = interval
	}
}

// defaultOptions returns Options with sensible defaults.
func defaultOptions() *Options {
	return &Options{
		Environment:        "development",
		LoggerLevel:        "info",
		TracerProvider:     "stdout",
		TracerSampleRatio:  1.0,
		TracerBatchTimeout: 5 * time.Second,
		MetricProvider:     "stdout",
		MetricInterval:     60 * time.Second,
	}
}
