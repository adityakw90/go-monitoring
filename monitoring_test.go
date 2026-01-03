package monitoring

import (
	"context"
	"testing"
	"time"
)

func TestMonitoring_Monitoring_Shutdown(t *testing.T) {
	monitoring, err := NewMonitoring(
		WithServiceName("test-service"),
	)
	if err != nil {
		t.Fatalf("NewMonitoring() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := monitoring.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}
}
