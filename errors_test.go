package monitoring

import (
	"errors"
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
			if tt.err.Error() != tt.want {
				t.Errorf("error message = %q, want %q", tt.err.Error(), tt.want)
			}
			if !errors.Is(tt.err, tt.err) {
				t.Errorf("errors.Is() = false, want true")
			}
		})
	}
}
