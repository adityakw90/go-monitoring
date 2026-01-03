package tracer

import (
	"context"
	"strings"
	"testing"
	"time"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

func TestTracer_Tracer_StartSpan(t *testing.T) {
	tracer, err := NewTracer(WithServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewTracer() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = tracer.Shutdown(ctx)
	}()

	ctx := context.Background()
	ctx, span := tracer.StartSpan(ctx, "test-operation")
	if span == nil {
		t.Errorf("StartSpan() returned nil span")
	}
	if !span.SpanContext().IsValid() {
		t.Errorf("StartSpan() returned invalid span context")
	}

	// Test with options
	_, span2 := tracer.StartSpan(ctx, "test-operation-2", trace.WithSpanKind(trace.SpanKindServer))
	if span2 == nil {
		t.Errorf("StartSpan() with options returned nil span")
	}
	span2.End()
}

func TestTracer_Tracer_EndSpan(t *testing.T) {
	tracer, err := NewTracer(WithServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewTracer() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = tracer.Shutdown(ctx)
	}()

	ctx := context.Background()
	_, span := tracer.StartSpan(ctx, "test-operation")

	// EndSpan should not panic
	tracer.EndSpan(span)
}

func TestTracer_Tracer_Shutdown(t *testing.T) {
	tracer, err := NewTracer(WithServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewTracer() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := tracer.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}

	// Shutdown should be idempotent
	if err := tracer.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown() second call error = %v", err)
	}
}

func TestTracer_Tracer_StartChildSpan(t *testing.T) {
	tracer, err := NewTracer(WithServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewTracer() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = tracer.Shutdown(ctx)
	}()

	ctx := context.Background()
	ctx, parentSpan := tracer.StartSpan(ctx, "parent-operation")
	defer parentSpan.End()

	ctx2, childSpan := tracer.StartChildSpan(ctx, "child-operation", parentSpan)
	if childSpan == nil {
		t.Errorf("StartChildSpan() returned nil span")
	}
	if !childSpan.SpanContext().IsValid() {
		t.Errorf("StartChildSpan() returned invalid span context")
	}

	// Verify parent-child relationship
	parentCtx := parentSpan.SpanContext()
	childCtx := childSpan.SpanContext()

	// Child should have the same TraceID as parent
	if childCtx.TraceID() != parentCtx.TraceID() {
		t.Errorf("StartChildSpan() child TraceID = %s, want %s", childCtx.TraceID().String(), parentCtx.TraceID().String())
	}

	// Verify the parent span context is correctly propagated in the returned context
	retrievedSpan := trace.SpanFromContext(ctx2)
	if retrievedSpan == nil {
		t.Errorf("StartChildSpan() context does not contain span")
	}
	retrievedCtx := retrievedSpan.SpanContext()
	if retrievedCtx.TraceID() != parentCtx.TraceID() {
		t.Errorf("StartChildSpan() retrieved span TraceID = %s, want %s", retrievedCtx.TraceID().String(), parentCtx.TraceID().String())
	}

	childSpan.End()
}

func TestTracer_Tracer_NewSpanFromContext(t *testing.T) {
	tracer, err := NewTracer(WithServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewTracer() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = tracer.Shutdown(ctx)
	}()

	ctx := context.Background()
	ctx, span := tracer.StartSpan(ctx, "test-operation")
	defer span.End()

	// Get span from context
	retrievedSpan := tracer.NewSpanFromContext(ctx)
	if retrievedSpan == nil {
		t.Errorf("NewSpanFromContext() returned nil span")
	}
	if retrievedSpan.SpanContext().TraceID() != span.SpanContext().TraceID() {
		t.Errorf("NewSpanFromContext() returned different trace ID")
	}
}

func TestTracer_Tracer_ExtractContext(t *testing.T) {
	tracer, err := NewTracer(WithServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewTracer() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = tracer.Shutdown(ctx)
	}()

	ctx := context.Background()
	ctx, span := tracer.StartSpan(ctx, "test-operation")
	defer span.End()

	// Inject context into metadata
	md := tracer.InjectContext(ctx)
	if len(md) == 0 {
		t.Errorf("InjectContext() returned empty metadata")
	}

	// Extract context from metadata
	ctx2 := context.Background()
	ctx2 = tracer.ExtractContext(ctx2, md)

	// Verify span context was extracted
	span2 := trace.SpanFromContext(ctx2)
	if !span2.SpanContext().IsValid() {
		t.Errorf("ExtractContext() did not extract valid span context")
	}
	if span2.SpanContext().TraceID() != span.SpanContext().TraceID() {
		t.Errorf("ExtractContext() extracted different trace ID")
	}
}

func TestTracer_Tracer_InjectContext(t *testing.T) {
	tracer, err := NewTracer(WithServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewTracer() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = tracer.Shutdown(ctx)
	}()

	ctx := context.Background()
	ctx, span := tracer.StartSpan(ctx, "test-operation")
	defer span.End()

	md := tracer.InjectContext(ctx)
	if len(md) == 0 {
		t.Errorf("InjectContext() returned empty metadata")
	}

	// Verify metadata keys are lowercase
	for k := range md {
		if k != strings.ToLower(k) {
			t.Errorf("InjectContext() metadata key %q is not lowercase", k)
		}
	}

	// Verify trace context is in metadata
	hasTraceContext := false
	for k := range md {
		if k == "traceparent" || k == "tracestate" {
			hasTraceContext = true
			break
		}
	}
	if !hasTraceContext {
		t.Errorf("InjectContext() did not include trace context in metadata")
	}
}

func TestTracer_Tracer_ExtractContext_EmptyMetadata(t *testing.T) {
	tracer, err := NewTracer(WithServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewTracer() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = tracer.Shutdown(ctx)
	}()

	ctx := context.Background()
	emptyMD := metadata.New(nil)
	ctx = tracer.ExtractContext(ctx, emptyMD)

	// Should not panic and return context (may or may not have span)
	_ = ctx
}

func TestTracer_Tracer_ExtractContext_WithMultipleValues(t *testing.T) {
	tracer, err := NewTracer(WithServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewTracer() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = tracer.Shutdown(ctx)
	}()

	ctx := context.Background()
	ctx, span := tracer.StartSpan(ctx, "test-operation")
	defer span.End()

	md := tracer.InjectContext(ctx)

	// Add multiple values to metadata
	md["test-key"] = []string{"value1", "value2"}

	ctx2 := context.Background()
	ctx2 = tracer.ExtractContext(ctx2, md)

	// Should handle multiple values gracefully (takes first value)
	span2 := trace.SpanFromContext(ctx2)
	if !span2.SpanContext().IsValid() {
		t.Errorf("ExtractContext() did not extract valid span context with multiple values")
	}
}

func TestTracer_Tracer_MultipleTracersCoexist(t *testing.T) {
	// Create multiple tracers with different configurations
	tracer1, err := NewTracer(
		WithServiceName("service-1"),
		WithEnvironment("env-1"),
	)
	if err != nil {
		t.Fatalf("NewTracer() for tracer1 error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = tracer1.Shutdown(ctx)
	}()

	tracer2, err := NewTracer(
		WithServiceName("service-2"),
		WithEnvironment("env-2"),
	)
	if err != nil {
		t.Fatalf("NewTracer() for tracer2 error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = tracer2.Shutdown(ctx)
	}()

	// Verify both tracers have their own providers and propagators
	if tracer1.(*tracer).provider == nil {
		t.Errorf("tracer1.provider is nil")
	}
	if tracer1.(*tracer).propagator == nil {
		t.Errorf("tracer1.propagator is nil")
	}
	if tracer2.(*tracer).provider == nil {
		t.Errorf("tracer2.provider is nil")
	}
	if tracer2.(*tracer).propagator == nil {
		t.Errorf("tracer2.propagator is nil")
	}

	// Verify they are different instances
	if tracer1.(*tracer).provider == tracer2.(*tracer).provider {
		t.Errorf("tracer1 and tracer2 share the same provider instance")
	}

	// Test that both tracers can create spans independently
	ctx1 := context.Background()
	ctx1, span1 := tracer1.(*tracer).StartSpan(ctx1, "span-1")
	if span1 == nil {
		t.Errorf("tracer1.StartSpan() returned nil")
	}
	span1.End()

	ctx2 := context.Background()
	ctx2, span2 := tracer2.(*tracer).StartSpan(ctx2, "span-2")
	if span2 == nil {
		t.Errorf("tracer2.StartSpan() returned nil")
	}
	span2.End()

	// Verify spans have valid contexts
	if !span1.SpanContext().IsValid() {
		t.Errorf("span1 has invalid span context")
	}
	if !span2.SpanContext().IsValid() {
		t.Errorf("span2 has invalid span context")
	}

	// Test that each tracer's propagator works independently
	md1 := tracer1.(*tracer).InjectContext(ctx1)
	md2 := tracer2.(*tracer).InjectContext(ctx2)

	if len(md1) == 0 {
		t.Errorf("tracer1.InjectContext() returned empty metadata")
	}
	if len(md2) == 0 {
		t.Errorf("tracer2.(*tracer).InjectContext() returned empty metadata")
	}
}
