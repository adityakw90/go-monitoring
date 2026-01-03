package monitoring

import "time"

// Options contains all configuration for monitoring components.
// It is used internally by NewMonitoring and should be configured using Option functions.
type Options struct {
	ServiceName        string        // ServiceName is the name of the service (required).
	Environment        string        // Environment is the deployment environment (e.g., "development", "production").
	InstanceName       string        // InstanceName is the unique identifier for this service instance.
	InstanceHost       string        // InstanceHost is the hostname where this service instance is running.
	LoggerLevel        string        // LoggerLevel is the minimum log level to output. Valid values: "debug", "info", "warn", "error", "fatal".
	LoggerOutputPath   string        // LoggerOutputPath is the file path where logs will be written. If empty, logs will be written to stdout.
	TracerProvider     string        // TracerProvider specifies the trace exporter to use ("stdout" or "otlp").
	TracerProviderHost string        // TracerProviderHost is the hostname of the OTLP trace collector.
	TracerProviderPort int           // TracerProviderPort is the port of the OTLP trace collector.
	TracerSampleRatio  float64       // TracerSampleRatio controls the sampling rate for traces (0.0 to 1.0). 0.0 means never sample, 1.0 means always sample.
	TracerBatchTimeout time.Duration // TracerBatchTimeout is the maximum time to wait before exporting a batch of spans.
	TracerInsecure     bool          // TracerInsecure controls whether to use an insecure (non-TLS) connection for OTLP exporter.
	MetricProvider     string        // MetricProvider specifies the metric exporter to use ("stdout" or "otlp").
	MetricProviderHost string        // MetricProviderHost is the hostname of the OTLP metric collector.
	MetricProviderPort int           // MetricProviderPort is the port of the OTLP metric collector.
	MetricInterval     time.Duration // MetricInterval is the time interval between metric exports.
	MetricInsecure     bool          // MetricInsecure controls whether to use an insecure (non-TLS) connection for OTLP exporter.
}

// Option is a function that configures Options.
// It follows the functional options pattern for flexible configuration.
type Option func(*Options)

// WithServiceName sets the service name.
// This is a required option and must be provided when creating a Monitoring instance.
//
// Parameters:
//   - name: The name of the service (e.g., "user-service", "api-gateway")
//
// Example:
//
//	mon, err := NewMonitoring(
//	    WithServiceName("my-service"),
//	)
func WithServiceName(name string) Option {
	return func(o *Options) {
		o.ServiceName = name
	}
}

// WithEnvironment sets the deployment environment.
// This is used to tag traces and metrics with environment information.
//
// Parameters:
//   - env: The environment name (e.g., "development", "staging", "production")
//
// Example:
//
//	mon, err := NewMonitoring(
//	    WithServiceName("my-service"),
//	    WithEnvironment("production"),
//	)
func WithEnvironment(env string) Option {
	return func(o *Options) {
		o.Environment = env
	}
}

// WithInstance sets the instance name and host.
// This is used to identify the specific service instance in distributed systems.
//
// Parameters:
//   - name: The unique identifier for this instance (e.g., "instance-1", "pod-abc123")
//   - host: The hostname where this instance is running (e.g., "localhost", "10.0.0.1")
//
// Example:
//
//	mon, err := NewMonitoring(
//	    WithServiceName("my-service"),
//	    WithInstance("instance-1", "localhost"),
//	)
func WithInstance(name, host string) Option {
	return func(o *Options) {
		o.InstanceName = name
		o.InstanceHost = host
	}
}

// WithLoggerLevel sets the logger level.
// Only log messages at or above this level will be output.
//
// Parameters:
//   - level: The log level ("debug", "info", "warn", "error", "fatal")
//
// Example:
//
//	mon, err := NewMonitoring(
//	    WithServiceName("my-service"),
//	    WithLoggerLevel("debug"),
// WithLoggerLevel returns an Option that sets the logger minimum level for monitoring
// (e.g., "debug", "info", "warn", "error", "fatal").
func WithLoggerLevel(level string) Option {
	return func(o *Options) {
		o.LoggerLevel = level
	}
}

// WithLoggerOutputPath sets the logger output path.
// Logs will be written to the specified file path. If not set, logs will be written to stdout.
//
// Parameters:
//   - path: The file path where logs will be written (e.g., "/var/log/app.log", "./logs/app.log")
//
// Example:
//
//	mon, err := NewMonitoring(
//	    WithServiceName("my-service"),
//	    WithLoggerOutputPath("/var/log/app.log"),
// WithLoggerOutputPath returns an Option that sets the file path used for log output.
// If the provided path is empty, logs will be written to stdout.
func WithLoggerOutputPath(path string) Option {
	return func(o *Options) {
		o.LoggerOutputPath = path
	}
}

// WithTracerProvider sets the tracer provider configuration.
// This determines where traces are exported (stdout for development, OTLP for production).
//
// Parameters:
//   - provider: The provider type ("stdout" or "otlp")
//   - host: The hostname of the OTLP collector (ignored for "stdout")
//   - port: The port of the OTLP collector (ignored for "stdout")
//
// Example:
//
//	mon, err := NewMonitoring(
//	    WithServiceName("my-service"),
//	    WithTracerProvider("otlp", "localhost", 4317),
//	)
func WithTracerProvider(provider, host string, port int) Option {
	return func(o *Options) {
		o.TracerProvider = provider
		o.TracerProviderHost = host
		o.TracerProviderPort = port
	}
}

// WithTracerSampleRatio sets the tracer sampling ratio.
// This controls what percentage of traces are sampled and exported.
//
// Parameters:
//   - ratio: Sampling ratio between 0.0 and 1.0
//   - 0.0: Never sample (no traces exported)
//   - 1.0: Always sample (all traces exported)
//   - 0.1: Sample 10% of traces
//
// Example:
//
//	mon, err := NewMonitoring(
//	    WithServiceName("my-service"),
//	    WithTracerSampleRatio(0.1), // Sample 10% of traces
//	)
func WithTracerSampleRatio(ratio float64) Option {
	return func(o *Options) {
		o.TracerSampleRatio = ratio
	}
}

// WithTracerBatchTimeout sets the tracer batch timeout.
// This is the maximum time to wait before exporting a batch of spans.
// Longer timeouts allow more spans to be batched together, improving efficiency.
//
// Parameters:
//   - timeout: The maximum time to wait before exporting a batch
//
// Example:
//
//	mon, err := NewMonitoring(
//	    WithServiceName("my-service"),
//	    WithTracerBatchTimeout(10*time.Second),
//	)
func WithTracerBatchTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.TracerBatchTimeout = timeout
	}
}

// WithTracerInsecure sets whether to use an insecure (non-TLS) connection for OTLP exporter.
// When false (default), a secure TLS connection is used. When true, connections are made without TLS.
// This should only be used in development or when TLS is handled by a proxy.
//
// Parameters:
//   - insecure: Whether to use an insecure connection
//
// Example:
//
//	mon, err := NewMonitoring(
//	    WithServiceName("my-service"),
//	    WithTracerProvider("otlp", "localhost", 4317),
//	    WithTracerInsecure(true), // Use insecure connection
//	)
func WithTracerInsecure(insecure bool) Option {
	return func(o *Options) {
		o.TracerInsecure = insecure
	}
}

// WithMetricProvider sets the metric provider configuration.
// This determines where metrics are exported (stdout for development, OTLP for production).
//
// Parameters:
//   - provider: The provider type ("stdout" or "otlp")
//   - host: The hostname of the OTLP collector (ignored for "stdout")
//   - port: The port of the OTLP collector (ignored for "stdout")
//
// Example:
//
//	mon, err := NewMonitoring(
//	    WithServiceName("my-service"),
//	    WithMetricProvider("otlp", "localhost", 4318),
//	)
func WithMetricProvider(provider, host string, port int) Option {
	return func(o *Options) {
		o.MetricProvider = provider
		o.MetricProviderHost = host
		o.MetricProviderPort = port
	}
}

// WithMetricInterval sets the metric export interval.
// This determines how frequently metrics are exported to the configured provider.
// Shorter intervals provide more real-time metrics but increase overhead.
//
// Parameters:
//   - interval: The time interval between metric exports
//
// Example:
//
//	mon, err := NewMonitoring(
//	    WithServiceName("my-service"),
//	    WithMetricInterval(30*time.Second), // Export every 30 seconds
//	)
func WithMetricInterval(interval time.Duration) Option {
	return func(o *Options) {
		o.MetricInterval = interval
	}
}

// WithMetricInsecure sets whether to use an insecure (non-TLS) connection for OTLP exporter.
// When false (default), a secure TLS connection is used. When true, connections are made without TLS.
// This should only be used in development or when TLS is handled by a proxy.
//
// Parameters:
//   - insecure: Whether to use an insecure connection
//
// Example:
//
//	mon, err := NewMonitoring(
//	    WithServiceName("my-service"),
//	    WithMetricProvider("otlp", "localhost", 4318),
//	    WithMetricInsecure(true), // Use insecure connection
//	)
func WithMetricInsecure(insecure bool) Option {
	return func(o *Options) {
		o.MetricInsecure = insecure
	}
}

// defaultOptions returns a pointer to Options populated with sensible defaults for monitoring components.
// The defaults set the environment to "development", logger level to "info" with an empty LoggerOutputPath (use stdout),
// tracer and metric providers to "stdout", tracer sample ratio to 1.0, tracer batch timeout to 5s, and metric export
// interval to 60s.
func defaultOptions() *Options {
	return &Options{
		Environment:        "development",
		LoggerLevel:        "info",
		LoggerOutputPath:   "",
		TracerProvider:     "stdout",
		TracerSampleRatio:  1.0,
		TracerBatchTimeout: 5 * time.Second,
		MetricProvider:     "stdout",
		MetricInterval:     60 * time.Second,
	}
}