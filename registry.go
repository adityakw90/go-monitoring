package monitoring

import (
	"context"
	"errors"
	"fmt"

	"github.com/adityakw90/go-monitoring/internal/logger"
	"github.com/adityakw90/go-monitoring/internal/metric"
	"github.com/adityakw90/go-monitoring/internal/tracer"
)

// parseOptions applies the given options to default options and returns the configured Options.
// It follows the functional options pattern, starting with defaults and applying each option in order.
func parseOptions(opts ...Option) *Options {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// parseError maps internal package errors to public API errors.
// It checks for known sentinel errors from internal packages and returns the corresponding
// public error. If the error is not a known sentinel error, it wraps it with the provided message.
func parseError(err error, message string) error {
	if errors.Is(err, logger.ErrInvalidLogLevel) {
		return ErrLoggerInvalidLogLevel
	}
	if errors.Is(err, tracer.ErrInvalidProvider) {
		return ErrTracerInvalidProvider
	}
	if errors.Is(err, metric.ErrInvalidProvider) {
		return ErrMetricInvalidProvider
	}
	return fmt.Errorf("%s: %w", message, err)
}

// NewLogger initializes a Logger component with the given options.
//
// It creates a structured logger using Zap that can be used for application logging.
// The logger supports multiple log levels and outputs structured JSON logs.
//
// Optional options:
//   - WithLoggerLevel: Log level (default: "info")
//     Valid values: "debug", "info", "warn", "error", "fatal"
//
// Returns an error if:
//   - Logger initialization fails
//   - Invalid log level is provided
//
// Example:
//
//	logger, err := NewLogger(
//	    WithLoggerLevel("debug"),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer logger.Sync()
func NewLogger(opts ...Option) (Logger, error) {
	options := parseOptions(opts...)
	loggerInstance, err := logger.NewLogger(logger.WithLevel(options.LoggerLevel))
	if err != nil {
		return nil, parseError(err, "failed to initialize logger")
	}
	return loggerInstance, nil
}

// NewTracer initializes a Tracer component with the given options.
//
// It creates a distributed tracing provider using OpenTelemetry that can be used
// to create and manage spans for request tracing. The tracer supports both stdout
// and OTLP exporters for development and production environments.
//
// Required options:
//   - WithServiceName: Service name must be provided
//
// Optional options:
//   - WithEnvironment: Deployment environment (default: "development")
//   - WithInstance: Instance name and host
//   - WithTracerProvider: Tracer exporter configuration (default: "stdout")
//     Valid values: "stdout" (for development) or "otlp" (for production)
//   - WithTracerSampleRatio: Sampling ratio (default: 1.0)
//     Controls what percentage of traces are sampled (0.0 to 1.0)
//   - WithTracerBatchTimeout: Batch timeout (default: 5 seconds)
//     Maximum time to wait before exporting a batch of spans
//   - WithTracerInsecure: Use insecure connection for OTLP (default: false)
//     Should only be used in development or when TLS is handled by a proxy
//
// Returns an error if:
//   - Service name is not provided
//   - Tracer initialization fails
//   - Invalid provider is specified
//
// Example:
//
//	tracer, err := NewTracer(
//	    WithServiceName("my-service"),
//	    WithEnvironment("production"),
//	    WithInstance("instance-1", "localhost"),
//	    WithTracerProvider("otlp", "localhost", 4317),
//	    WithTracerSampleRatio(0.1), // Sample 10% of traces
//	    WithTracerInsecure(false),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer tracer.Shutdown(context.Background())
func NewTracer(opts ...Option) (Tracer, error) {
	options := parseOptions(opts...)
	tracerInstance, err := tracer.NewTracer(
		tracer.WithServiceName(options.ServiceName),
		tracer.WithEnvironment(options.Environment),
		tracer.WithInstance(options.InstanceName, options.InstanceHost),
		tracer.WithProvider(options.TracerProvider, options.TracerProviderHost, options.TracerProviderPort),
		tracer.WithSampleRatio(options.TracerSampleRatio),
		tracer.WithBatchTimeout(options.TracerBatchTimeout),
		tracer.WithInsecure(options.TracerInsecure),
	)
	if err != nil {
		return nil, parseError(err, "failed to initialize tracer")
	}
	return tracerInstance, nil
}

// NewMetric initializes a Metric component with the given options.
//
// It creates a metrics collection provider using OpenTelemetry that can be used
// to record and export application metrics. The metric provider supports both stdout
// and OTLP exporters for development and production environments.
//
// Required options:
//   - WithServiceName: Service name must be provided
//
// Optional options:
//   - WithEnvironment: Deployment environment (default: "development")
//   - WithInstance: Instance name and host
//   - WithMetricProvider: Metric exporter configuration (default: "stdout")
//     Valid values: "stdout" (for development) or "otlp" (for production)
//   - WithMetricInterval: Export interval (default: 60 seconds)
//     Time interval between metric exports to the configured provider
//   - WithMetricInsecure: Use insecure connection for OTLP (default: false)
//     Should only be used in development or when TLS is handled by a proxy
//
// Returns an error if:
//   - Service name is not provided
//   - Metric initialization fails
//   - Invalid provider is specified
//
// Example:
//
//	metric, err := NewMetric(
//	    WithServiceName("my-service"),
//	    WithEnvironment("production"),
//	    WithInstance("instance-1", "localhost"),
//	    WithMetricProvider("otlp", "localhost", 4318),
//	    WithMetricInterval(30*time.Second), // Export every 30 seconds
//	    WithMetricInsecure(false),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer metric.Shutdown(context.Background())
func NewMetric(opts ...Option) (Metric, error) {
	options := parseOptions(opts...)
	metricInstance, err := metric.NewMetric(
		metric.WithServiceName(options.ServiceName),
		metric.WithEnvironment(options.Environment),
		metric.WithInstance(options.InstanceName, options.InstanceHost),
		metric.WithProvider(options.MetricProvider, options.MetricProviderHost, options.MetricProviderPort),
		metric.WithInterval(options.MetricInterval),
		metric.WithInsecure(options.MetricInsecure),
	)
	if err != nil {
		return nil, parseError(err, "failed to initialize metric")
	}
	return metricInstance, nil
}

// NewMonitoring initializes all monitoring components (Logger, Tracer, Metric) with the given options.
//
// It creates a unified monitoring setup that can be used throughout the application.
// If initialization of any component fails, all previously initialized components
// are cleaned up before returning an error.
//
// Required options:
//   - WithServiceName: Service name must be provided
//
// Optional options:
//   - WithEnvironment: Deployment environment (default: "development")
//   - WithInstance: Instance name and host
//   - WithLoggerLevel: Log level (default: "info")
//   - WithTracerProvider: Tracer exporter configuration (default: "stdout")
//   - WithTracerSampleRatio: Sampling ratio (default: 1.0)
//   - WithTracerBatchTimeout: Batch timeout (default: 5 seconds)
//   - WithMetricProvider: Metric exporter configuration (default: "stdout")
//   - WithMetricInterval: Export interval (default: 60 seconds)
//
// Returns an error if:
//   - Service name is not provided
//   - Logger initialization fails
//   - Tracer initialization fails
//   - Metric initialization fails
//
// Example:
//
//	mon, err := NewMonitoring(
//	    WithServiceName("my-service"),
//	    WithEnvironment("production"),
//	    WithInstance("instance-1", "localhost"),
//	    WithLoggerLevel("info"),
//	    WithTracerProvider("otlp", "localhost", 4317),
//	    WithMetricProvider("otlp", "localhost", 4318),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer mon.Shutdown(context.Background())
func NewMonitoring(opts ...Option) (*Monitoring, error) {
	options := parseOptions(opts...)

	// Validate required options
	if options.ServiceName == "" {
		return nil, ErrServiceNameRequired
	}

	// Initialize logger
	loggerInstance, err := logger.NewLogger(logger.WithLevel(options.LoggerLevel))
	if err != nil {
		return nil, parseError(err, "failed to initialize logger")
	}

	// Initialize tracer
	tracerInstance, err := tracer.NewTracer(
		tracer.WithServiceName(options.ServiceName),
		tracer.WithEnvironment(options.Environment),
		tracer.WithInstance(options.InstanceName, options.InstanceHost),
		tracer.WithProvider(options.TracerProvider, options.TracerProviderHost, options.TracerProviderPort),
		tracer.WithSampleRatio(options.TracerSampleRatio),
		tracer.WithBatchTimeout(options.TracerBatchTimeout),
		tracer.WithInsecure(options.TracerInsecure),
	)
	if err != nil {
		// Cleanup logger before returning
		if loggerInstance != nil {
			_ = loggerInstance.Sync() // Ignore cleanup errors when returning initialization error
		}
		return nil, parseError(err, "failed to initialize tracer")
	}

	// Initialize metric
	metricInstance, err := metric.NewMetric(
		metric.WithServiceName(options.ServiceName),
		metric.WithEnvironment(options.Environment),
		metric.WithInstance(options.InstanceName, options.InstanceHost),
		metric.WithProvider(options.MetricProvider, options.MetricProviderHost, options.MetricProviderPort),
		metric.WithInterval(options.MetricInterval),
		metric.WithInsecure(options.MetricInsecure),
	)
	if err != nil {
		// Cleanup tracer and logger before returning (in reverse order of initialization)
		if tracerInstance != nil {
			_ = tracerInstance.Shutdown(context.Background()) // Ignore cleanup errors when returning initialization error
		}
		if loggerInstance != nil {
			_ = loggerInstance.Sync() // Ignore cleanup errors when returning initialization error
		}
		return nil, parseError(err, "failed to initialize metric")
	}

	return &Monitoring{
		Logger: loggerInstance,
		Tracer: tracerInstance,
		Metric: metricInstance,
	}, nil
}
