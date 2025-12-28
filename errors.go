package monitoring

import "errors"

// Error definitions for the monitoring library.
var (
	// ErrServiceNameRequired is returned when service name is not provided.
	ErrServiceNameRequired = errors.New("service name is required")

	// ErrInvalidLogLevel is returned when an invalid log level is provided.
	ErrInvalidLogLevel = errors.New("invalid log level")

	// ErrInvalidProvider is returned when an invalid provider type is specified.
	ErrInvalidProvider = errors.New("invalid provider")

	// ErrInvalidSampleRatio is returned when sample ratio is not between 0 and 1.
	ErrInvalidSampleRatio = errors.New("sample ratio must be between 0 and 1")
)
