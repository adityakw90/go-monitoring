package logger

import (
	"errors"
	"testing"
)

func TestLogger_Error_ErrInvalidLogLevel(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "ErrInvalidLogLevel is defined",
			err:  ErrInvalidLogLevel,
			want: "invalid log level",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Error("ErrInvalidLogLevel is nil")
				return
			}
			if tt.err.Error() != tt.want {
				t.Errorf("ErrInvalidLogLevel.Error() = %v, want %v", tt.err.Error(), tt.want)
			}
			// Verify it's a standard error
			if !errors.Is(tt.err, ErrInvalidLogLevel) {
				t.Error("ErrInvalidLogLevel should match itself")
			}
		})
	}
}

func TestLogger_Error_ErrInvalidLogLevel_Usage(t *testing.T) {
	// Verify ErrInvalidLogLevel is used correctly in NewLogger
	loggerInstance, err := NewLogger(WithLevel("invalid_level_xyz"))
	if err == nil {
		t.Fatal("NewLogger() with invalid level expected error")
	}
	if err != ErrInvalidLogLevel {
		t.Errorf("NewLogger() error = %v, want ErrInvalidLogLevel", err)
	}
	if loggerInstance != nil {
		t.Error("NewLogger() with invalid level expected nil logger")
	}

	// Verify error can be compared
	if !errors.Is(err, ErrInvalidLogLevel) {
		t.Error("errors.Is() should return true for ErrInvalidLogLevel")
	}
}
