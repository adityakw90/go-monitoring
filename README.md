# Go Monitoring Library

[![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A comprehensive observability library for Go that provides structured logging, distributed tracing, and metrics collection. This library consolidates monitoring functionality into a single, reusable package.

**Version**: v0.2.0

> **⚠️ Pre-1.0.0 Notice**: This library is currently in version 0.x.x. **Backward compatibility is not guaranteed** until version 1.0.0 is released. Breaking changes may occur in minor version updates (0.x → 0.y). Please pin your dependency to a specific version or use a dependency management tool that locks versions.

## Features

- **Structured Logging**: Zap-based JSON logging with trace context integration
- **Distributed Tracing**: OpenTelemetry tracing with OTLP and stdout exporters
- **Metrics Collection**: OpenTelemetry metrics with counter and histogram support
- **Unified Initialization**: Single entry point for all monitoring components
- **Functional Options**: Flexible configuration using the options pattern
- **Service-Independent**: No service-specific dependencies

## Installation

```bash
go get github.com/adityakw90/go-monitoring
```

> **⚠️ Version Compatibility**: For production use before v1.0.0, pin to a specific version:
> ```bash
> go get github.com/adityakw90/go-monitoring@v0.2.0
> ```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "time"
    
    "github.com/adityakw90/go-monitoring"
)

func main() {
    // Initialize monitoring with required service name
    mon, err := monitoring.NewMonitoring(
        monitoring.WithServiceName("my-service"),
        monitoring.WithEnvironment("production"),
        monitoring.WithInstance("instance-1", "localhost"),
    )
    if err != nil {
        panic(err)
    }
    defer mon.Shutdown(context.Background())

    // Use logger
    mon.Logger.Info("Service started", map[string]interface{}{
        "port": 8080,
    })

    // Use tracer
    ctx, span := mon.Tracer.StartSpan(context.Background(), "operation")
    defer span.End()

    // Use metrics
    counter, _ := mon.Metric.CreateCounter("requests_total", "1", "Total requests")
    mon.Metric.RecordCounter(ctx, counter, 1)
}
```

### Advanced Configuration

```go
mon, err := monitoring.NewMonitoring(
    monitoring.WithServiceName("my-service"),
    monitoring.WithEnvironment("production"),
    monitoring.WithInstance("instance-1", "localhost"),
    monitoring.WithLoggerLevel("debug"),
    monitoring.WithTracerProvider("otlp", "localhost", 4317),
    monitoring.WithTracerSampleRatio(0.1),
    monitoring.WithMetricProvider("otlp", "localhost", 4318),
    monitoring.WithMetricInterval(30 * time.Second),
)
```

### Individual Components

You can also initialize components individually:

```go
// Logger
logger, err := monitoring.NewLogger(
    monitoring.WithLoggerLevel("info"),
)

// Tracer
tracer, err := monitoring.NewTracer(
    monitoring.WithServiceName("my-service"),
    monitoring.WithTracerProvider("stdout", "", 0),
)

// Metric
metric, err := monitoring.NewMetric(
    monitoring.WithServiceName("my-service"),
    monitoring.WithMetricProvider("stdout", "", 0),
)
```

## API Reference

### Monitoring

#### `NewMonitoring(opts ...Option) (*Monitoring, error)`

Initializes all monitoring components (Logger, Tracer, Metric) with the given options.

**Required Options:**
- `WithServiceName(name string)` - Service name (required)

**Optional Options:**
- `WithEnvironment(env string)` - Environment (default: "development")
- `WithInstance(name, host string)` - Instance name and host
- `WithLoggerLevel(level string)` - Log level (default: "info")
- `WithTracerProvider(provider, host string, port int)` - Tracer provider (default: "stdout")
- `WithTracerSampleRatio(ratio float64)` - Sampling ratio 0.0-1.0 (default: 1.0)
- `WithMetricProvider(provider, host string, port int)` - Metric provider (default: "stdout")
- `WithMetricInterval(interval time.Duration)` - Export interval (default: 60s)

### Logger

The Logger provides structured logging with Zap.

**Methods:**
- `Debug(message string, fields map[string]interface{})`
- `Info(message string, fields map[string]interface{})`
- `Warn(message string, fields map[string]interface{})`
- `Error(message string, fields map[string]interface{})`
- `Fatal(message string, fields map[string]interface{})`
- `SetLogLevel(level string)` - Change log level at runtime (invalid levels default to INFO)
- `WithSpanContext(span trace.SpanContext) *Logger` - Add trace context to logs

### Tracer

The Tracer provides distributed tracing with OpenTelemetry.

**Methods:**
- `StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span)`
- `EndSpan(span trace.Span)`
- `Shutdown(ctx context.Context) error`
- `ExtractContext(ctx context.Context, md metadata.MD) context.Context` - Extract from gRPC metadata
- `InjectContext(ctx context.Context) metadata.MD` - Inject into gRPC metadata

### Metric

The Metric provides metrics collection with OpenTelemetry.

**Methods:**
- `CreateCounter(name, unit, description string) (metric.Int64Counter, error)`
- `RecordCounter(ctx context.Context, counter metric.Int64Counter, value int64, labels ...attribute.KeyValue)`
- `CreateHistogram(name, unit, description string) (metric.Int64Histogram, error)`
- `RecordHistogram(ctx context.Context, histogram metric.Int64Histogram, value int64, labels ...attribute.KeyValue)`
- `CreateAttributeInt(key string, value int) attribute.KeyValue`
- `CreateAttributeString(key string, value string) attribute.KeyValue`
- `Shutdown(ctx context.Context) error`

## Examples

### Logging with Trace Context

```go
ctx, span := mon.Tracer.StartSpan(ctx, "handle-request")
defer span.End()

// Add trace context to logger
logger := mon.Logger.WithSpanContext(span.SpanContext())
logger.Info("Processing request", map[string]interface{}{
    "request_id": "123",
})
```

### Recording Metrics

```go
// Create counter
counter, _ := mon.Metric.CreateCounter(
    "http_requests_total",
    "1",
    "Total HTTP requests",
)

// Record metric
mon.Metric.RecordCounter(ctx, counter, 1,
    mon.Metric.CreateAttributeString("method", "GET"),
    mon.Metric.CreateAttributeString("status", "200"),
)

// Create histogram
histogram, _ := mon.Metric.CreateHistogram(
    "http_request_duration_ms",
    "ms",
    "HTTP request duration",
)

// Record duration
mon.Metric.RecordHistogram(ctx, histogram, 150,
    mon.Metric.CreateAttributeString("method", "GET"),
)
```

### gRPC Context Propagation

```go
// Server: Extract context
ctx := mon.Tracer.ExtractContext(ctx, md)

// Client: Inject context
md := mon.Tracer.InjectContext(ctx)
```

## Configuration

### Log Levels

Supported log levels: `debug`, `info`, `warn`, `error`, `fatal`

### Tracer Providers

- `stdout` - Output traces to stdout (for development)
- `otlp` - Send traces via OTLP/gRPC

### Metric Providers

- `stdout` - Output metrics to stdout (for development)
- `otlp` - Send metrics via OTLP/gRPC

## Troubleshooting

### Common Issues

#### Traces not appearing in OTLP collector
- **Check connection**: Ensure the OTLP collector is running and accessible
- **Verify endpoint**: Check that the host and port are correct
- **Check TLS**: If using TLS, ensure certificates are properly configured
- **Sample ratio**: Verify `TracerSampleRatio` is not set to 0.0

#### Metrics not exporting
- **Check interval**: Metrics are exported periodically (default: 60s). Wait for the interval to pass
- **Verify provider**: Ensure the metric provider is correctly configured
- **Check OTLP connection**: Same as trace troubleshooting above

#### Logs not appearing in file
- **File permissions**: Ensure the application has write permissions to the log file path
- **Path validation**: Verify the file path is valid and the directory exists
- **Log level**: Check that the log level allows the messages you're trying to log

#### Invalid log level error
- **Valid levels**: Use only `debug`, `info`, `warn`, `error`, or `fatal` (case-insensitive)
- **SetLogLevel behavior**: Invalid levels in `SetLogLevel` will default to INFO and log a warning

### Performance Considerations

- **High-frequency logging**: For applications with very high log volume, consider using async logging or adjusting log levels
- **Trace sampling**: Use `TracerSampleRatio` < 1.0 in production to reduce overhead
- **Metric intervals**: Adjust `MetricInterval` based on your needs (shorter = more real-time but higher overhead)

## Requirements

- Go 1.22 or later
- OpenTelemetry SDK v1.34.0
- Zap v1.21.0

## Version History

See [CHANGELOG.md](CHANGELOG.md) for detailed version history and changes.

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Quick Contribution Guide

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`make test`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Development Setup

```bash
# Clone the repository
git clone https://github.com/adityakw90/go-monitoring.git
cd go-monitoring

# Run tests
make test

# Run tests with coverage
make test-cover

# Clean test cache
make test-clean
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
