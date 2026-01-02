package logger

import "testing"

func TestLogger_NewLogger(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantErr bool
	}{
		{
			name:    "default logger",
			opts:    nil,
			wantErr: false,
		},
		{
			name:    "with valid log level",
			opts:    []Option{WithLevel("debug")},
			wantErr: false,
		},
		{
			name:    "with invalid log level",
			opts:    []Option{WithLevel("invalid")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loggerInstance, err := NewLogger(tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && loggerInstance == nil {
				t.Errorf("NewLogger() returned nil logger")
			}
		})
	}
}
