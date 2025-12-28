package monitoring

import (
	"context"
	"fmt"
)

// Monitoring contains all observability components.
type Monitoring struct {
	Logger *Logger
	Tracer *Tracer
	Metric *Metric
}

// NewMonitoring initializes all monitoring components with the given options.
func NewMonitoring(opts ...Option) (*Monitoring, error) {
	options := defaultOptions()

	// Apply options
	for _, opt := range opts {
		opt(options)
	}

	// Validate required options
	if options.ServiceName == "" {
		return nil, ErrServiceNameRequired
	}

	// Initialize logger
	logger, err := NewLogger(withLoggerLevel(options.LoggerLevel))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Initialize tracer
	tracer, err := NewTracer(
		withTracerServiceName(options.ServiceName),
		withTracerEnvironment(options.Environment),
		withTracerInstance(options.InstanceName, options.InstanceHost),
		withTracerProvider(options.TracerProvider, options.TracerProviderHost, options.TracerProviderPort),
		withTracerSampleRatio(options.TracerSampleRatio),
		withTracerBatchTimeout(options.TracerBatchTimeout),
		withTracerInsecure(options.TracerInsecure),
	)
	if err != nil {
		// Cleanup logger before returning
		if logger != nil {
			_ = logger.Sync() // Ignore cleanup errors when returning initialization error
		}
		return nil, fmt.Errorf("failed to initialize tracer: %w", err)
	}

	// Initialize metric
	metric, err := NewMetric(
		withMetricServiceName(options.ServiceName),
		withMetricEnvironment(options.Environment),
		withMetricInstance(options.InstanceName, options.InstanceHost),
		withMetricProvider(options.MetricProvider, options.MetricProviderHost, options.MetricProviderPort),
		withMetricInterval(options.MetricInterval),
		withMetricInsecure(options.MetricInsecure),
	)
	if err != nil {
		// Cleanup tracer and logger before returning (in reverse order of initialization)
		if tracer != nil {
			_ = tracer.Shutdown(context.Background()) // Ignore cleanup errors when returning initialization error
		}
		if logger != nil {
			_ = logger.Sync() // Ignore cleanup errors when returning initialization error
		}
		return nil, fmt.Errorf("failed to initialize metric: %w", err)
	}

	return &Monitoring{
		Logger: logger,
		Tracer: tracer,
		Metric: metric,
	}, nil
}

// Shutdown gracefully shuts down all monitoring components.
func (m *Monitoring) Shutdown(ctx context.Context) error {
	if m.Tracer != nil {
		if err := m.Tracer.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown tracer: %w", err)
		}
	}
	if m.Metric != nil {
		if err := m.Metric.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown metric: %w", err)
		}
	}
	return nil
}
