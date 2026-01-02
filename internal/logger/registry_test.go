package logger

import (
	"testing"
)

func TestLogger_Registry_NewLogger(t *testing.T) {
	tests := []struct {
		name      string
		opts      []Option
		wantErr   bool
		checkFunc func(t *testing.T, logger Logger)
	}{
		{
			name:    "default logger",
			opts:    nil,
			wantErr: false,
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil logger")
				}
			},
		},
		{
			name:    "with debug level",
			opts:    []Option{WithLevel("debug")},
			wantErr: false,
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil logger")
				}
				logger.Debug("test message", nil)
			},
		},
		{
			name:    "with info level",
			opts:    []Option{WithLevel("info")},
			wantErr: false,
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil logger")
				}
				logger.Info("test message", nil)
			},
		},
		{
			name:    "with warn level",
			opts:    []Option{WithLevel("warn")},
			wantErr: false,
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil logger")
				}
				logger.Warn("test message", nil)
			},
		},
		{
			name:    "with error level",
			opts:    []Option{WithLevel("error")},
			wantErr: false,
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil logger")
				}
				logger.Error("test message", nil)
			},
		},
		{
			name:    "with fatal level",
			opts:    []Option{WithLevel("fatal")},
			wantErr: false,
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil logger")
				}
			},
		},
		{
			name:      "with invalid log level",
			opts:      []Option{WithLevel("invalid")},
			wantErr:   true,
			checkFunc: nil,
		},
		{
			name:    "with empty level",
			opts:    []Option{WithLevel("")},
			wantErr: false,
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil logger")
				}
			},
		},
		{
			name:    "with multiple options last wins",
			opts:    []Option{WithLevel("debug"), WithLevel("warn")},
			wantErr: false,
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil logger")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loggerInstance, err := NewLogger(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if loggerInstance != nil {
					t.Errorf("NewLogger() returned logger = %v, want nil", loggerInstance)
				}
				return
			}
			if loggerInstance == nil {
				t.Errorf("NewLogger() returned nil logger")
				return
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, loggerInstance)
			}
		})
	}
}
