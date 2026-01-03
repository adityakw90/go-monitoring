package monitoring

import (
	"testing"
	"time"
)

func TestMonitoring_Options_DefaultOptions(t *testing.T) {
	opts := defaultOptions()

	tests := []struct {
		name string
		got  interface{}
		want interface{}
	}{
		{"Environment", opts.Environment, "development"},
		{"LoggerLevel", opts.LoggerLevel, "info"},
		{"LoggerOutputPath", opts.LoggerOutputPath, ""},
		{"TracerProvider", opts.TracerProvider, "stdout"},
		{"TracerSampleRatio", opts.TracerSampleRatio, 1.0},
		{"TracerBatchTimeout", opts.TracerBatchTimeout, 5 * time.Second},
		{"TracerInsecure", opts.TracerInsecure, false},
		{"MetricProvider", opts.MetricProvider, "stdout"},
		{"MetricInterval", opts.MetricInterval, 60 * time.Second},
		{"MetricInsecure", opts.MetricInsecure, false},
		{"ServiceName", opts.ServiceName, ""},
		{"InstanceName", opts.InstanceName, ""},
		{"InstanceHost", opts.InstanceHost, ""},
		{"TracerProviderHost", opts.TracerProviderHost, ""},
		{"TracerProviderPort", opts.TracerProviderPort, 0},
		{"MetricProviderHost", opts.MetricProviderHost, ""},
		{"MetricProviderPort", opts.MetricProviderPort, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("defaultOptions() %s = %v, want %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestMonitoring_Options_WithServiceName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"test-service", "test-service"},
		{"my-api-service", "my-api-service"},
		{"", ""},
		{"service-with-dashes", "service-with-dashes"},
		{"service_with_underscores", "service_with_underscores"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := defaultOptions()
			WithServiceName(tt.name)(opts)
			if opts.ServiceName != tt.want {
				t.Errorf("WithServiceName(%q) ServiceName = %v, want %v", tt.name, opts.ServiceName, tt.want)
			}
		})
	}
}

func TestMonitoring_Options_WithEnvironment(t *testing.T) {
	tests := []struct {
		env  string
		want string
	}{
		{"development", "development"},
		{"staging", "staging"},
		{"production", "production"},
		{"test", "test"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			opts := defaultOptions()
			WithEnvironment(tt.env)(opts)
			if opts.Environment != tt.want {
				t.Errorf("WithEnvironment(%q) Environment = %v, want %v", tt.env, opts.Environment, tt.want)
			}
		})
	}
}

func TestMonitoring_Options_WithInstance(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		wantName string
		wantHost string
	}{
		{"instance-1", "localhost", "instance-1", "localhost"},
		{"pod-abc123", "10.0.0.1", "pod-abc123", "10.0.0.1"},
		{"", "", "", ""},
		{"instance-name", "", "instance-name", ""},
		{"", "hostname", "", "hostname"},
		{"instance.with.dots", "host.with.dots", "instance.with.dots", "host.with.dots"},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_"+tt.host, func(t *testing.T) {
			opts := defaultOptions()
			WithInstance(tt.name, tt.host)(opts)
			if opts.InstanceName != tt.wantName {
				t.Errorf("WithInstance(%q, %q) InstanceName = %v, want %v", tt.name, tt.host, opts.InstanceName, tt.wantName)
			}
			if opts.InstanceHost != tt.wantHost {
				t.Errorf("WithInstance(%q, %q) InstanceHost = %v, want %v", tt.name, tt.host, opts.InstanceHost, tt.wantHost)
			}
		})
	}
}

func TestMonitoring_Options_WithLoggerLevel(t *testing.T) {
	tests := []struct {
		level string
		want  string
	}{
		{"debug", "debug"},
		{"info", "info"},
		{"warn", "warn"},
		{"error", "error"},
		{"fatal", "fatal"},
		{"", ""},
		{"DEBUG", "DEBUG"},
		{"invalid", "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			opts := defaultOptions()
			WithLoggerLevel(tt.level)(opts)
			if opts.LoggerLevel != tt.want {
				t.Errorf("WithLoggerLevel(%q) LoggerLevel = %v, want %v", tt.level, opts.LoggerLevel, tt.want)
			}
		})
	}
}

func TestMonitoring_Options_WithLoggerOutputPath(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/var/log/app.log", "/var/log/app.log"},
		{"./logs/app.log", "./logs/app.log"},
		{"logs/app.log", "logs/app.log"},
		{"", ""},
		{"/tmp/test.log", "/tmp/test.log"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			opts := defaultOptions()
			WithLoggerOutputPath(tt.path)(opts)
			if opts.LoggerOutputPath != tt.want {
				t.Errorf("WithLoggerOutputPath(%q) LoggerOutputPath = %v, want %v", tt.path, opts.LoggerOutputPath, tt.want)
			}
		})
	}
}

func TestMonitoring_Options_WithTracerProvider(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		host     string
		port     int
		wantProv string
		wantHost string
		wantPort int
	}{
		{"stdout", "stdout", "localhost", 4317, "stdout", "localhost", 4317},
		{"otlp", "otlp", "localhost", 4317, "otlp", "localhost", 4317},
		{"otlp_custom", "otlp", "collector.example.com", 4318, "otlp", "collector.example.com", 4318},
		{"empty_provider", "", "", 0, "", "", 0},
		{"zero_port", "otlp", "localhost", 0, "otlp", "localhost", 0},
		{"negative_port", "otlp", "localhost", -1, "otlp", "localhost", -1},
		{"large_port", "otlp", "localhost", 65535, "otlp", "localhost", 65535},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := defaultOptions()
			WithTracerProvider(tt.provider, tt.host, tt.port)(opts)
			if opts.TracerProvider != tt.wantProv {
				t.Errorf("WithTracerProvider() TracerProvider = %v, want %v", opts.TracerProvider, tt.wantProv)
			}
			if opts.TracerProviderHost != tt.wantHost {
				t.Errorf("WithTracerProvider() TracerProviderHost = %v, want %v", opts.TracerProviderHost, tt.wantHost)
			}
			if opts.TracerProviderPort != tt.wantPort {
				t.Errorf("WithTracerProvider() TracerProviderPort = %v, want %v", opts.TracerProviderPort, tt.wantPort)
			}
		})
	}
}

func TestMonitoring_Options_WithTracerSampleRatio(t *testing.T) {
	tests := []struct {
		name  string
		ratio float64
		want  float64
	}{
		{"zero", 0.0, 0.0},
		{"one", 1.0, 1.0},
		{"half", 0.5, 0.5},
		{"ten_percent", 0.1, 0.1},
		{"ninety_percent", 0.9, 0.9},
		{"negative", -0.1, -0.1},
		{"over_one", 1.5, 1.5},
		{"very_small", 0.0001, 0.0001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := defaultOptions()
			WithTracerSampleRatio(tt.ratio)(opts)
			if opts.TracerSampleRatio != tt.want {
				t.Errorf("WithTracerSampleRatio(%v) TracerSampleRatio = %v, want %v", tt.ratio, opts.TracerSampleRatio, tt.want)
			}
		})
	}
}

func TestMonitoring_Options_WithTracerBatchTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		want    time.Duration
	}{
		{"one_second", 1 * time.Second, 1 * time.Second},
		{"five_seconds", 5 * time.Second, 5 * time.Second},
		{"ten_seconds", 10 * time.Second, 10 * time.Second},
		{"one_minute", 1 * time.Minute, 1 * time.Minute},
		{"zero", 0, 0},
		{"negative", -1 * time.Second, -1 * time.Second},
		{"millisecond", 100 * time.Millisecond, 100 * time.Millisecond},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := defaultOptions()
			WithTracerBatchTimeout(tt.timeout)(opts)
			if opts.TracerBatchTimeout != tt.want {
				t.Errorf("WithTracerBatchTimeout(%v) TracerBatchTimeout = %v, want %v", tt.timeout, opts.TracerBatchTimeout, tt.want)
			}
		})
	}
}

func TestMonitoring_Options_WithTracerInsecure(t *testing.T) {
	tests := []struct {
		name     string
		insecure bool
		want     bool
	}{
		{"true", true, true},
		{"false", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := defaultOptions()
			WithTracerInsecure(tt.insecure)(opts)
			if opts.TracerInsecure != tt.want {
				t.Errorf("WithTracerInsecure(%v) TracerInsecure = %v, want %v", tt.insecure, opts.TracerInsecure, tt.want)
			}
		})
	}
}

func TestMonitoring_Options_WithMetricProvider(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		host     string
		port     int
		wantProv string
		wantHost string
		wantPort int
	}{
		{"stdout", "stdout", "localhost", 4318, "stdout", "localhost", 4318},
		{"otlp", "otlp", "localhost", 4318, "otlp", "localhost", 4318},
		{"otlp_custom", "otlp", "collector.example.com", 4319, "otlp", "collector.example.com", 4319},
		{"empty_provider", "", "", 0, "", "", 0},
		{"zero_port", "otlp", "localhost", 0, "otlp", "localhost", 0},
		{"negative_port", "otlp", "localhost", -1, "otlp", "localhost", -1},
		{"large_port", "otlp", "localhost", 65535, "otlp", "localhost", 65535},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := defaultOptions()
			WithMetricProvider(tt.provider, tt.host, tt.port)(opts)
			if opts.MetricProvider != tt.wantProv {
				t.Errorf("WithMetricProvider() MetricProvider = %v, want %v", opts.MetricProvider, tt.wantProv)
			}
			if opts.MetricProviderHost != tt.wantHost {
				t.Errorf("WithMetricProvider() MetricProviderHost = %v, want %v", opts.MetricProviderHost, tt.wantHost)
			}
			if opts.MetricProviderPort != tt.wantPort {
				t.Errorf("WithMetricProvider() MetricProviderPort = %v, want %v", opts.MetricProviderPort, tt.wantPort)
			}
		})
	}
}

func TestMonitoring_Options_WithMetricInterval(t *testing.T) {
	tests := []struct {
		name     string
		interval time.Duration
		want     time.Duration
	}{
		{"thirty_seconds", 30 * time.Second, 30 * time.Second},
		{"sixty_seconds", 60 * time.Second, 60 * time.Second},
		{"one_minute", 1 * time.Minute, 1 * time.Minute},
		{"five_minutes", 5 * time.Minute, 5 * time.Minute},
		{"zero", 0, 0},
		{"negative", -1 * time.Second, -1 * time.Second},
		{"one_second", 1 * time.Second, 1 * time.Second},
		{"millisecond", 500 * time.Millisecond, 500 * time.Millisecond},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := defaultOptions()
			WithMetricInterval(tt.interval)(opts)
			if opts.MetricInterval != tt.want {
				t.Errorf("WithMetricInterval(%v) MetricInterval = %v, want %v", tt.interval, opts.MetricInterval, tt.want)
			}
		})
	}
}

func TestMonitoring_Options_WithMetricInsecure(t *testing.T) {
	tests := []struct {
		name     string
		insecure bool
		want     bool
	}{
		{"true", true, true},
		{"false", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := defaultOptions()
			WithMetricInsecure(tt.insecure)(opts)
			if opts.MetricInsecure != tt.want {
				t.Errorf("WithMetricInsecure(%v) MetricInsecure = %v, want %v", tt.insecure, opts.MetricInsecure, tt.want)
			}
		})
	}
}

func TestMonitoring_Options_Combined(t *testing.T) {
	opts := defaultOptions()

	WithServiceName("test-service")(opts)
	WithEnvironment("production")(opts)
	WithInstance("instance-1", "localhost")(opts)
	WithLoggerLevel("debug")(opts)
	WithLoggerOutputPath("/var/log/app.log")(opts)
	WithTracerProvider("otlp", "localhost", 4317)(opts)
	WithTracerSampleRatio(0.5)(opts)
	WithTracerBatchTimeout(10 * time.Second)(opts)
	WithTracerInsecure(true)(opts)
	WithMetricProvider("otlp", "localhost", 4318)(opts)
	WithMetricInterval(30 * time.Second)(opts)
	WithMetricInsecure(true)(opts)

	if opts.ServiceName != "test-service" {
		t.Errorf("ServiceName = %v, want test-service", opts.ServiceName)
	}
	if opts.Environment != "production" {
		t.Errorf("Environment = %v, want production", opts.Environment)
	}
	if opts.InstanceName != "instance-1" {
		t.Errorf("InstanceName = %v, want instance-1", opts.InstanceName)
	}
	if opts.InstanceHost != "localhost" {
		t.Errorf("InstanceHost = %v, want localhost", opts.InstanceHost)
	}
	if opts.LoggerLevel != "debug" {
		t.Errorf("LoggerLevel = %v, want debug", opts.LoggerLevel)
	}
	if opts.LoggerOutputPath != "/var/log/app.log" {
		t.Errorf("LoggerOutputPath = %v, want /var/log/app.log", opts.LoggerOutputPath)
	}
	if opts.TracerProvider != "otlp" {
		t.Errorf("TracerProvider = %v, want otlp", opts.TracerProvider)
	}
	if opts.TracerProviderHost != "localhost" {
		t.Errorf("TracerProviderHost = %v, want localhost", opts.TracerProviderHost)
	}
	if opts.TracerProviderPort != 4317 {
		t.Errorf("TracerProviderPort = %v, want 4317", opts.TracerProviderPort)
	}
	if opts.TracerSampleRatio != 0.5 {
		t.Errorf("TracerSampleRatio = %v, want 0.5", opts.TracerSampleRatio)
	}
	if opts.TracerBatchTimeout != 10*time.Second {
		t.Errorf("TracerBatchTimeout = %v, want 10s", opts.TracerBatchTimeout)
	}
	if opts.TracerInsecure != true {
		t.Errorf("TracerInsecure = %v, want true", opts.TracerInsecure)
	}
	if opts.MetricProvider != "otlp" {
		t.Errorf("MetricProvider = %v, want otlp", opts.MetricProvider)
	}
	if opts.MetricProviderHost != "localhost" {
		t.Errorf("MetricProviderHost = %v, want localhost", opts.MetricProviderHost)
	}
	if opts.MetricProviderPort != 4318 {
		t.Errorf("MetricProviderPort = %v, want 4318", opts.MetricProviderPort)
	}
	if opts.MetricInterval != 30*time.Second {
		t.Errorf("MetricInterval = %v, want 30s", opts.MetricInterval)
	}
	if opts.MetricInsecure != true {
		t.Errorf("MetricInsecure = %v, want true", opts.MetricInsecure)
	}
}

func TestMonitoring_Options_Override(t *testing.T) {
	opts := defaultOptions()

	WithServiceName("first-service")(opts)
	WithServiceName("second-service")(opts)
	if opts.ServiceName != "second-service" {
		t.Errorf("ServiceName = %v, want second-service", opts.ServiceName)
	}

	WithEnvironment("development")(opts)
	WithEnvironment("staging")(opts)
	WithEnvironment("production")(opts)
	if opts.Environment != "production" {
		t.Errorf("Environment = %v, want production", opts.Environment)
	}

	WithTracerSampleRatio(0.1)(opts)
	WithTracerSampleRatio(0.5)(opts)
	WithTracerSampleRatio(1.0)(opts)
	if opts.TracerSampleRatio != 1.0 {
		t.Errorf("TracerSampleRatio = %v, want 1.0", opts.TracerSampleRatio)
	}
}

func TestMonitoring_Options_Isolation(t *testing.T) {
	opts1 := defaultOptions()
	opts2 := defaultOptions()

	WithServiceName("service-1")(opts1)
	WithServiceName("service-2")(opts2)

	if opts1.ServiceName != "service-1" {
		t.Errorf("opts1.ServiceName = %v, want service-1", opts1.ServiceName)
	}
	if opts2.ServiceName != "service-2" {
		t.Errorf("opts2.ServiceName = %v, want service-2", opts2.ServiceName)
	}
	if opts1.ServiceName == opts2.ServiceName {
		t.Error("Options instances should be isolated from each other")
	}
}
