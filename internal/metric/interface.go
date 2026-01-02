package metric

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
)

type Metric interface {
	CreateCounter(name, unit, description string) (otelmetric.Int64Counter, error)
	RecordCounter(ctx context.Context, counter otelmetric.Int64Counter, value int64, labels ...attribute.KeyValue)
	CreateHistogram(name, unit, description string) (otelmetric.Int64Histogram, error)
	RecordHistogram(ctx context.Context, histogram otelmetric.Int64Histogram, value int64, labels ...attribute.KeyValue)
	CreateAttributeInt(key string, value int) attribute.KeyValue
	CreateAttributeString(key string, value string) attribute.KeyValue
	Shutdown(ctx context.Context) error
}
