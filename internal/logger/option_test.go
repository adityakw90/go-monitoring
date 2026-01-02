package logger

import (
	"testing"
)

func TestLogger_Option_WithLevel(t *testing.T) {
	tests := []struct {
		name      string
		level     string
		checkFunc func(t *testing.T, opts *Options)
	}{
		{
			name:  "set debug level",
			level: "debug",
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.Level != "debug" {
					t.Errorf("WithLevel() set level = %v, want %v", opts.Level, "debug")
				}
			},
		},
		{
			name:  "set info level",
			level: "info",
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.Level != "info" {
					t.Errorf("WithLevel() set level = %v, want %v", opts.Level, "info")
				}
			},
		},
		{
			name:  "set warn level",
			level: "warn",
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.Level != "warn" {
					t.Errorf("WithLevel() set level = %v, want %v", opts.Level, "warn")
				}
			},
		},
		{
			name:  "set error level",
			level: "error",
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.Level != "error" {
					t.Errorf("WithLevel() set level = %v, want %v", opts.Level, "error")
				}
			},
		},
		{
			name:  "set fatal level",
			level: "fatal",
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.Level != "fatal" {
					t.Errorf("WithLevel() set level = %v, want %v", opts.Level, "fatal")
				}
			},
		},
		{
			name:  "set empty level",
			level: "",
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.Level != "" {
					t.Errorf("WithLevel() set level = %v, want %v", opts.Level, "")
				}
			},
		},
		{
			name:  "set invalid level",
			level: "invalid",
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.Level != "invalid" {
					t.Errorf("WithLevel() set level = %v, want %v", opts.Level, "invalid")
				}
			},
		},
		{
			name:  "override existing level",
			level: "warn",
			checkFunc: func(t *testing.T, opts *Options) {
				opts.Level = "info"
				opt := WithLevel("warn")
				opt(opts)
				if opts.Level != "warn" {
					t.Errorf("WithLevel() set level = %v, want %v", opts.Level, "warn")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &Options{
				Level: "info",
			}
			opt := WithLevel(tt.level)
			opt(opts)
			if tt.checkFunc != nil {
				tt.checkFunc(t, opts)
			}
		})
	}
}
