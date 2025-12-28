package monitoring

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "ErrServiceNameRequired",
			err:  ErrServiceNameRequired,
			want: "service name is required",
		},
		{
			name: "ErrInvalidLogLevel",
			err:  ErrInvalidLogLevel,
			want: "invalid log level",
		},
		{
			name: "ErrInvalidProvider",
			err:  ErrInvalidProvider,
			want: "invalid provider",
		},
		{
			name: "ErrInvalidSampleRatio",
			err:  ErrInvalidSampleRatio,
			want: "sample ratio must be between 0 and 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test error message
			if tt.err.Error() != tt.want {
				t.Errorf("error message = %q, want %q", tt.err.Error(), tt.want)
			}

			// Test that wrapped error reports the sentinel via errors.Is
			wrappedErr := fmt.Errorf("operation failed: %w", tt.err)
			if !errors.Is(wrappedErr, tt.err) {
				t.Errorf("errors.Is(wrappedErr, sentinel) = false, want true")
			}

			// Test that unrelated errors return false
			unrelatedErr := errors.New("unrelated error")
			if errors.Is(unrelatedErr, tt.err) {
				t.Errorf("errors.Is(unrelatedErr, sentinel) = true, want false")
			}
		})
	}
}
