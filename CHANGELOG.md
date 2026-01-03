# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

> **⚠️ Pre-1.0.0 Notice**: Versions before 1.0.0 are considered unstable. **Backward compatibility is not guaranteed** between minor versions (0.x → 0.y). Breaking changes may be introduced without major version bumps. Once version 1.0.0 is released, this project will follow semantic versioning strictly, and breaking changes will only occur in major version updates.

## [Unreleased]

### Added
- Comprehensive code review documentation
- Troubleshooting section in README
- Contributing guidelines

### Changed
- Updated README with version information
- Fixed documentation inconsistencies

## [0.1.0] - 2025-01-27

### Added
- Initial release of go-monitoring library
- Structured logging with Zap integration
- Distributed tracing with OpenTelemetry
- Metrics collection with OpenTelemetry
- Unified Monitoring struct for all components
- Functional options pattern for configuration
- Support for stdout and OTLP exporters
- Trace context propagation for gRPC
- Runtime log level adjustment
- Comprehensive test coverage (97.6%)

### Features
- **Logger**: JSON-structured logging with trace context integration
- **Tracer**: OpenTelemetry tracing with configurable sampling
- **Metric**: Counter and histogram metrics with OpenTelemetry
- **Configuration**: Flexible options pattern with sensible defaults

### Documentation
- Comprehensive README with examples
- GoDoc comments on all exported functions
- API reference documentation
- Usage examples for common scenarios

---

[Unreleased]: https://github.com/adityakw90/go-monitoring/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/adityakw90/go-monitoring/releases/tag/v0.1.0

