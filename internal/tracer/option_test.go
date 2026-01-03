package tracer

import (
	"testing"
	"time"
)

func TestTracer_Option_WithServiceName(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		checkFunc   func(t *testing.T, opts *Options)
	}{
		{
			name:        "set service name",
			serviceName: "test-service",
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.ServiceName != "test-service" {
					t.Errorf("WithServiceName() set ServiceName = %v, want %v", opts.ServiceName, "test-service")
				}
			},
		},
		{
			name:        "set empty service name",
			serviceName: "",
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.ServiceName != "" {
					t.Errorf("WithServiceName() set ServiceName = %v, want %v", opts.ServiceName, "")
				}
			},
		},
		{
			name:        "override existing service name",
			serviceName: "new-service",
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.ServiceName != "new-service" {
					t.Errorf("WithServiceName() set ServiceName = %v, want %v", opts.ServiceName, "new-service")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &Options{}
			WithServiceName(tt.serviceName)(opts)
			if tt.checkFunc != nil {
				tt.checkFunc(t, opts)
			}
		})
	}
}

func TestTracer_Option_WithEnvironment(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		checkFunc   func(t *testing.T, opts *Options)
	}{
		{
			name:        "set development environment",
			environment: "development",
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.Environment != "development" {
					t.Errorf("WithEnvironment() set Environment = %v, want %v", opts.Environment, "development")
				}
			},
		},
		{
			name:        "set production environment",
			environment: "production",
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.Environment != "production" {
					t.Errorf("WithEnvironment() set Environment = %v, want %v", opts.Environment, "production")
				}
			},
		},
		{
			name:        "set staging environment",
			environment: "staging",
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.Environment != "staging" {
					t.Errorf("WithEnvironment() set Environment = %v, want %v", opts.Environment, "staging")
				}
			},
		},
		{
			name:        "set empty environment",
			environment: "",
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.Environment != "" {
					t.Errorf("WithEnvironment() set Environment = %v, want %v", opts.Environment, "")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &Options{}
			WithEnvironment(tt.environment)(opts)
			if tt.checkFunc != nil {
				tt.checkFunc(t, opts)
			}
		})
	}
}

func TestTracer_Option_WithInstance(t *testing.T) {
	tests := []struct {
		name         string
		instanceName string
		instanceHost string
		checkFunc    func(t *testing.T, opts *Options)
	}{
		{
			name:         "set instance name and host",
			instanceName: "instance-1",
			instanceHost: "localhost",
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.InstanceName != "instance-1" {
					t.Errorf("WithInstance() set InstanceName = %v, want %v", opts.InstanceName, "instance-1")
				}
				if opts.InstanceHost != "localhost" {
					t.Errorf("WithInstance() set InstanceHost = %v, want %v", opts.InstanceHost, "localhost")
				}
			},
		},
		{
			name:         "set empty instance name and host",
			instanceName: "",
			instanceHost: "",
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.InstanceName != "" {
					t.Errorf("WithInstance() set InstanceName = %v, want %v", opts.InstanceName, "")
				}
				if opts.InstanceHost != "" {
					t.Errorf("WithInstance() set InstanceHost = %v, want %v", opts.InstanceHost, "")
				}
			},
		},
		{
			name:         "set instance with IP address",
			instanceName: "pod-123",
			instanceHost: "192.168.1.100",
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.InstanceName != "pod-123" {
					t.Errorf("WithInstance() set InstanceName = %v, want %v", opts.InstanceName, "pod-123")
				}
				if opts.InstanceHost != "192.168.1.100" {
					t.Errorf("WithInstance() set InstanceHost = %v, want %v", opts.InstanceHost, "192.168.1.100")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &Options{}
			WithInstance(tt.instanceName, tt.instanceHost)(opts)
			if tt.checkFunc != nil {
				tt.checkFunc(t, opts)
			}
		})
	}
}

func TestTracer_Option_WithProvider(t *testing.T) {
	tests := []struct {
		name      string
		provider  string
		host      string
		port      int
		checkFunc func(t *testing.T, opts *Options)
	}{
		{
			name:     "set stdout provider",
			provider: "stdout",
			host:     "",
			port:     0,
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.Provider != "stdout" {
					t.Errorf("WithProvider() set Provider = %v, want %v", opts.Provider, "stdout")
				}
				if opts.ProviderHost != "" {
					t.Errorf("WithProvider() set ProviderHost = %v, want %v", opts.ProviderHost, "")
				}
				if opts.ProviderPort != 0 {
					t.Errorf("WithProvider() set ProviderPort = %v, want %v", opts.ProviderPort, 0)
				}
			},
		},
		{
			name:     "set otlp provider with host and port",
			provider: "otlp",
			host:     "localhost",
			port:     4317,
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.Provider != "otlp" {
					t.Errorf("WithProvider() set Provider = %v, want %v", opts.Provider, "otlp")
				}
				if opts.ProviderHost != "localhost" {
					t.Errorf("WithProvider() set ProviderHost = %v, want %v", opts.ProviderHost, "localhost")
				}
				if opts.ProviderPort != 4317 {
					t.Errorf("WithProvider() set ProviderPort = %v, want %v", opts.ProviderPort, 4317)
				}
			},
		},
		{
			name:     "set otlp provider with custom host and port",
			provider: "otlp",
			host:     "collector.example.com",
			port:     4318,
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.Provider != "otlp" {
					t.Errorf("WithProvider() set Provider = %v, want %v", opts.Provider, "otlp")
				}
				if opts.ProviderHost != "collector.example.com" {
					t.Errorf("WithProvider() set ProviderHost = %v, want %v", opts.ProviderHost, "collector.example.com")
				}
				if opts.ProviderPort != 4318 {
					t.Errorf("WithProvider() set ProviderPort = %v, want %v", opts.ProviderPort, 4318)
				}
			},
		},
		{
			name:     "override existing provider",
			provider: "otlp",
			host:     "new-host",
			port:     9999,
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.Provider != "otlp" {
					t.Errorf("WithProvider() set Provider = %v, want %v", opts.Provider, "otlp")
				}
				if opts.ProviderHost != "new-host" {
					t.Errorf("WithProvider() set ProviderHost = %v, want %v", opts.ProviderHost, "new-host")
				}
				if opts.ProviderPort != 9999 {
					t.Errorf("WithProvider() set ProviderPort = %v, want %v", opts.ProviderPort, 9999)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &Options{}
			WithProvider(tt.provider, tt.host, tt.port)(opts)
			if tt.checkFunc != nil {
				tt.checkFunc(t, opts)
			}
		})
	}
}

func TestTracer_Option_WithSampleRatio(t *testing.T) {
	tests := []struct {
		name      string
		ratio     float64
		checkFunc func(t *testing.T, opts *Options)
	}{
		{
			name:  "set sample ratio 0.0",
			ratio: 0.0,
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.SampleRatio != 0.0 {
					t.Errorf("WithSampleRatio() set SampleRatio = %v, want %v", opts.SampleRatio, 0.0)
				}
			},
		},
		{
			name:  "set sample ratio 0.5",
			ratio: 0.5,
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.SampleRatio != 0.5 {
					t.Errorf("WithSampleRatio() set SampleRatio = %v, want %v", opts.SampleRatio, 0.5)
				}
			},
		},
		{
			name:  "set sample ratio 1.0",
			ratio: 1.0,
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.SampleRatio != 1.0 {
					t.Errorf("WithSampleRatio() set SampleRatio = %v, want %v", opts.SampleRatio, 1.0)
				}
			},
		},
		{
			name:  "set sample ratio > 1.0",
			ratio: 1.5,
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.SampleRatio != 1.5 {
					t.Errorf("WithSampleRatio() set SampleRatio = %v, want %v", opts.SampleRatio, 1.5)
				}
			},
		},
		{
			name:  "set sample ratio < 0",
			ratio: -0.1,
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.SampleRatio != -0.1 {
					t.Errorf("WithSampleRatio() set SampleRatio = %v, want %v", opts.SampleRatio, -0.1)
				}
			},
		},
		{
			name:  "set small sample ratio",
			ratio: 0.001,
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.SampleRatio != 0.001 {
					t.Errorf("WithSampleRatio() set SampleRatio = %v, want %v", opts.SampleRatio, 0.001)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &Options{}
			WithSampleRatio(tt.ratio)(opts)
			if tt.checkFunc != nil {
				tt.checkFunc(t, opts)
			}
		})
	}
}

func TestTracer_Option_WithBatchTimeout(t *testing.T) {
	tests := []struct {
		name      string
		timeout   time.Duration
		checkFunc func(t *testing.T, opts *Options)
	}{
		{
			name:    "set batch timeout 1 second",
			timeout: 1 * time.Second,
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.BatchTimeout != 1*time.Second {
					t.Errorf("WithBatchTimeout() set BatchTimeout = %v, want %v", opts.BatchTimeout, 1*time.Second)
				}
			},
		},
		{
			name:    "set batch timeout 5 seconds",
			timeout: 5 * time.Second,
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.BatchTimeout != 5*time.Second {
					t.Errorf("WithBatchTimeout() set BatchTimeout = %v, want %v", opts.BatchTimeout, 5*time.Second)
				}
			},
		},
		{
			name:    "set batch timeout 10 seconds",
			timeout: 10 * time.Second,
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.BatchTimeout != 10*time.Second {
					t.Errorf("WithBatchTimeout() set BatchTimeout = %v, want %v", opts.BatchTimeout, 10*time.Second)
				}
			},
		},
		{
			name:    "set batch timeout 100 milliseconds",
			timeout: 100 * time.Millisecond,
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.BatchTimeout != 100*time.Millisecond {
					t.Errorf("WithBatchTimeout() set BatchTimeout = %v, want %v", opts.BatchTimeout, 100*time.Millisecond)
				}
			},
		},
		{
			name:    "set batch timeout 1 minute",
			timeout: 1 * time.Minute,
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.BatchTimeout != 1*time.Minute {
					t.Errorf("WithBatchTimeout() set BatchTimeout = %v, want %v", opts.BatchTimeout, 1*time.Minute)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &Options{}
			WithBatchTimeout(tt.timeout)(opts)
			if tt.checkFunc != nil {
				tt.checkFunc(t, opts)
			}
		})
	}
}

func TestTracer_Option_WithInsecure(t *testing.T) {
	tests := []struct {
		name      string
		insecure  bool
		checkFunc func(t *testing.T, opts *Options)
	}{
		{
			name:     "set insecure true",
			insecure: true,
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.Insecure != true {
					t.Errorf("WithInsecure() set Insecure = %v, want %v", opts.Insecure, true)
				}
			},
		},
		{
			name:     "set insecure false",
			insecure: false,
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.Insecure != false {
					t.Errorf("WithInsecure() set Insecure = %v, want %v", opts.Insecure, false)
				}
			},
		},
		{
			name:     "override existing insecure setting",
			insecure: true,
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.Insecure != true {
					t.Errorf("WithInsecure() set Insecure = %v, want %v", opts.Insecure, true)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &Options{}
			WithInsecure(tt.insecure)(opts)
			if tt.checkFunc != nil {
				tt.checkFunc(t, opts)
			}
		})
	}
}

func TestTracer_Option_MultipleOptions(t *testing.T) {
	tests := []struct {
		name      string
		opts      []Option
		checkFunc func(t *testing.T, opts *Options)
	}{
		{
			name: "apply multiple options",
			opts: []Option{
				WithServiceName("test-service"),
				WithEnvironment("production"),
				WithInstance("instance-1", "localhost"),
				WithProvider("otlp", "collector.example.com", 4317),
				WithSampleRatio(0.1),
				WithBatchTimeout(10 * time.Second),
				WithInsecure(true),
			},
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.ServiceName != "test-service" {
					t.Errorf("ServiceName = %v, want %v", opts.ServiceName, "test-service")
				}
				if opts.Environment != "production" {
					t.Errorf("Environment = %v, want %v", opts.Environment, "production")
				}
				if opts.InstanceName != "instance-1" {
					t.Errorf("InstanceName = %v, want %v", opts.InstanceName, "instance-1")
				}
				if opts.InstanceHost != "localhost" {
					t.Errorf("InstanceHost = %v, want %v", opts.InstanceHost, "localhost")
				}
				if opts.Provider != "otlp" {
					t.Errorf("Provider = %v, want %v", opts.Provider, "otlp")
				}
				if opts.ProviderHost != "collector.example.com" {
					t.Errorf("ProviderHost = %v, want %v", opts.ProviderHost, "collector.example.com")
				}
				if opts.ProviderPort != 4317 {
					t.Errorf("ProviderPort = %v, want %v", opts.ProviderPort, 4317)
				}
				if opts.SampleRatio != 0.1 {
					t.Errorf("SampleRatio = %v, want %v", opts.SampleRatio, 0.1)
				}
				if opts.BatchTimeout != 10*time.Second {
					t.Errorf("BatchTimeout = %v, want %v", opts.BatchTimeout, 10*time.Second)
				}
				if opts.Insecure != true {
					t.Errorf("Insecure = %v, want %v", opts.Insecure, true)
				}
			},
		},
		{
			name: "apply options in different order",
			opts: []Option{
				WithInsecure(false),
				WithBatchTimeout(5 * time.Second),
				WithSampleRatio(1.0),
				WithProvider("stdout", "", 0),
				WithInstance("pod-123", "192.168.1.1"),
				WithEnvironment("development"),
				WithServiceName("my-service"),
			},
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.ServiceName != "my-service" {
					t.Errorf("ServiceName = %v, want %v", opts.ServiceName, "my-service")
				}
				if opts.Environment != "development" {
					t.Errorf("Environment = %v, want %v", opts.Environment, "development")
				}
				if opts.InstanceName != "pod-123" {
					t.Errorf("InstanceName = %v, want %v", opts.InstanceName, "pod-123")
				}
				if opts.InstanceHost != "192.168.1.1" {
					t.Errorf("InstanceHost = %v, want %v", opts.InstanceHost, "192.168.1.1")
				}
				if opts.Provider != "stdout" {
					t.Errorf("Provider = %v, want %v", opts.Provider, "stdout")
				}
				if opts.SampleRatio != 1.0 {
					t.Errorf("SampleRatio = %v, want %v", opts.SampleRatio, 1.0)
				}
				if opts.BatchTimeout != 5*time.Second {
					t.Errorf("BatchTimeout = %v, want %v", opts.BatchTimeout, 5*time.Second)
				}
				if opts.Insecure != false {
					t.Errorf("Insecure = %v, want %v", opts.Insecure, false)
				}
			},
		},
		{
			name: "override options with later options",
			opts: []Option{
				WithServiceName("first-service"),
				WithEnvironment("first-env"),
				WithServiceName("second-service"),
				WithEnvironment("second-env"),
			},
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.ServiceName != "second-service" {
					t.Errorf("ServiceName = %v, want %v", opts.ServiceName, "second-service")
				}
				if opts.Environment != "second-env" {
					t.Errorf("Environment = %v, want %v", opts.Environment, "second-env")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &Options{}
			for _, opt := range tt.opts {
				opt(opts)
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, opts)
			}
		})
	}
}
