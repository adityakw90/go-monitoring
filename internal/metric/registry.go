package metric

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc/credentials"
)

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
//	    WithServiceName("my-service"),
//	    WithProvider("otlp", "localhost", 4318),
//	    WithInterval(30*time.Second),
// NewMetric creates and returns a Metric configured according to the provided Options.
// It builds an OpenTelemetry MeterProvider backed by a PeriodicReader and an exporter
// selected by the Options.Provider (supported: "stdout", "otlp"), and attaches a Resource
// populated from the service attributes in Options.
//
// Errors returned include:
// - ErrIntervalInvalid when Options.Interval is less than or equal to zero.
// - ErrProviderHostRequired, ErrProviderPortRequired, ErrProviderPortInvalid for missing/invalid OTLP host/port.
// - ErrInvalidProvider when Options.Provider is not supported.
// Other errors wrap failures that occur while creating the resource or the exporter.
func NewMetric(opts ...Option) (Metric, error) {
	options := &Options{
		Provider: "stdout",
		Interval: 60 * time.Second,
	}

	for _, opt := range opts {
		opt(options)
	}

	// validate interval
	if options.Interval <= 0 {
		return nil, ErrIntervalInvalid
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
		if options.ProviderHost == "" {
			return nil, ErrProviderHostRequired
		}
		if options.ProviderPort == 0 {
			return nil, ErrProviderPortRequired
		}
		if options.ProviderPort < 0 {
			return nil, ErrProviderPortInvalid
		}
		otlpOpts := []otlpmetricgrpc.Option{
			otlpmetricgrpc.WithEndpoint(
				fmt.Sprintf("%s:%d", options.ProviderHost, options.ProviderPort),
			),
		}
		if options.Insecure {
			otlpOpts = append(otlpOpts, otlpmetricgrpc.WithInsecure())
		} else {
			otlpOpts = append(otlpOpts, otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, options.ProviderHost)))
		}
		exporter, err = otlpmetricgrpc.New(context.Background(), otlpOpts...)
	default:
		return nil, ErrInvalidProvider
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

	return &metric{
		provider: mp,
		meter:    mp.Meter(options.ServiceName),
	}, nil
}