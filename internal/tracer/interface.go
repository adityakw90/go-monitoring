package tracer

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

type Tracer interface {
	StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span)
	EndSpan(span trace.Span)
	Shutdown(ctx context.Context) error
	StartChildSpan(ctx context.Context, name string, parent trace.Span) (context.Context, trace.Span)
	NewSpanFromContext(ctx context.Context) trace.Span
	ExtractContext(ctx context.Context, md metadata.MD) context.Context
	InjectContext(ctx context.Context) metadata.MD
}
