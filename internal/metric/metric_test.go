package metric

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

func TestMetric_Metric_CreateCounter(t *testing.T) {
	metricInstance, err := NewMetric(WithServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewMetric() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = metricInstance.Shutdown(ctx)
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
			counter, err := metricInstance.CreateCounter(tt.counterName, tt.unit, tt.description)
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

func TestMetric_Metric_RecordCounter(t *testing.T) {
	metricInstance, err := NewMetric(WithServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewMetric() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = metricInstance.Shutdown(ctx)
	}()

	counter, err := metricInstance.CreateCounter("test_counter", "1", "Test counter")
	if err != nil {
		t.Fatalf("CreateCounter() error = %v", err)
	}

	ctx := context.Background()

	// Test recording without labels
	metricInstance.RecordCounter(ctx, counter, 1)

	// Test recording with labels
	metricInstance.RecordCounter(ctx, counter, 1,
		attribute.String("method", "GET"),
		attribute.String("status", "200"),
	)

	// Test recording with multiple values
	metricInstance.RecordCounter(ctx, counter, 5,
		attribute.String("method", "POST"),
		attribute.Int("code", 201),
	)
}

func TestMetric_Metric_CreateHistogram(t *testing.T) {
	metricInstance, err := NewMetric(WithServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewMetric() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = metricInstance.Shutdown(ctx)
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
			histogram, err := metricInstance.CreateHistogram(tt.histogramName, tt.unit, tt.description)
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

func TestMetric_Metric_RecordHistogram(t *testing.T) {
	metricInstance, err := NewMetric(WithServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewMetric() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = metricInstance.Shutdown(ctx)
	}()

	histogram, err := metricInstance.CreateHistogram("test_histogram", "ms", "Test histogram")
	if err != nil {
		t.Fatalf("CreateHistogram() error = %v", err)
	}

	ctx := context.Background()

	// Test recording without labels
	metricInstance.RecordHistogram(ctx, histogram, 100)

	// Test recording with labels
	metricInstance.RecordHistogram(ctx, histogram, 150,
		attribute.String("method", "GET"),
		attribute.String("endpoint", "/api/users"),
	)

	// Test recording with different values
	metricInstance.RecordHistogram(ctx, histogram, 200,
		attribute.String("method", "POST"),
		attribute.Int("status", 201),
	)

	// Test recording zero value
	metricInstance.RecordHistogram(ctx, histogram, 0)

	// Test recording large value
	metricInstance.RecordHistogram(ctx, histogram, 999999)
}

func TestMetric_Metric_CreateAttributeInt(t *testing.T) {
	metricInstance, err := NewMetric(WithServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewMetric() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = metricInstance.Shutdown(ctx)
	}()

	attr := metricInstance.CreateAttributeInt("test_key", 42)
	if attr.Key != "test_key" {
		t.Errorf("CreateAttributeInt() key = %v, want test_key", attr.Key)
	}
	if attr.Value.AsInt64() != 42 {
		t.Errorf("CreateAttributeInt() value = %v, want 42", attr.Value.AsInt64())
	}

	// Test with zero
	attrZero := metricInstance.CreateAttributeInt("zero", 0)
	if attrZero.Value.AsInt64() != 0 {
		t.Errorf("CreateAttributeInt() zero value = %v, want 0", attrZero.Value.AsInt64())
	}

	// Test with negative
	attrNeg := metricInstance.CreateAttributeInt("negative", -10)
	if attrNeg.Value.AsInt64() != -10 {
		t.Errorf("CreateAttributeInt() negative value = %v, want -10", attrNeg.Value.AsInt64())
	}
}

func TestMetric_Metric_CreateAttributeString(t *testing.T) {
	metricInstance, err := NewMetric(WithServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewMetric() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = metricInstance.Shutdown(ctx)
	}()

	attr := metricInstance.CreateAttributeString("test_key", "test_value")
	if attr.Key != "test_key" {
		t.Errorf("CreateAttributeString() key = %v, want test_key", attr.Key)
	}
	if attr.Value.AsString() != "test_value" {
		t.Errorf("CreateAttributeString() value = %v, want test_value", attr.Value.AsString())
	}

	// Test with empty string
	attrEmpty := metricInstance.CreateAttributeString("empty", "")
	if attrEmpty.Value.AsString() != "" {
		t.Errorf("CreateAttributeString() empty value = %v, want empty string", attrEmpty.Value.AsString())
	}

	// Test with special characters
	attrSpecial := metricInstance.CreateAttributeString("special", "test-value_123")
	if attrSpecial.Value.AsString() != "test-value_123" {
		t.Errorf("CreateAttributeString() special value = %v, want test-value_123", attrSpecial.Value.AsString())
	}
}

func TestMetric_Metric_Shutdown(t *testing.T) {
	metricInstance, err := NewMetric(WithServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewMetric() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := metricInstance.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}

	// Second shutdown may return an error (reader is shutdown)
	// This is expected behavior from OpenTelemetry
	_ = metricInstance.Shutdown(ctx)
}

func TestMetric_Metric_Integration(t *testing.T) {
	metricInstance, err := NewMetric(WithServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewMetric() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = metricInstance.Shutdown(ctx)
	}()

	// Create counter and histogram
	counter, err := metricInstance.CreateCounter("requests_total", "1", "Total requests")
	if err != nil {
		t.Fatalf("CreateCounter() error = %v", err)
	}

	histogram, err := metricInstance.CreateHistogram("request_duration_ms", "ms", "Request duration")
	if err != nil {
		t.Fatalf("CreateHistogram() error = %v", err)
	}

	ctx := context.Background()

	// Record metrics with attributes
	metricInstance.RecordCounter(ctx, counter, 1,
		metricInstance.CreateAttributeString("method", "GET"),
		metricInstance.CreateAttributeString("status", "200"),
	)

	metricInstance.RecordHistogram(ctx, histogram, 150,
		metricInstance.CreateAttributeString("method", "GET"),
		metricInstance.CreateAttributeInt("status_code", 200),
	)

	// Record multiple times
	for i := 0; i < 10; i++ {
		metricInstance.RecordCounter(ctx, counter, 1,
			metricInstance.CreateAttributeString("method", "POST"),
		)
		metricInstance.RecordHistogram(ctx, histogram, int64(100+i*10),
			metricInstance.CreateAttributeString("method", "POST"),
		)
	}
}

func TestMetric_Metric_MultipleCounters(t *testing.T) {
	metricInstance, err := NewMetric(WithServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewMetric() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = metricInstance.Shutdown(ctx)
	}()

	// Create multiple counters
	counter1, err := metricInstance.CreateCounter("counter1", "1", "First counter")
	if err != nil {
		t.Fatalf("CreateCounter() error = %v", err)
	}

	counter2, err := metricInstance.CreateCounter("counter2", "1", "Second counter")
	if err != nil {
		t.Fatalf("CreateCounter() error = %v", err)
	}

	ctx := context.Background()
	metricInstance.RecordCounter(ctx, counter1, 1)
	metricInstance.RecordCounter(ctx, counter2, 2)

	// Verify they are different instances
	if counter1 == counter2 {
		t.Errorf("CreateCounter() returned same instance for different counters")
	}
}

func TestMetric_Metric_MultipleHistograms(t *testing.T) {
	metricInstance, err := NewMetric(WithServiceName("test-service"))
	if err != nil {
		t.Fatalf("NewMetric() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = metricInstance.Shutdown(ctx)
	}()

	// Create multiple histograms
	histogram1, err := metricInstance.CreateHistogram("histogram1", "ms", "First histogram")
	if err != nil {
		t.Fatalf("CreateHistogram() error = %v", err)
	}

	histogram2, err := metricInstance.CreateHistogram("histogram2", "s", "Second histogram")
	if err != nil {
		t.Fatalf("CreateHistogram() error = %v", err)
	}

	ctx := context.Background()
	metricInstance.RecordHistogram(ctx, histogram1, 100)
	metricInstance.RecordHistogram(ctx, histogram2, 200)

	// Verify they are different instances
	if histogram1 == histogram2 {
		t.Errorf("CreateHistogram() returned same instance for different histograms")
	}
}

func TestMetric_Metric_MultipleInstancesCoexist(t *testing.T) {
	// Test that multiple Metric instances can coexist without global state conflicts
	// This verifies that we removed global state mutations
	metricInstance1, err := NewMetric(
		WithServiceName("service-1"),
		WithProvider("stdout", "", 0),
	)
	if err != nil {
		t.Fatalf("NewMetric() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = metricInstance1.Shutdown(ctx)
	}()

	metricInstance2, err := NewMetric(
		WithServiceName("service-2"),
		WithProvider("stdout", "", 0),
	)
	if err != nil {
		t.Fatalf("NewMetric() error = %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = metricInstance2.Shutdown(ctx)
	}()

	// Verify they are different instances
	if metricInstance1 == metricInstance2 {
		t.Errorf("NewMetric() returned same instance for different metrics")
	}

	// Verify they have different providers
	if metricInstance1.(*metric).provider == metricInstance2.(*metric).provider {
		t.Errorf("NewMetric() returned same provider for different metrics")
	}

	// Verify they have different meters
	if metricInstance1.(*metric).meter == metricInstance2.(*metric).meter {
		t.Errorf("NewMetric() returned same meter for different metrics")
	}

	// Create counters from both instances and verify they work independently
	counter1, err := metricInstance1.CreateCounter("counter1", "1", "Counter from metric1")
	if err != nil {
		t.Fatalf("metric1.CreateCounter() error = %v", err)
	}

	counter2, err := metricInstance2.CreateCounter("counter2", "1", "Counter from metric2")
	if err != nil {
		t.Fatalf("metric2.CreateCounter() error = %v", err)
	}

	ctx := context.Background()
	metricInstance1.RecordCounter(ctx, counter1, 1)
	metricInstance2.RecordCounter(ctx, counter2, 2)

	// Verify counters are different
	if counter1 == counter2 {
		t.Errorf("CreateCounter() returned same instance from different metrics")
	}
}
