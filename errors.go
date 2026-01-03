package monitoring

import (
	"errors"
	"fmt"

	"github.com/adityakw90/go-monitoring/internal/logger"
	"github.com/adityakw90/go-monitoring/internal/metric"
	"github.com/adityakw90/go-monitoring/internal/tracer"
)

// Error definitions for the monitoring library.
var (
	// ErrServiceNameRequired is returned when service name is not provided.
	ErrServiceNameRequired = errors.New("service name is required")
)

// re-export errors from internal packages
var (
	// logger
	ErrLoggerInvalidLogLevel = logger.ErrInvalidLogLevel

	// tracer
	ErrTracerInvalidProvider      = tracer.ErrInvalidProvider
	ErrTracerProviderHostRequired = tracer.ErrProviderHostRequired
	ErrTracerProviderPortRequired = tracer.ErrProviderPortRequired
	ErrTracerProviderPortInvalid  = tracer.ErrProviderPortInvalid
	ErrTracerBatchTimeoutInvalid  = tracer.ErrBatchTimeoutInvalid

	// metric
	ErrMetricInvalidProvider      = metric.ErrInvalidProvider
	ErrMetricProviderHostRequired = metric.ErrProviderHostRequired
	ErrMetricProviderPortRequired = metric.ErrProviderPortRequired
	ErrMetricProviderPortInvalid  = metric.ErrProviderPortInvalid
	ErrMetricIntervalInvalid      = metric.ErrIntervalInvalid
)

// parseError maps internal package errors to public API errors.
// It checks for known sentinel errors from internal packages and returns the corresponding
// parseError maps known internal sentinel errors to the package's public API error aliases.
// If err is nil it returns an error formatted as "<message>: unknown error".
// If err matches a recognized internal sentinel, it returns the corresponding exported error alias.
// For any other error it returns the original error wrapped with the provided message.
func parseError(err error, message string) error {
	if err == nil {
		return fmt.Errorf("%s: unknown error", message)
	}

	// logger
	if errors.Is(err, logger.ErrInvalidLogLevel) {
		return ErrLoggerInvalidLogLevel
	}

	// tracer
	if errors.Is(err, tracer.ErrInvalidProvider) {
		return ErrTracerInvalidProvider
	}
	if errors.Is(err, tracer.ErrProviderHostRequired) {
		return ErrTracerProviderHostRequired
	}
	if errors.Is(err, tracer.ErrProviderPortRequired) {
		return ErrTracerProviderPortRequired
	}
	if errors.Is(err, tracer.ErrProviderPortInvalid) {
		return ErrTracerProviderPortInvalid
	}
	if errors.Is(err, tracer.ErrBatchTimeoutInvalid) {
		return ErrTracerBatchTimeoutInvalid
	}

	// metric
	if errors.Is(err, metric.ErrInvalidProvider) {
		return ErrMetricInvalidProvider
	}
	if errors.Is(err, metric.ErrProviderHostRequired) {
		return ErrMetricProviderHostRequired
	}
	if errors.Is(err, metric.ErrProviderPortRequired) {
		return ErrMetricProviderPortRequired
	}
	if errors.Is(err, metric.ErrProviderPortInvalid) {
		return ErrMetricProviderPortInvalid
	}
	if errors.Is(err, metric.ErrIntervalInvalid) {
		return ErrMetricIntervalInvalid
	}

	return fmt.Errorf("%s: %w", message, err)
}