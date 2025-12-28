package monitoring

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc/credentials"
)

// Metric wraps OpenTelemetry meter and provides metrics collection functionality.
// It supports counters and histograms with configurable exporters (stdout, OTLP).
type Metric struct {
	provider *sdkmetric.MeterProvider
	meter    metric.Meter
}

// MetricOptions contains configuration options for creating a Metric.
// All fields are optional and have sensible defaults.
type MetricOptions struct {
	ServiceName  string        // ServiceName is the name of the service collecting metrics.
	Environment  string        // Environment is the deployment environment (e.g., "development", "production").
	InstanceName string        // InstanceName is the unique identifier for this service instance.
	InstanceHost string        // InstanceHost is the hostname where this service instance is running.
	Provider     string        // Provider specifies the metric exporter to use ("stdout" or "otlp").
	ProviderHost string        // ProviderHost is the hostname of the OTLP metric collector (only used when Provider is "otlp").
	ProviderPort int           // ProviderPort is the port of the OTLP metric collector (only used when Provider is "otlp").
	Interval     time.Duration // Interval is the time interval between metric exports.
	Insecure     bool          // Insecure controls whether to use an insecure (non-TLS) connection for OTLP exporter. When true, connections are made without TLS. Default is false (secure TLS connection).
}

// MetricOption is a function that configures MetricOptions.
// It follows the functional options pattern for flexible metric configuration.
type MetricOption func(*MetricOptions)

// withMetricServiceName sets the service name for the metric (internal use).
func withMetricServiceName(name string) MetricOption {
	return func(o *MetricOptions) {
		o.ServiceName = name
	}
}

// withMetricEnvironment sets the environment (internal use).
func withMetricEnvironment(env string) MetricOption {
	return func(o *MetricOptions) {
		o.Environment = env
	}
}

// withMetricInstance sets the instance name and host (internal use).
func withMetricInstance(name, host string) MetricOption {
	return func(o *MetricOptions) {
		o.InstanceName = name
		o.InstanceHost = host
	}
}

// withMetricProvider sets the metric provider configuration (internal use).
func withMetricProvider(provider, host string, port int) MetricOption {
	return func(o *MetricOptions) {
		o.Provider = provider
		o.ProviderHost = host
		o.ProviderPort = port
	}
}

// withMetricInterval sets the export interval (internal use).
func withMetricInterval(interval time.Duration) MetricOption {
	return func(o *MetricOptions) {
		o.Interval = interval
	}
}

// withMetricInsecure sets whether to use an insecure connection for OTLP exporter (internal use).
func withMetricInsecure(insecure bool) MetricOption {
	return func(o *MetricOptions) {
		o.Insecure = insecure
	}
}

// NewMetric initializes a new OpenTelemetry metric with the given options.
//
// It creates a meter provider with the specified exporter (stdout or OTLP),
// configures periodic metric export, and sets up resource attributes
// for service identification.
//
// Default configuration:
//   - Provider: "stdout"
//   - Interval: 60 seconds
//
// Returns an error if:
//   - The provider type is invalid (not "stdout" or "otlp")
//   - Resource creation fails
//   - Exporter creation fails
//
// Example:
//
//	metric, err := NewMetric(
//	    withMetricServiceName("my-service"),
//	    withMetricProvider("otlp", "localhost", 4318),
//	    withMetricInterval(30*time.Second),
//	)
func NewMetric(opts ...MetricOption) (*Metric, error) {
	options := &MetricOptions{
		Provider: "stdout",
		Interval: 60 * time.Second,
	}

	for _, opt := range opts {
		opt(options)
	}

	// Create resource with service name and other attributes
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceInstanceIDKey.String(options.InstanceName),
			semconv.HostNameKey.String(options.InstanceHost),
			semconv.DeploymentEnvironmentKey.String(options.Environment),
			semconv.ServiceNameKey.String(options.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Select the exporter based on the config
	var exporter sdkmetric.Exporter
	switch options.Provider {
	case "stdout":
		exporter, err = stdoutmetric.New(
			stdoutmetric.WithPrettyPrint(),
		)
	case "otlp":
		opts := []otlpmetricgrpc.Option{
			otlpmetricgrpc.WithEndpoint(
				fmt.Sprintf("%s:%d", options.ProviderHost, options.ProviderPort),
			),
		}
		if options.Insecure {
			opts = append(opts, otlpmetricgrpc.WithInsecure())
		} else {
			opts = append(opts, otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")))
		}
		exporter, err = otlpmetricgrpc.New(context.Background(), opts...)
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidProvider, options.Provider)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	// Create the MeterProvider with the exporter
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				exporter,
				sdkmetric.WithInterval(options.Interval),
			),
		),
	)

	return &Metric{
		provider: mp,
		meter:    mp.Meter(options.ServiceName),
	}, nil
}

// CreateCounter creates a new counter metric.
// Counters are monotonically increasing metrics that track cumulative values.
//
// Parameters:
//   - name: The metric name (should follow OpenTelemetry naming conventions)
//   - unit: The unit of measurement (e.g., "1", "ms", "bytes")
//   - description: A human-readable description of what the counter measures
//
// Returns:
//   - The created counter metric
//   - An error if counter creation fails
//
// Example:
//
//	counter, err := metric.CreateCounter(
//	    "http_requests_total",
//	    "1",
//	    "Total number of HTTP requests",
//	)
func (m *Metric) CreateCounter(name, unit, description string) (metric.Int64Counter, error) {
	counter, err := m.meter.Int64Counter(
		name,
		metric.WithDescription(description),
		metric.WithUnit(unit),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create counter: %w", err)
	}
	return counter, nil
}

// RecordCounter increments a counter by a given value.
// The counter must have been created using CreateCounter.
//
// Parameters:
//   - ctx: Context for the metric recording
//   - counter: The counter metric to increment
//   - value: The value to add to the counter (must be non-negative)
//   - labels: Optional key-value pairs for metric dimensions
//
// Example:
//
//	metric.RecordCounter(ctx, counter, 1,
//	    metric.CreateAttributeString("method", "GET"),
//	    metric.CreateAttributeString("status", "200"),
//	)
func (m *Metric) RecordCounter(ctx context.Context, counter metric.Int64Counter, value int64, labels ...attribute.KeyValue) {
	counter.Add(ctx, value, metric.WithAttributes(labels...))
}

// CreateHistogram creates a new histogram metric.
// Histograms track the distribution of values over time.
//
// Parameters:
//   - name: The metric name (should follow OpenTelemetry naming conventions)
//   - unit: The unit of measurement (e.g., "ms", "bytes", "seconds")
//   - description: A human-readable description of what the histogram measures
//
// Returns:
//   - The created histogram metric
//   - An error if histogram creation fails
//
// Example:
//
//	histogram, err := metric.CreateHistogram(
//	    "http_request_duration_ms",
//	    "ms",
//	    "HTTP request duration in milliseconds",
//	)
func (m *Metric) CreateHistogram(name, unit, description string) (metric.Int64Histogram, error) {
	histogram, err := m.meter.Int64Histogram(
		name,
		metric.WithDescription(description),
		metric.WithUnit(unit),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create histogram: %w", err)
	}
	return histogram, nil
}

// RecordHistogram records a value in a histogram.
// The histogram must have been created using CreateHistogram.
//
// Parameters:
//   - ctx: Context for the metric recording
//   - histogram: The histogram metric to record to
//   - value: The value to record (e.g., request duration, response size)
//   - labels: Optional key-value pairs for metric dimensions
//
// Example:
//
//	start := time.Now()
//	// ... perform operation ...
//	duration := time.Since(start).Milliseconds()
//	metric.RecordHistogram(ctx, histogram, duration,
//	    metric.CreateAttributeString("endpoint", "/api/users"),
//	)
func (m *Metric) RecordHistogram(ctx context.Context, histogram metric.Int64Histogram, value int64, labels ...attribute.KeyValue) {
	histogram.Record(ctx, value, metric.WithAttributes(labels...))
}

// CreateAttributeInt creates an integer attribute for metric labels.
// Attributes are used to add dimensions to metrics for filtering and aggregation.
//
// Parameters:
//   - key: The attribute key (should follow OpenTelemetry naming conventions)
//   - value: The integer value
//
// Returns:
//   - An attribute key-value pair
//
// Example:
//
//	attr := metric.CreateAttributeInt("status_code", 200)
//	metric.RecordCounter(ctx, counter, 1, attr)
func (m *Metric) CreateAttributeInt(key string, value int) attribute.KeyValue {
	return attribute.Int(key, value)
}

// CreateAttributeString creates a string attribute for metric labels.
// Attributes are used to add dimensions to metrics for filtering and aggregation.
//
// Parameters:
//   - key: The attribute key (should follow OpenTelemetry naming conventions)
//   - value: The string value
//
// Returns:
//   - An attribute key-value pair
//
// Example:
//
//	attr := metric.CreateAttributeString("method", "GET")
//	metric.RecordCounter(ctx, counter, 1, attr)
func (m *Metric) CreateAttributeString(key string, value string) attribute.KeyValue {
	return attribute.String(key, value)
}

// Shutdown gracefully shuts down the meter provider.
// It flushes any pending metrics and releases resources.
// This should be called before application shutdown to ensure all metrics are exported.
//
// Parameters:
//   - ctx: Context for controlling shutdown timeout
//
// Returns an error if shutdown fails or times out.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	if err := metric.Shutdown(ctx); err != nil {
//	    log.Printf("Failed to shutdown metric: %v", err)
//	}
func (m *Metric) Shutdown(ctx context.Context) error {
	return m.provider.Shutdown(ctx)
}
