package metric

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// metric wraps OpenTelemetry meter and provides metrics collection functionality.
// It supports counters and histograms with configurable exporters (stdout, OTLP).
type metric struct {
	provider *sdkmetric.MeterProvider
	meter    otelmetric.Meter
}

// CreateCounter creates a new counter metric.
// Counters are monotonically increasing metrics that track cumulative values.
//
// Parameters:
//   - name: The metric name (should follow OpenTelemetry naming conventions)
//   - unit: The unit of measurement (e.g., "1", "ms", "bytes")
//   - description: A human-readable description of what the counter measures
//
// Returns:
//   - The created counter metric
//   - An error if counter creation fails
//
// Example:
//
//	counter, err := metric.CreateCounter(
//	    "http_requests_total",
//	    "1",
//	    "Total number of HTTP requests",
//	)
func (m *metric) CreateCounter(name, unit, description string) (otelmetric.Int64Counter, error) {
	counter, err := m.meter.Int64Counter(
		name,
		otelmetric.WithDescription(description),
		otelmetric.WithUnit(unit),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create counter: %w", err)
	}
	return counter, nil
}

// RecordCounter increments a counter by a given value.
// The counter must have been created using CreateCounter.
//
// Parameters:
//   - ctx: Context for the metric recording
//   - counter: The counter metric to increment
//   - value: The value to add to the counter (must be non-negative)
//   - labels: Optional key-value pairs for metric dimensions
//
// Example:
//
//	metric.RecordCounter(ctx, counter, 1,
//	    metric.CreateAttributeString("method", "GET"),
//	    metric.CreateAttributeString("status", "200"),
//	)
func (m *metric) RecordCounter(ctx context.Context, counter otelmetric.Int64Counter, value int64, labels ...attribute.KeyValue) {
	counter.Add(ctx, value, otelmetric.WithAttributes(labels...))
}

// CreateHistogram creates a new histogram metric.
// Histograms track the distribution of values over time.
//
// Parameters:
//   - name: The metric name (should follow OpenTelemetry naming conventions)
//   - unit: The unit of measurement (e.g., "ms", "bytes", "seconds")
//   - description: A human-readable description of what the histogram measures
//
// Returns:
//   - The created histogram metric
//   - An error if histogram creation fails
//
// Example:
//
//	histogram, err := metric.CreateHistogram(
//	    "http_request_duration_ms",
//	    "ms",
//	    "HTTP request duration in milliseconds",
//	)
func (m *metric) CreateHistogram(name, unit, description string) (otelmetric.Int64Histogram, error) {
	histogram, err := m.meter.Int64Histogram(
		name,
		otelmetric.WithDescription(description),
		otelmetric.WithUnit(unit),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create histogram: %w", err)
	}
	return histogram, nil
}

// RecordHistogram records a value in a histogram.
// The histogram must have been created using CreateHistogram.
//
// Parameters:
//   - ctx: Context for the metric recording
//   - histogram: The histogram metric to record to
//   - value: The value to record (e.g., request duration, response size)
//   - labels: Optional key-value pairs for metric dimensions
//
// Example:
//
//	start := time.Now()
//	// ... perform operation ...
//	duration := time.Since(start).Milliseconds()
//	metric.RecordHistogram(ctx, histogram, duration,
//	    metric.CreateAttributeString("endpoint", "/api/users"),
//	)
func (m *metric) RecordHistogram(ctx context.Context, histogram otelmetric.Int64Histogram, value int64, labels ...attribute.KeyValue) {
	histogram.Record(ctx, value, otelmetric.WithAttributes(labels...))
}

// CreateAttributeInt creates an integer attribute for metric labels.
// Attributes are used to add dimensions to metrics for filtering and aggregation.
//
// Parameters:
//   - key: The attribute key (should follow OpenTelemetry naming conventions)
//   - value: The integer value
//
// Returns:
//   - An attribute key-value pair
//
// Example:
//
//	attr := metric.CreateAttributeInt("status_code", 200)
//	metric.RecordCounter(ctx, counter, 1, attr)
func (m *metric) CreateAttributeInt(key string, value int) attribute.KeyValue {
	return attribute.Int(key, value)
}

// CreateAttributeString creates a string attribute for metric labels.
// Attributes are used to add dimensions to metrics for filtering and aggregation.
//
// Parameters:
//   - key: The attribute key (should follow OpenTelemetry naming conventions)
//   - value: The string value
//
// Returns:
//   - An attribute key-value pair
//
// Example:
//
//	attr := metric.CreateAttributeString("method", "GET")
//	metric.RecordCounter(ctx, counter, 1, attr)
func (m *metric) CreateAttributeString(key string, value string) attribute.KeyValue {
	return attribute.String(key, value)
}

// Shutdown gracefully shuts down the meter provider.
// It flushes any pending metrics and releases resources.
// This should be called before application shutdown to ensure all metrics are exported.
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
//	if err := metric.Shutdown(ctx); err != nil {
//	    log.Printf("Failed to shutdown metric: %v", err)
//	}
func (m *metric) Shutdown(ctx context.Context) error {
	return m.provider.Shutdown(ctx)
}
