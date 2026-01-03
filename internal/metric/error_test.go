package metric

import (
	"errors"
	"testing"
)

func TestMetric_Error_ErrInvalidProvider(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "ErrInvalidProvider is defined",
			err:  ErrInvalidProvider,
			want: "invalid provider",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Error("ErrInvalidProvider is nil")
				return
			}
			if tt.err.Error() != tt.want {
				t.Errorf("ErrInvalidProvider.Error() = %v, want %v", tt.err.Error(), tt.want)
			}
			if !errors.Is(tt.err, ErrInvalidProvider) {
				t.Error("ErrInvalidProvider should match itself")
			}
		})
	}
}

func TestMetric_Error_ErrProviderHostRequired(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "ErrProviderHostRequired is defined",
			err:  ErrProviderHostRequired,
			want: "provider host is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Error("ErrProviderHostRequired is nil")
				return
			}
			if tt.err.Error() != tt.want {
				t.Errorf("ErrProviderHostRequired.Error() = %v, want %v", tt.err.Error(), tt.want)
			}
			if !errors.Is(tt.err, ErrProviderHostRequired) {
				t.Error("ErrProviderHostRequired should match itself")
			}
		})
	}
}

func TestMetric_Error_ErrProviderPortRequired(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "ErrProviderPortRequired is defined",
			err:  ErrProviderPortRequired,
			want: "provider port is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Error("ErrProviderPortRequired is nil")
				return
			}
			if tt.err.Error() != tt.want {
				t.Errorf("ErrProviderPortRequired.Error() = %v, want %v", tt.err.Error(), tt.want)
			}
			if !errors.Is(tt.err, ErrProviderPortRequired) {
				t.Error("ErrProviderPortRequired should match itself")
			}
		})
	}
}

func TestMetric_Error_ErrProviderPortInvalid(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "ErrProviderPortInvalid is defined",
			err:  ErrProviderPortInvalid,
			want: "provider port must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Error("ErrProviderPortInvalid is nil")
				return
			}
			if tt.err.Error() != tt.want {
				t.Errorf("ErrProviderPortInvalid.Error() = %v, want %v", tt.err.Error(), tt.want)
			}
			if !errors.Is(tt.err, ErrProviderPortInvalid) {
				t.Error("ErrProviderPortInvalid should match itself")
			}
		})
	}
}

func TestMetric_Error_ErrIntervalInvalid(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "ErrIntervalInvalid is defined",
			err:  ErrIntervalInvalid,
			want: "interval must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Error("ErrIntervalInvalid is nil")
				return
			}
			if tt.err.Error() != tt.want {
				t.Errorf("ErrIntervalInvalid.Error() = %v, want %v", tt.err.Error(), tt.want)
			}
			if !errors.Is(tt.err, ErrIntervalInvalid) {
				t.Error("ErrIntervalInvalid should match itself")
			}
		})
	}
}

func TestMetric_Error_ErrInvalidProvider_Usage(t *testing.T) {
	metricInstance, err := NewMetric(
		WithServiceName("test-service"),
		WithProvider("invalid_provider", "", 0),
	)
	if err == nil {
		t.Fatal("NewMetric() with invalid provider expected error")
	}
	if err != ErrInvalidProvider {
		t.Errorf("NewMetric() error = %v, want ErrInvalidProvider", err)
	}
	if metricInstance != nil {
		t.Error("NewMetric() with invalid provider expected nil metric")
	}
	if !errors.Is(err, ErrInvalidProvider) {
		t.Error("errors.Is() should return true for ErrInvalidProvider")
	}
}

func TestMetric_Error_ErrProviderHostRequired_Usage(t *testing.T) {
	metricInstance, err := NewMetric(
		WithServiceName("test-service"),
		WithProvider("otlp", "", 4318),
	)
	if err == nil {
		t.Fatal("NewMetric() with missing host expected error")
	}
	if err != ErrProviderHostRequired {
		t.Errorf("NewMetric() error = %v, want ErrProviderHostRequired", err)
	}
	if metricInstance != nil {
		t.Error("NewMetric() with missing host expected nil metric")
	}
	if !errors.Is(err, ErrProviderHostRequired) {
		t.Error("errors.Is() should return true for ErrProviderHostRequired")
	}
}

func TestMetric_Error_ErrProviderPortRequired_Usage(t *testing.T) {
	metricInstance, err := NewMetric(
		WithServiceName("test-service"),
		WithProvider("otlp", "localhost", 0),
	)
	if err == nil {
		t.Fatal("NewMetric() with missing port expected error")
	}
	if err != ErrProviderPortRequired {
		t.Errorf("NewMetric() error = %v, want ErrProviderPortRequired", err)
	}
	if metricInstance != nil {
		t.Error("NewMetric() with missing port expected nil metric")
	}
	if !errors.Is(err, ErrProviderPortRequired) {
		t.Error("errors.Is() should return true for ErrProviderPortRequired")
	}
}

func TestMetric_Error_ErrProviderPortInvalid_Usage(t *testing.T) {
	metricInstance, err := NewMetric(
		WithServiceName("test-service"),
		WithProvider("otlp", "localhost", -1),
	)
	if err == nil {
		t.Fatal("NewMetric() with invalid port expected error")
	}
	if err != ErrProviderPortInvalid {
		t.Errorf("NewMetric() error = %v, want ErrProviderPortInvalid", err)
	}
	if metricInstance != nil {
		t.Error("NewMetric() with invalid port expected nil metric")
	}
	if !errors.Is(err, ErrProviderPortInvalid) {
		t.Error("errors.Is() should return true for ErrProviderPortInvalid")
	}
}

func TestMetric_Error_ErrIntervalInvalid_Usage(t *testing.T) {
	metricInstance, err := NewMetric(
		WithServiceName("test-service"),
		WithInterval(0),
	)
	if err == nil {
		t.Fatal("NewMetric() with invalid interval expected error")
	}
	if err != ErrIntervalInvalid {
		t.Errorf("NewMetric() error = %v, want ErrIntervalInvalid", err)
	}
	if metricInstance != nil {
		t.Error("NewMetric() with invalid interval expected nil metric")
	}
	if !errors.Is(err, ErrIntervalInvalid) {
		t.Error("errors.Is() should return true for ErrIntervalInvalid")
	}
}

func TestMetric_Error_AllErrorsAreDistinct(t *testing.T) {
	errList := []error{
		ErrInvalidProvider,
		ErrProviderHostRequired,
		ErrProviderPortRequired,
		ErrProviderPortInvalid,
		ErrIntervalInvalid,
	}

	for i, err1 := range errList {
		for j, err2 := range errList {
			if i != j && err1 == err2 {
				t.Errorf("Error %v and %v should be distinct", err1, err2)
			}
			if i != j && errors.Is(err1, err2) {
				t.Errorf("Error %v should not match %v", err1, err2)
			}
		}
	}
}
