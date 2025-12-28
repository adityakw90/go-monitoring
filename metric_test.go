package monitoring

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

func TestNewMetric(t *testing.T) {
	tests := []struct {
		name    string
		opts    []MetricOption
		wantErr bool
	}{
		{
			name:    "default metric with stdout",
			opts:    []MetricOption{withMetricServiceName("test-service")},
			wantErr: false,
		},
		{
			name: "with all options",
			opts: []MetricOption{
				withMetricServiceName("test-service"),
				withMetricEnvironment("test"),
				withMetricInstance("instance-1", "localhost"),
				withMetricProvider("stdout", "", 0),
				withMetricInterval(30 * time.Second),
			},
			wantErr: false,
		},
		{
			name: "with otlp provider",
			opts: []MetricOption{
				withMetricServiceName("test-service"),
				withMetricProvider("otlp", "localhost", 4318),
			},
			wantErr: false,
		},
		{
			name: "with invalid provider",
			opts: []MetricOption{
				withMetricServiceName("test-service"),
				withMetricProvider("invalid", "", 0),
			},
			wantErr: true,
		},
		{
			name: "with custom interval",
			opts: []MetricOption{
				withMetricServiceName("test-service"),
				withMetricInterval(10 * time.Second),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metric, err := NewMetric(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if metric == nil {
					t.Errorf("NewMetric() returned nil")
					return
				}
				if metric.provider == nil {
					t.Errorf("NewMetric() provider is nil")
				}
				if metric.meter == nil {
					t.Errorf("NewMetric() meter is nil")
				}
				// Cleanup
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_ = metric.Shutdown(ctx)
			}
		})
	}
}

func TestMetric_CreateCounter(t *testing.T) {
	m, err := NewMetric(withMetricServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewMetric() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = m.Shutdown(ctx)
	}()

	tests := []struct {
		name        string
		counterName string
		unit        string
		description string
		wantErr     bool
	}{
		{
			name:        "valid counter",
			counterName: "test_counter",
			unit:        "1",
			description: "Test counter description",
			wantErr:     false,
		},
		{
			name:        "counter with empty name",
			counterName: "",
			unit:        "1",
			description: "Test counter",
			wantErr:     true, // OpenTelemetry doesn't allow empty names
		},
		{
			name:        "counter with custom unit",
			counterName: "requests_total",
			unit:        "req",
			description: "Total requests",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter, err := m.CreateCounter(tt.counterName, tt.unit, tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && counter == nil {
				t.Errorf("CreateCounter() returned nil counter")
			}
		})
	}
}

func TestMetric_RecordCounter(t *testing.T) {
	m, err := NewMetric(withMetricServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewMetric() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = m.Shutdown(ctx)
	}()

	counter, err := m.CreateCounter("test_counter", "1", "Test counter")
	if err != nil {
		t.Fatalf("CreateCounter() error = %v", err)
	}

	ctx := context.Background()

	// Test recording without labels
	m.RecordCounter(ctx, counter, 1)

	// Test recording with labels
	m.RecordCounter(ctx, counter, 1,
		attribute.String("method", "GET"),
		attribute.String("status", "200"),
	)

	// Test recording with multiple values
	m.RecordCounter(ctx, counter, 5,
		attribute.String("method", "POST"),
		attribute.Int("code", 201),
	)
}

func TestMetric_CreateHistogram(t *testing.T) {
	m, err := NewMetric(withMetricServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewMetric() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = m.Shutdown(ctx)
	}()

	tests := []struct {
		name          string
		histogramName string
		unit          string
		description   string
		wantErr       bool
	}{
		{
			name:          "valid histogram",
			histogramName: "test_histogram",
			unit:          "ms",
			description:   "Test histogram description",
			wantErr:       false,
		},
		{
			name:          "histogram with duration unit",
			histogramName: "request_duration",
			unit:          "s",
			description:   "Request duration",
			wantErr:       false,
		},
		{
			name:          "histogram with bytes unit",
			histogramName: "response_size",
			unit:          "By",
			description:   "Response size in bytes",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			histogram, err := m.CreateHistogram(tt.histogramName, tt.unit, tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateHistogram() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && histogram == nil {
				t.Errorf("CreateHistogram() returned nil histogram")
			}
		})
	}
}

func TestMetric_RecordHistogram(t *testing.T) {
	m, err := NewMetric(withMetricServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewMetric() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = m.Shutdown(ctx)
	}()

	histogram, err := m.CreateHistogram("test_histogram", "ms", "Test histogram")
	if err != nil {
		t.Fatalf("CreateHistogram() error = %v", err)
	}

	ctx := context.Background()

	// Test recording without labels
	m.RecordHistogram(ctx, histogram, 100)

	// Test recording with labels
	m.RecordHistogram(ctx, histogram, 150,
		attribute.String("method", "GET"),
		attribute.String("endpoint", "/api/users"),
	)

	// Test recording with different values
	m.RecordHistogram(ctx, histogram, 200,
		attribute.String("method", "POST"),
		attribute.Int("status", 201),
	)

	// Test recording zero value
	m.RecordHistogram(ctx, histogram, 0)

	// Test recording large value
	m.RecordHistogram(ctx, histogram, 999999)
}

func TestMetric_CreateAttributeInt(t *testing.T) {
	m, err := NewMetric(withMetricServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewMetric() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = m.Shutdown(ctx)
	}()

	attr := m.CreateAttributeInt("test_key", 42)
	if attr.Key != "test_key" {
		t.Errorf("CreateAttributeInt() key = %v, want test_key", attr.Key)
	}
	if attr.Value.AsInt64() != 42 {
		t.Errorf("CreateAttributeInt() value = %v, want 42", attr.Value.AsInt64())
	}

	// Test with zero
	attrZero := m.CreateAttributeInt("zero", 0)
	if attrZero.Value.AsInt64() != 0 {
		t.Errorf("CreateAttributeInt() zero value = %v, want 0", attrZero.Value.AsInt64())
	}

	// Test with negative
	attrNeg := m.CreateAttributeInt("negative", -10)
	if attrNeg.Value.AsInt64() != -10 {
		t.Errorf("CreateAttributeInt() negative value = %v, want -10", attrNeg.Value.AsInt64())
	}
}

func TestMetric_CreateAttributeString(t *testing.T) {
	m, err := NewMetric(withMetricServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewMetric() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = m.Shutdown(ctx)
	}()

	attr := m.CreateAttributeString("test_key", "test_value")
	if attr.Key != "test_key" {
		t.Errorf("CreateAttributeString() key = %v, want test_key", attr.Key)
	}
	if attr.Value.AsString() != "test_value" {
		t.Errorf("CreateAttributeString() value = %v, want test_value", attr.Value.AsString())
	}

	// Test with empty string
	attrEmpty := m.CreateAttributeString("empty", "")
	if attrEmpty.Value.AsString() != "" {
		t.Errorf("CreateAttributeString() empty value = %v, want empty string", attrEmpty.Value.AsString())
	}

	// Test with special characters
	attrSpecial := m.CreateAttributeString("special", "test-value_123")
	if attrSpecial.Value.AsString() != "test-value_123" {
		t.Errorf("CreateAttributeString() special value = %v, want test-value_123", attrSpecial.Value.AsString())
	}
}

func TestMetric_Shutdown(t *testing.T) {
	m, err := NewMetric(withMetricServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewMetric() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}

	// Second shutdown may return an error (reader is shutdown)
	// This is expected behavior from OpenTelemetry
	_ = m.Shutdown(ctx)
}

func TestMetric_Integration(t *testing.T) {
	m, err := NewMetric(withMetricServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewMetric() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = m.Shutdown(ctx)
	}()

	// Create counter and histogram
	counter, err := m.CreateCounter("requests_total", "1", "Total requests")
	if err != nil {
		t.Fatalf("CreateCounter() error = %v", err)
	}

	histogram, err := m.CreateHistogram("request_duration_ms", "ms", "Request duration")
	if err != nil {
		t.Fatalf("CreateHistogram() error = %v", err)
	}

	ctx := context.Background()

	// Record metrics with attributes
	m.RecordCounter(ctx, counter, 1,
		m.CreateAttributeString("method", "GET"),
		m.CreateAttributeString("status", "200"),
	)

	m.RecordHistogram(ctx, histogram, 150,
		m.CreateAttributeString("method", "GET"),
		m.CreateAttributeInt("status_code", 200),
	)

	// Record multiple times
	for i := 0; i < 10; i++ {
		m.RecordCounter(ctx, counter, 1,
			m.CreateAttributeString("method", "POST"),
		)
		m.RecordHistogram(ctx, histogram, int64(100+i*10),
			m.CreateAttributeString("method", "POST"),
		)
	}
}

func TestMetric_MultipleCounters(t *testing.T) {
	m, err := NewMetric(withMetricServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewMetric() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = m.Shutdown(ctx)
	}()

	// Create multiple counters
	counter1, err := m.CreateCounter("counter1", "1", "First counter")
	if err != nil {
		t.Fatalf("CreateCounter() error = %v", err)
	}

	counter2, err := m.CreateCounter("counter2", "1", "Second counter")
	if err != nil {
		t.Fatalf("CreateCounter() error = %v", err)
	}

	ctx := context.Background()
	m.RecordCounter(ctx, counter1, 1)
	m.RecordCounter(ctx, counter2, 2)

	// Verify they are different instances
	if counter1 == counter2 {
		t.Errorf("CreateCounter() returned same instance for different counters")
	}
}

func TestMetric_MultipleHistograms(t *testing.T) {
	m, err := NewMetric(withMetricServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewMetric() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = m.Shutdown(ctx)
	}()

	// Create multiple histograms
	histogram1, err := m.CreateHistogram("histogram1", "ms", "First histogram")
	if err != nil {
		t.Fatalf("CreateHistogram() error = %v", err)
	}

	histogram2, err := m.CreateHistogram("histogram2", "s", "Second histogram")
	if err != nil {
		t.Fatalf("CreateHistogram() error = %v", err)
	}

	ctx := context.Background()
	m.RecordHistogram(ctx, histogram1, 100)
	m.RecordHistogram(ctx, histogram2, 200)

	// Verify they are different instances
	if histogram1 == histogram2 {
		t.Errorf("CreateHistogram() returned same instance for different histograms")
	}
}
