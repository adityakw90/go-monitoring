package monitoring

import (
	"errors"

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
	ErrMetricInvalidProvider = metric.ErrInvalidProvider
	ErrTracerInvalidProvider = tracer.ErrInvalidProvider
	ErrLoggerInvalidLogLevel = logger.ErrInvalidLogLevel
)
