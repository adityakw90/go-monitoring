package monitoring

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// Metric wraps OpenTelemetry meter.
type Metric struct {
	provider *sdkmetric.MeterProvider
	meter    metric.Meter
}

// MetricOptions contains metric configuration.
type MetricOptions struct {
	ServiceName  string
	Environment  string
	InstanceName string
	InstanceHost string
	Provider     string
	ProviderHost string
	ProviderPort int
	Interval     time.Duration
}

// MetricOption configures MetricOptions.
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

// NewMetric initializes a new OpenTelemetry metric with the given options.
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
			attribute.String("instance_name", options.InstanceName),
			attribute.String("instance_host", options.InstanceHost),
			attribute.String("service", options.ServiceName),
			attribute.String("environment", options.Environment),
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
		exporter, err = otlpmetricgrpc.New(
			context.Background(),
			otlpmetricgrpc.WithEndpoint(
				fmt.Sprintf("%s:%d", options.ProviderHost, options.ProviderPort),
			),
			otlpmetricgrpc.WithInsecure(),
		)
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

	// Set the global meter provider
	otel.SetMeterProvider(mp)

	return &Metric{
		provider: mp,
		meter:    otel.Meter(options.ServiceName),
	}, nil
}

// CreateCounter creates a new counter metric.
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
func (m *Metric) RecordCounter(ctx context.Context, counter metric.Int64Counter, value int64, labels ...attribute.KeyValue) {
	counter.Add(ctx, value, metric.WithAttributes(labels...))
}

// CreateHistogram creates a new histogram metric.
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
func (m *Metric) RecordHistogram(ctx context.Context, histogram metric.Int64Histogram, value int64, labels ...attribute.KeyValue) {
	histogram.Record(ctx, value, metric.WithAttributes(labels...))
}

// CreateAttributeInt creates an integer attribute.
func (m *Metric) CreateAttributeInt(key string, value int) attribute.KeyValue {
	return attribute.Int(key, value)
}

// CreateAttributeString creates a string attribute.
func (m *Metric) CreateAttributeString(key string, value string) attribute.KeyValue {
	return attribute.String(key, value)
}

// Shutdown gracefully shuts down the meter provider.
func (m *Metric) Shutdown(ctx context.Context) error {
	return m.provider.Shutdown(ctx)
}
