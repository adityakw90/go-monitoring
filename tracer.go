package monitoring

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

// Tracer wraps OpenTelemetry tracer.
type Tracer struct {
	provider *sdktrace.TracerProvider
	tracer   trace.Tracer
}

// TracerOptions contains tracer configuration.
type TracerOptions struct {
	ServiceName  string
	Environment  string
	InstanceName string
	InstanceHost string
	Provider     string
	ProviderHost string
	ProviderPort int
	SampleRatio  float64
	BatchTimeout time.Duration
}

// TracerOption configures TracerOptions.
type TracerOption func(*TracerOptions)

// withTracerServiceName sets the service name for the tracer (internal use).
func withTracerServiceName(name string) TracerOption {
	return func(o *TracerOptions) {
		o.ServiceName = name
	}
}

// withTracerEnvironment sets the environment (internal use).
func withTracerEnvironment(env string) TracerOption {
	return func(o *TracerOptions) {
		o.Environment = env
	}
}

// withTracerInstance sets the instance name and host (internal use).
func withTracerInstance(name, host string) TracerOption {
	return func(o *TracerOptions) {
		o.InstanceName = name
		o.InstanceHost = host
	}
}

// withTracerProvider sets the tracer provider configuration (internal use).
func withTracerProvider(provider, host string, port int) TracerOption {
	return func(o *TracerOptions) {
		o.Provider = provider
		o.ProviderHost = host
		o.ProviderPort = port
	}
}

// withTracerSampleRatio sets the sampling ratio (internal use).
func withTracerSampleRatio(ratio float64) TracerOption {
	return func(o *TracerOptions) {
		o.SampleRatio = ratio
	}
}

// withTracerBatchTimeout sets the batch timeout (internal use).
func withTracerBatchTimeout(timeout time.Duration) TracerOption {
	return func(o *TracerOptions) {
		o.BatchTimeout = timeout
	}
}

// NewTracer initializes a new OpenTelemetry tracer with the given options.
func NewTracer(opts ...TracerOption) (*Tracer, error) {
	options := &TracerOptions{
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
	var exporter sdktrace.SpanExporter
	switch options.Provider {
	case "stdout":
		exporter, err = stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
		)
	case "otlp":
		exporter, err = otlptracegrpc.New(
			context.Background(),
			otlptracegrpc.WithEndpoint(
				fmt.Sprintf("%s:%d", options.ProviderHost, options.ProviderPort),
			),
			otlptracegrpc.WithInsecure(),
		)
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidProvider, options.Provider)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	// Create a sampler with the ratio from config
	var sampler sdktrace.Sampler
	if options.SampleRatio > 0 && options.SampleRatio <= 1.0 {
		sampler = sdktrace.TraceIDRatioBased(options.SampleRatio)
	} else {
		sampler = sdktrace.AlwaysSample()
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(
			exporter,
			sdktrace.WithBatchTimeout(options.BatchTimeout),
		),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	// Register W3C Trace Context Propagator
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Set the global tracer provider
	otel.SetTracerProvider(tp)

	return &Tracer{
		provider: tp,
		tracer:   otel.Tracer(options.ServiceName),
	}, nil
}

// StartSpan starts a span with the given name and context.
func (t *Tracer) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name, opts...)
}

// EndSpan ends the given span.
func (t *Tracer) EndSpan(span trace.Span) {
	span.End()
}

// Shutdown gracefully shuts down the tracer provider.
func (t *Tracer) Shutdown(ctx context.Context) error {
	return t.provider.Shutdown(ctx)
}

// NewSpanFromSpan creates a new span from a parent span.
func (t *Tracer) NewSpanFromSpan(ctx context.Context, name string, parent trace.Span) (context.Context, trace.Span) {
	newCtx := trace.ContextWithSpanContext(ctx, parent.SpanContext())
	return t.StartSpan(newCtx, name, trace.WithLinks(trace.Link{
		SpanContext: parent.SpanContext(),
	}))
}

// NewSpanFromContext gets a span from context.
func (t *Tracer) NewSpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// ExtractContext extracts trace context from gRPC metadata.
func (t *Tracer) ExtractContext(ctx context.Context, md metadata.MD) context.Context {
	carrier := propagation.HeaderCarrier{}
	for k, v := range md {
		if len(v) > 0 {
			carrier.Set(k, v[0])
		}
	}
	return otel.GetTextMapPropagator().Extract(ctx, carrier)
}

// InjectContext injects trace context into gRPC metadata.
func (t *Tracer) InjectContext(ctx context.Context) metadata.MD {
	md := metadata.New(nil)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(md))

	// Force metadata keys to lowercase
	mdLower := metadata.New(nil)
	for k, v := range md {
		mdLower[strings.ToLower(k)] = v
	}

	return mdLower
}
