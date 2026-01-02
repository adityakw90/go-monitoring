package monitoring

import "testing"

func TestMonitoring_NewMonitoring(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
	}{
		{
			name:    "missing service name",
			opts:    nil,
			wantErr: true,
		},
		{
			name: "with service name",
			opts: []Option{
				WithServiceName("test-service"),
			},
			wantErr: false,
		},
		{
			name: "with all options",
			opts: []Option{
				WithServiceName("test-service"),
				WithEnvironment("test"),
				WithInstance("instance-1", "localhost"),
				WithLoggerLevel("debug"),
				WithTracerProvider("stdout", "", 0),
				WithMetricProvider("stdout", "", 0),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monitoring, err := NewMonitoring(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMonitoring() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if monitoring == nil {
					t.Errorf("NewMonitoring() returned nil")
					return
				}
				if monitoring.Logger == nil {
					t.Errorf("NewMonitoring() Logger is nil")
				}
				if monitoring.Tracer == nil {
					t.Errorf("NewMonitoring() Tracer is nil")
				}
				if monitoring.Metric == nil {
					t.Errorf("NewMonitoring() Metric is nil")
				}
			}
		})
	}
}
