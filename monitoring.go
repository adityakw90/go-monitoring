package monitoring

import (
	"context"
	"fmt"
)

// Monitoring contains all observability components in a single unified structure.
// It provides access to logging, tracing, and metrics functionality.
type Monitoring struct {
	Logger Logger // Logger provides structured logging capabilities.
	Tracer Tracer // Tracer provides distributed tracing capabilities.
	Metric Metric // Metric provides metrics collection capabilities.
}

// Shutdown gracefully shuts down all monitoring components.
// It shuts down the Tracer and Metric providers in order, ensuring all
// pending traces and metrics are exported before termination.
//
// This should be called before application shutdown to ensure proper cleanup.
// The Logger does not require explicit shutdown.
//
// Parameters:
//   - ctx: Context for controlling shutdown timeout
//
// Returns an error if shutdown of any component fails.
// Errors from individual components are wrapped with context.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	if err := mon.Shutdown(ctx); err != nil {
//	    log.Printf("Failed to shutdown monitoring: %v", err)
//	}
func (m *Monitoring) Shutdown(ctx context.Context) error {
	if m.Tracer != nil {
		if err := m.Tracer.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown tracer: %w", err)
		}
	}
	if m.Metric != nil {
		if err := m.Metric.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown metric: %w", err)
		}
	}
	return nil
}
