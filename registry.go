package monitoring

import (
	"context"

	"github.com/adityakw90/go-monitoring/internal/logger"
	"github.com/adityakw90/go-monitoring/internal/metric"
	"github.com/adityakw90/go-monitoring/internal/tracer"
)

// parseOptions applies the provided functional options to a copy of the package default Options
// and returns the resulting configuration. Options are applied in order; later options override earlier ones.
func parseOptions(opts ...Option) *Options {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// NewLogger creates a Logger configured by the provided functional options.
// It returns the initialized Logger or an error if initialization fails.
func NewLogger(opts ...Option) (Logger, error) {
	options := parseOptions(opts...)
	loggerInstance, err := logger.NewLogger(
		logger.WithLevel(options.LoggerLevel),
		logger.WithOutputPath(options.LoggerOutputPath),
	)
	if err != nil {
		return nil, parseError(err, "failed to initialize logger")
	}
	return loggerInstance, nil
}

// NewTracer creates and returns a Tracer configured by the provided options.
// It applies functional options to the default configuration and initializes an
// underlying tracer instance with service name, environment, instance info,
// provider settings, sampling ratio, batch timeout, and insecure flag.
// Returns a non-nil error if tracer initialization fails.
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

// NewMetric creates a Metric configured by the provided functional options.
// It applies the options to defaults and initializes the metric backend accordingly.
// On success it returns the initialized Metric. If initialization fails it returns
// nil and an error describing the failure (prefixed with "failed to initialize metric").
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

// NewMonitoring initializes and returns a Monitoring containing Logger, Tracer, and Metric configured by the provided options.
// It requires the ServiceName option; when ServiceName is empty it returns ErrServiceNameRequired.
// If initialization of any component fails, previously initialized components are cleaned up (logger Sync, tracer Shutdown) and the error is returned wrapped via parseError.
func NewMonitoring(opts ...Option) (*Monitoring, error) {
	options := parseOptions(opts...)

	// Validate required options
	if options.ServiceName == "" {
		return nil, ErrServiceNameRequired
	}

	// Initialize logger
	loggerInstance, err := logger.NewLogger(
		logger.WithLevel(options.LoggerLevel),
		logger.WithOutputPath(options.LoggerOutputPath),
	)
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
