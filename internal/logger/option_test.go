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

func TestLogger_Option_WithOutputPath(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		checkFunc func(t *testing.T, opts *Options)
	}{
		{
			name: "set file path",
			path: "/var/log/app.log",
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.OutputPath != "/var/log/app.log" {
					t.Errorf("WithOutputPath() set path = %v, want %v", opts.OutputPath, "/var/log/app.log")
				}
			},
		},
		{
			name: "set relative path",
			path: "./logs/app.log",
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.OutputPath != "./logs/app.log" {
					t.Errorf("WithOutputPath() set path = %v, want %v", opts.OutputPath, "./logs/app.log")
				}
			},
		},
		{
			name: "set empty path",
			path: "",
			checkFunc: func(t *testing.T, opts *Options) {
				if opts.OutputPath != "" {
					t.Errorf("WithOutputPath() set path = %v, want %v", opts.OutputPath, "")
				}
			},
		},
		{
			name: "override existing path",
			path: "/new/path.log",
			checkFunc: func(t *testing.T, opts *Options) {
				opts.OutputPath = "/old/path.log"
				opt := WithOutputPath("/new/path.log")
				opt(opts)
				if opts.OutputPath != "/new/path.log" {
					t.Errorf("WithOutputPath() set path = %v, want %v", opts.OutputPath, "/new/path.log")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &Options{
				OutputPath: "",
			}
			opt := WithOutputPath(tt.path)
			opt(opts)
			if tt.checkFunc != nil {
				tt.checkFunc(t, opts)
			}
		})
	}
}
