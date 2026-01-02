package monitoring

import (
	"context"
	"fmt"

	"github.com/adityakw90/go-monitoring/internal/logger"
	"github.com/adityakw90/go-monitoring/internal/metric"
	"github.com/adityakw90/go-monitoring/internal/tracer"
)

// Monitoring contains all observability components in a single unified structure.
// It provides access to logging, tracing, and metrics functionality.
type Monitoring struct {
	Logger Logger // Logger provides structured logging capabilities.
	Tracer Tracer // Tracer provides distributed tracing capabilities.
	Metric Metric // Metric provides metrics collection capabilities.
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
	loggerInstance, err := logger.NewLogger(logger.WithLevel(options.LoggerLevel))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
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
		return nil, fmt.Errorf("failed to initialize tracer: %w", err)
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
		return nil, fmt.Errorf("failed to initialize metric: %w", err)
	}

	return &Monitoring{
		Logger: loggerInstance,
		Tracer: tracerInstance,
		Metric: metricInstance,
	}, nil
}

// Shutdown gracefully shuts down all monitoring components.
// It shuts down the Tracer and Metric providers in order, ensuring all
// pending traces and metrics are exported before termination.
//
// This should be called before application shutdown to ensure proper cleanup.
// The Logger does not require explicit shutdown.
//
// Parameters:
//   - ctx: Context for controlling shutdown timeout
//
// Returns an error if shutdown of any component fails.
// Errors from individual components are wrapped with context.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	if err := mon.Shutdown(ctx); err != nil {
//	    log.Printf("Failed to shutdown monitoring: %v", err)
//	}
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
