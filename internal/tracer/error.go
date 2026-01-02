package tracer

import "errors"

var (
	// ErrInvalidProvider is returned when an invalid provider type is specified.
	ErrInvalidProvider = errors.New("invalid provider")
)
