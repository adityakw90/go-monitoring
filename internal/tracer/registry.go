package tracer

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc/credentials"
)

// NewTracer initializes a new OpenTelemetry tracer with the given options.
//
// It creates a tracer provider with the specified exporter (stdout or OTLP),
// configures sampling based on the sample ratio, and sets up resource attributes
// for service identification.
//
// Default configuration:
//   - Provider: "stdout"
//   - SampleRatio: 1.0 (always sample)
//   - BatchTimeout: 5 seconds
//
// Returns an error if:
//   - The provider type is invalid (not "stdout" or "otlp")
//   - Resource creation fails
//   - Exporter creation fails
//
// Example:
//
//	tracer, err := NewTracer(
//	    withTracerServiceName("my-service"),
//	    withTracerProvider("otlp", "localhost", 4317),
//	    withTracerSampleRatio(0.1),
//	)
func NewTracer(opts ...Option) (Tracer, error) {
	options := &Options{
		Provider:     "stdout",
		SampleRatio:  1.0,
		BatchTimeout: 5 * time.Second,
	}

	for _, opt := range opts {
		opt(options)
	}

	// Create resource with service name
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
	var exporter sdktrace.SpanExporter
	switch options.Provider {
	case "stdout":
		exporter, err = stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
		)
	case "otlp":
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(
				fmt.Sprintf("%s:%d", options.ProviderHost, options.ProviderPort),
			),
		}
		if options.Insecure {
			opts = append(opts, otlptracegrpc.WithInsecure())
		} else {
			opts = append(opts, otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")))
		}
		exporter, err = otlptracegrpc.New(context.Background(), opts...)
	default:
		return nil, ErrInvalidProvider
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	// Create a sampler with the ratio from config
	var sampler sdktrace.Sampler
	switch {
	case options.SampleRatio <= 0:
		sampler = sdktrace.NeverSample()
	case options.SampleRatio >= 1.0:
		sampler = sdktrace.AlwaysSample()
	default:
		sampler = sdktrace.TraceIDRatioBased(options.SampleRatio)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(
			exporter,
			sdktrace.WithBatchTimeout(options.BatchTimeout),
		),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	return &tracer{
		provider:   tp,
		tracer:     tp.Tracer(options.ServiceName),
		propagator: propagation.TraceContext{},
	}, nil
}
