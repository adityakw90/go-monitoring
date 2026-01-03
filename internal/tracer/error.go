package tracer

import "errors"

var (
	// ErrInvalidProvider is returned when an invalid provider type is specified.
	ErrInvalidProvider      = errors.New("invalid provider")
	ErrProviderHostRequired = errors.New("provider host is required")
	ErrProviderPortRequired = errors.New("provider port is required")
	ErrProviderPortInvalid  = errors.New("provider port must be greater than 0")
	ErrBatchTimeoutInvalid  = errors.New("batch timeout must be greater than 0")
)
