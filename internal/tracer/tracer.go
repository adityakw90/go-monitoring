package tracer

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

// tracer wraps OpenTelemetry tracer and provides distributed tracing functionality.
// It supports multiple exporters (stdout, OTLP) and configurable sampling.
type tracer struct {
	provider   *sdktrace.TracerProvider
	tracer     trace.Tracer
	propagator propagation.TextMapPropagator
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
func (t *tracer) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
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
func (t *tracer) EndSpan(span trace.Span) {
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
func (t *tracer) Shutdown(ctx context.Context) error {
	return t.provider.Shutdown(ctx)
}

// StartChildSpan creates a new child span from a parent span.
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
//	ctx, childSpan := tracer.StartChildSpan(ctx, "child-operation", parentSpan)
//	defer tracer.EndSpan(childSpan)
func (t *tracer) StartChildSpan(ctx context.Context, name string, parent trace.Span) (context.Context, trace.Span) {
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
func (t *tracer) NewSpanFromContext(ctx context.Context) trace.Span {
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
func (t *tracer) ExtractContext(ctx context.Context, md metadata.MD) context.Context {
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
func (t *tracer) InjectContext(ctx context.Context) metadata.MD {
	md := metadata.New(nil)
	t.propagator.Inject(ctx, propagation.HeaderCarrier(md))

	// Force metadata keys to lowercase
	mdLower := metadata.New(nil)
	for k, v := range md {
		mdLower[strings.ToLower(k)] = v
	}

	return mdLower
}
