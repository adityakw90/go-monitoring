package monitoring

import (
	"context"
	"strings"
	"testing"
	"time"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

func TestNewTracer(t *testing.T) {
	tests := []struct {
		name    string
		opts    []TracerOption
		wantErr bool
	}{
		{
			name:    "default tracer with stdout",
			opts:    []TracerOption{withTracerServiceName("test-service")},
			wantErr: false,
		},
		{
			name: "with all options",
			opts: []TracerOption{
				withTracerServiceName("test-service"),
				withTracerEnvironment("test"),
				withTracerInstance("instance-1", "localhost"),
				withTracerProvider("stdout", "", 0),
				withTracerSampleRatio(0.5),
				withTracerBatchTimeout(10 * time.Second),
			},
			wantErr: false,
		},
		{
			name: "with otlp provider (insecure)",
			opts: []TracerOption{
				withTracerServiceName("test-service"),
				withTracerProvider("otlp", "localhost", 4317),
				withTracerInsecure(true),
			},
			wantErr: false,
		},
		{
			name: "with otlp provider (secure)",
			opts: []TracerOption{
				withTracerServiceName("test-service"),
				withTracerProvider("otlp", "localhost", 4317),
				withTracerInsecure(false),
			},
			wantErr: false,
		},
		{
			name: "with invalid provider",
			opts: []TracerOption{
				withTracerServiceName("test-service"),
				withTracerProvider("invalid", "", 0),
			},
			wantErr: true,
		},
		{
			name: "with sample ratio 0",
			opts: []TracerOption{
				withTracerServiceName("test-service"),
				withTracerSampleRatio(0.0),
			},
			wantErr: false,
		},
		{
			name: "with sample ratio 1.0",
			opts: []TracerOption{
				withTracerServiceName("test-service"),
				withTracerSampleRatio(1.0),
			},
			wantErr: false,
		},
		{
			name: "with invalid sample ratio > 1.0",
			opts: []TracerOption{
				withTracerServiceName("test-service"),
				withTracerSampleRatio(1.5),
			},
			wantErr: false, // Should default to AlwaysSample
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracer, err := NewTracer(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTracer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tracer == nil {
					t.Errorf("NewTracer() returned nil")
					return
				}
				if tracer.provider == nil {
					t.Errorf("NewTracer() provider is nil")
				}
				if tracer.tracer == nil {
					t.Errorf("NewTracer() tracer is nil")
				}
				if tracer.propagator == nil {
					t.Errorf("NewTracer() propagator is nil")
				}
				// Cleanup
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_ = tracer.Shutdown(ctx)
			}
		})
	}
}

func TestTracer_StartSpan(t *testing.T) {
	tracer, err := NewTracer(withTracerServiceName("test-service"))
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

func TestTracer_EndSpan(t *testing.T) {
	tracer, err := NewTracer(withTracerServiceName("test-service"))
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

func TestTracer_Shutdown(t *testing.T) {
	tracer, err := NewTracer(withTracerServiceName("test-service"))
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

func TestTracer_NewSpanFromSpan(t *testing.T) {
	tracer, err := NewTracer(withTracerServiceName("test-service"))
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

	ctx2, childSpan := tracer.NewSpanFromSpan(ctx, "child-operation", parentSpan)
	if childSpan == nil {
		t.Errorf("NewSpanFromSpan() returned nil span")
	}
	if !childSpan.SpanContext().IsValid() {
		t.Errorf("NewSpanFromSpan() returned invalid span context")
	}

	// Verify parent-child relationship
	parentCtx := parentSpan.SpanContext()
	childCtx := childSpan.SpanContext()

	// Child should have the same TraceID as parent
	if childCtx.TraceID() != parentCtx.TraceID() {
		t.Errorf("NewSpanFromSpan() child TraceID = %s, want %s", childCtx.TraceID().String(), parentCtx.TraceID().String())
	}

	// Verify the parent span context is correctly propagated in the returned context
	retrievedSpan := trace.SpanFromContext(ctx2)
	if retrievedSpan == nil {
		t.Errorf("NewSpanFromSpan() context does not contain span")
	}
	retrievedCtx := retrievedSpan.SpanContext()
	if retrievedCtx.TraceID() != parentCtx.TraceID() {
		t.Errorf("NewSpanFromSpan() retrieved span TraceID = %s, want %s", retrievedCtx.TraceID().String(), parentCtx.TraceID().String())
	}

	childSpan.End()
}

func TestTracer_NewSpanFromContext(t *testing.T) {
	tracer, err := NewTracer(withTracerServiceName("test-service"))
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

func TestTracer_ExtractContext(t *testing.T) {
	tracer, err := NewTracer(withTracerServiceName("test-service"))
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

func TestTracer_InjectContext(t *testing.T) {
	tracer, err := NewTracer(withTracerServiceName("test-service"))
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

func TestTracer_ExtractContext_EmptyMetadata(t *testing.T) {
	tracer, err := NewTracer(withTracerServiceName("test-service"))
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

func TestTracer_ExtractContext_WithMultipleValues(t *testing.T) {
	tracer, err := NewTracer(withTracerServiceName("test-service"))
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

func TestTracer_MultipleTracersCoexist(t *testing.T) {
	// Create multiple tracers with different configurations
	tracer1, err := NewTracer(
		withTracerServiceName("service-1"),
		withTracerEnvironment("env-1"),
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
		withTracerServiceName("service-2"),
		withTracerEnvironment("env-2"),
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
	if tracer1.provider == nil {
		t.Errorf("tracer1.provider is nil")
	}
	if tracer1.propagator == nil {
		t.Errorf("tracer1.propagator is nil")
	}
	if tracer2.provider == nil {
		t.Errorf("tracer2.provider is nil")
	}
	if tracer2.propagator == nil {
		t.Errorf("tracer2.propagator is nil")
	}

	// Verify they are different instances
	if tracer1.provider == tracer2.provider {
		t.Errorf("tracer1 and tracer2 share the same provider instance")
	}

	// Test that both tracers can create spans independently
	ctx1 := context.Background()
	ctx1, span1 := tracer1.StartSpan(ctx1, "span-1")
	if span1 == nil {
		t.Errorf("tracer1.StartSpan() returned nil")
	}
	span1.End()

	ctx2 := context.Background()
	ctx2, span2 := tracer2.StartSpan(ctx2, "span-2")
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
	md1 := tracer1.InjectContext(ctx1)
	md2 := tracer2.InjectContext(ctx2)

	if len(md1) == 0 {
		t.Errorf("tracer1.InjectContext() returned empty metadata")
	}
	if len(md2) == 0 {
		t.Errorf("tracer2.InjectContext() returned empty metadata")
	}
}
