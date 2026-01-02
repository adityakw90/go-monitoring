package metric

import (
	"context"
	"testing"
	"time"
)

func TestMetric_NewMetric(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
	}{
		{
			name:    "default metric with stdout",
			opts:    []Option{WithServiceName("test-service")},
			wantErr: false,
		},
		{
			name: "with all options",
			opts: []Option{
				WithServiceName("test-service"),
				WithEnvironment("test"),
				WithInstance("instance-1", "localhost"),
				WithProvider("stdout", "", 0),
				WithInterval(30 * time.Second),
			},
			wantErr: false,
		},
		{
			name: "with otlp provider",
			opts: []Option{
				WithServiceName("test-service"),
				WithProvider("otlp", "localhost", 4318),
			},
			wantErr: false,
		},
		{
			name: "with invalid provider",
			opts: []Option{
				WithServiceName("test-service"),
				WithProvider("invalid", "", 0),
			},
			wantErr: true,
		},
		{
			name: "with custom interval",
			opts: []Option{
				WithServiceName("test-service"),
				WithInterval(10 * time.Second),
			},
			wantErr: false,
		},
		{
			name: "with insecure option",
			opts: []Option{
				WithServiceName("test-service"),
				WithProvider("otlp", "localhost", 4318),
				WithInsecure(true),
			},
			wantErr: false,
		},
		{
			name: "with secure option (default)",
			opts: []Option{
				WithServiceName("test-service"),
				WithProvider("otlp", "localhost", 4318),
				WithInsecure(false),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metricInstance, err := NewMetric(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if metricInstance == nil {
					t.Errorf("NewMetric() returned nil")
					return
				}
				if metricInstance.(*metric).provider == nil {
					t.Errorf("NewMetric() provider is nil")
				}
				if metricInstance.(*metric).meter == nil {
					t.Errorf("NewMetric() meter is nil")
				}
				// Cleanup
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_ = metricInstance.Shutdown(ctx)
			}
		})
	}
}
