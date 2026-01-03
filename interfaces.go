package monitoring

import (
	"github.com/adityakw90/go-monitoring/internal/logger"
	"github.com/adityakw90/go-monitoring/internal/metric"
	"github.com/adityakw90/go-monitoring/internal/tracer"
)

// Logger is the interface for logging.
// It is re-exported from the internal logger package for public API use.
type Logger = logger.Logger

// Tracer is the interface for tracing.
// It is re-exported from the internal tracer package for public API use.
type Tracer = tracer.Tracer

// Metric is the interface for metrics.
// It is re-exported from the internal metric package for public API use.
type Metric = metric.Metric
