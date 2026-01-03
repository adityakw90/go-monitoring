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
//	)
func NewMetric(opts ...Option) (Metric, error) {
	options := &Options{
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
