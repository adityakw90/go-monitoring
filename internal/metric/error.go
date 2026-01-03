package metric

import "errors"

var (
	// ErrInvalidProvider is returned when an invalid provider type is specified.
	ErrInvalidProvider      = errors.New("invalid provider")
	ErrProviderHostRequired = errors.New("provider host is required")
	ErrProviderPortRequired = errors.New("provider port is required")
	ErrProviderPortInvalid  = errors.New("provider port must be greater than 0")
	ErrIntervalInvalid      = errors.New("interval must be greater than 0")
)
