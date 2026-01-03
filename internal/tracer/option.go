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

// WithServiceName returns an Option that sets the tracer service name.
// The provided name is applied to Options.ServiceName.
func WithServiceName(name string) Option {
	return func(o *Options) {
		o.ServiceName = name
	}
}

// WithEnvironment returns an Option that sets the tracer's Environment field.
// The value typically identifies the deployment environment, e.g. "development" or "production".
func WithEnvironment(env string) Option {
	return func(o *Options) {
		o.Environment = env
	}
}

// WithInstance sets the tracer instance name and host.
// name is the instance identifier; host is the instance hostname.
func WithInstance(name, host string) Option {
	return func(o *Options) {
		o.InstanceName = name
		o.InstanceHost = host
	}
}

// collector endpoint to use when the provider requires a network collector.
func WithProvider(provider, host string, port int) Option {
	return func(o *Options) {
		o.Provider = provider
		o.ProviderHost = host
		o.ProviderPort = port
	}
}

// WithSampleRatio returns an Option that sets the tracer sampling ratio.
// Valid values are between 0.0 and 1.0 inclusive — 0.0 means never sample and 1.0 means always sample.
func WithSampleRatio(ratio float64) Option {
	return func(o *Options) {
		o.SampleRatio = ratio
	}
}

// WithBatchTimeout returns an Option that sets the maximum time to wait before exporting a batch of spans.
func WithBatchTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.BatchTimeout = timeout
	}
}

// WithInsecure sets whether the OTLP exporter uses an insecure (non‑TLS) connection.
func WithInsecure(insecure bool) Option {
	return func(o *Options) {
		o.Insecure = insecure
	}
}