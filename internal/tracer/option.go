package tracer

import "time"

// Options contains configuration options for creating a Tracer.
// All fields are optional and have sensible defaults.
type Options struct {
	ServiceName  string        // ServiceName is the name of the service being traced.
	Environment  string        // Environment is the deployment environment (e.g., "development", "production").
	InstanceName string        // InstanceName is the unique identifier for this service instance.
	InstanceHost string        // InstanceHost is the hostname where this service instance is running.
	Provider     string        // Provider specifies the trace exporter to use ("stdout" or "otlp").
	ProviderHost string        // ProviderHost is the hostname of the OTLP trace collector (only used when Provider is "otlp").
	ProviderPort int           // ProviderPort is the port of the OTLP trace collector (only used when Provider is "otlp").
	SampleRatio  float64       // SampleRatio controls the sampling rate for traces (0.0 to 1.0). 0.0 means never sample, 1.0 means always sample, values in between use probabilistic sampling.
	BatchTimeout time.Duration // BatchTimeout is the maximum time to wait before exporting a batch of spans.
	Insecure     bool          // Insecure controls whether to use an insecure (non-TLS) connection for OTLP exporter. When true, connections are made without TLS. Default is false (secure TLS connection).
}

// Option is a function that configures Options.
// It follows the functional options pattern for flexible tracer configuration.
type Option func(*Options)

// WithServiceName sets the service name for the tracer (internal use).
func WithServiceName(name string) Option {
	return func(o *Options) {
		o.ServiceName = name
	}
}

// WithEnvironment sets the environment (internal use).
func WithEnvironment(env string) Option {
	return func(o *Options) {
		o.Environment = env
	}
}

// WithInstance sets the instance name and host (internal use).
func WithInstance(name, host string) Option {
	return func(o *Options) {
		o.InstanceName = name
		o.InstanceHost = host
	}
}

// WithProvider sets the tracer provider configuration (internal use).
func WithProvider(provider, host string, port int) Option {
	return func(o *Options) {
		o.Provider = provider
		o.ProviderHost = host
		o.ProviderPort = port
	}
}

// WithSampleRatio sets the sampling ratio (internal use).
func WithSampleRatio(ratio float64) Option {
	return func(o *Options) {
		o.SampleRatio = ratio
	}
}

// WithBatchTimeout sets the batch timeout (internal use).
func WithBatchTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.BatchTimeout = timeout
	}
}

// WithInsecure sets whether to use an insecure connection for OTLP exporter (internal use).
func WithInsecure(insecure bool) Option {
	return func(o *Options) {
		o.Insecure = insecure
	}
}
