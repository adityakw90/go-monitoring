package tracer

import (
	"context"
	"testing"
	"time"
)

func TestTracer_NewTracer(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
	}{
		{
			name:    "default tracer with stdout",
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
				WithSampleRatio(0.5),
				WithBatchTimeout(10 * time.Second),
			},
			wantErr: false,
		},
		{
			name: "with otlp provider (insecure)",
			opts: []Option{
				WithServiceName("test-service"),
				WithProvider("otlp", "localhost", 4317),
				WithInsecure(true),
			},
			wantErr: false,
		},
		{
			name: "with otlp provider (secure)",
			opts: []Option{
				WithServiceName("test-service"),
				WithProvider("otlp", "localhost", 4317),
				WithInsecure(false),
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
			name: "with sample ratio 0",
			opts: []Option{
				WithServiceName("test-service"),
				WithSampleRatio(0.0),
			},
			wantErr: false,
		},
		{
			name: "with sample ratio 1.0",
			opts: []Option{
				WithServiceName("test-service"),
				WithSampleRatio(1.0),
			},
			wantErr: false,
		},
		{
			name: "with sample ratio > 1.0 (uses AlwaysSample)",
			opts: []Option{
				WithServiceName("test-service"),
				WithSampleRatio(1.5),
			},
			wantErr: false, // Uses AlwaysSample for ratios >= 1.0
		},
		{
			name: "with sample ratio < 0 (uses NeverSample)",
			opts: []Option{
				WithServiceName("test-service"),
				WithSampleRatio(-0.5),
			},
			wantErr: false, // Uses NeverSample for ratios <= 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracerInstance, err := NewTracer(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTracer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tracerInstance == nil {
					t.Errorf("NewTracer() returned nil")
					return
				}
				if tracerInstance.(*tracer).provider == nil {
					t.Errorf("NewTracer() provider is nil")
				}
				if tracerInstance.(*tracer).tracer == nil {
					t.Errorf("NewTracer() tracer is nil")
				}
				if tracerInstance.(*tracer).propagator == nil {
					t.Errorf("NewTracer() propagator is nil")
				}
				// Cleanup
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				_ = tracerInstance.Shutdown(ctx)
			}
		})
	}
}
