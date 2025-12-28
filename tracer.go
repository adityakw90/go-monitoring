package monitoring

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

// Tracer wraps OpenTelemetry tracer and provides distributed tracing functionality.
// It supports multiple exporters (stdout, OTLP) and configurable sampling.
type Tracer struct {
	provider   *sdktrace.TracerProvider
	tracer     trace.Tracer
	propagator propagation.TextMapPropagator
}

// TracerOptions contains configuration options for creating a Tracer.
// All fields are optional and have sensible defaults.
type TracerOptions struct {
	ServiceName  string        // ServiceName is the name of the service being traced.
	Environment  string        // Environment is the deployment environment (e.g., "development", "production").
	InstanceName string        // InstanceName is the unique identifier for this service instance.
	InstanceHost string        // InstanceHost is the hostname where this service instance is running.
	Provider     string        // Provider specifies the trace exporter to use ("stdout" or "otlp").
	ProviderHost string        // ProviderHost is the hostname of the OTLP trace collector (only used when Provider is "otlp").
	ProviderPort int           // ProviderPort is the port of the OTLP trace collector (only used when Provider is "otlp").
	SampleRatio  float64       // SampleRatio controls the sampling rate for traces (0.0 to 1.0). 0.0 means never sample, 1.0 means always sample, values in between use probabilistic sampling.
	BatchTimeout time.Duration // BatchTimeout is the maximum time to wait before exporting a batch of spans.
	Insecure     bool          // Insecure controls whether to use an insecure (non-TLS) connection for OTLP exporter. When true, connections are made without TLS. Default is false (secure TLS connection).
}

// TracerOption is a function that configures TracerOptions.
// It follows the functional options pattern for flexible tracer configuration.
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

// withTracerInsecure sets whether to use an insecure connection for OTLP exporter (internal use).
func withTracerInsecure(insecure bool) TracerOption {
	return func(o *TracerOptions) {
		o.Insecure = insecure
	}
}

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
		}
		exporter, err = otlptracegrpc.New(context.Background(), opts...)
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidProvider, options.Provider)
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

	return &Tracer{
		provider:   tp,
		tracer:     tp.Tracer(options.ServiceName),
		propagator: propagation.TraceContext{},
	}, nil
}

// StartSpan starts a new span with the given name and context.
// It returns a new context containing the span and the span itself.
// The span should be ended by calling EndSpan or span.End().
//
// Parameters:
//   - ctx: The parent context (may contain a parent span)
//   - name: The name of the span (should be descriptive, e.g., "handle-request")
//   - opts: Optional span start options (e.g., trace.WithSpanKind)
//
// Returns:
//   - A new context containing the span
//   - The created span
//
// Example:
//
//	ctx, span := tracer.StartSpan(ctx, "process-payment")
//	defer tracer.EndSpan(span)
func (t *Tracer) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name, opts...)
}

// EndSpan ends the given span, recording its completion time.
// This should be called when the operation represented by the span is complete.
// Typically used with defer to ensure spans are always ended.
//
// Example:
//
//	ctx, span := tracer.StartSpan(ctx, "operation")
//	defer tracer.EndSpan(span)
func (t *Tracer) EndSpan(span trace.Span) {
	span.End()
}

// Shutdown gracefully shuts down the tracer provider.
// It flushes any pending spans and releases resources.
// This should be called before application shutdown to ensure all traces are exported.
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
//	if err := tracer.Shutdown(ctx); err != nil {
//	    log.Printf("Failed to shutdown tracer: %v", err)
//	}
func (t *Tracer) Shutdown(ctx context.Context) error {
	return t.provider.Shutdown(ctx)
}

// NewSpanFromSpan creates a new child span from a parent span.
// The new span will be linked to the parent span's trace context.
//
// Parameters:
//   - ctx: The context to use for the new span
//   - name: The name of the new span
//   - parent: The parent span to create a child from
//
// Returns:
//   - A new context containing the child span
//   - The created child span
//
// Example:
//
//	ctx, parentSpan := tracer.StartSpan(ctx, "parent-operation")
//	defer tracer.EndSpan(parentSpan)
//
//	ctx, childSpan := tracer.NewSpanFromSpan(ctx, "child-operation", parentSpan)
//	defer tracer.EndSpan(childSpan)
func (t *Tracer) NewSpanFromSpan(ctx context.Context, name string, parent trace.Span) (context.Context, trace.Span) {
	newCtx := trace.ContextWithSpanContext(ctx, parent.SpanContext())
	return t.StartSpan(newCtx, name)
}

// NewSpanFromContext retrieves the span from the given context.
// Returns nil if no span is present in the context.
//
// Parameters:
//   - ctx: The context that may contain a span
//
// Returns:
//   - The span from the context, or nil if no span exists
//
// Example:
//
//	span := tracer.NewSpanFromContext(ctx)
//	if span != nil {
//	    span.SetAttributes(attribute.String("key", "value"))
//	}
func (t *Tracer) NewSpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// ExtractContext extracts trace context from gRPC metadata.
// This is used on the server side to extract trace context from incoming gRPC requests.
// The extracted context can be used to continue the trace across service boundaries.
//
// Parameters:
//   - ctx: The base context
//   - md: gRPC metadata containing trace propagation headers
//
// Returns:
//   - A new context containing the extracted trace context
//
// Example:
//
//	// In gRPC server handler
//	ctx := tracer.ExtractContext(ctx, md)
//	ctx, span := tracer.StartSpan(ctx, "handle-request")
//	defer tracer.EndSpan(span)
func (t *Tracer) ExtractContext(ctx context.Context, md metadata.MD) context.Context {
	carrier := propagation.HeaderCarrier{}
	for k, v := range md {
		if len(v) > 0 {
			carrier.Set(k, v[0])
		}
	}
	return t.propagator.Extract(ctx, carrier)
}

// InjectContext injects trace context into gRPC metadata.
// This is used on the client side to propagate trace context to downstream services.
// The returned metadata should be attached to outgoing gRPC requests.
//
// Parameters:
//   - ctx: The context containing the trace context to inject
//
// Returns:
//   - gRPC metadata with trace propagation headers (keys are lowercase)
//
// Example:
//
//	// In gRPC client
//	md := tracer.InjectContext(ctx)
//	ctx := metadata.NewOutgoingContext(ctx, md)
//	resp, err := client.Call(ctx, req)
func (t *Tracer) InjectContext(ctx context.Context) metadata.MD {
	md := metadata.New(nil)
	t.propagator.Inject(ctx, propagation.HeaderCarrier(md))

	// Force metadata keys to lowercase
	mdLower := metadata.New(nil)
	for k, v := range md {
		mdLower[strings.ToLower(k)] = v
	}

	return mdLower
}
