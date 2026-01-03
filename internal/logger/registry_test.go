package logger

import (
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogger_Registry_NewLogger(t *testing.T) {
	tests := []struct {
		name        string
		opts        []Option
		wantErr     bool
		wantErrType error
		wantErrMsg  string
		checkFunc   func(t *testing.T, logger Logger)
	}{
		{
			name:        "default logger",
			opts:        nil,
			wantErr:     false,
			wantErrType: nil,
			wantErrMsg:  "",
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil logger")
				}
			},
		},
		{
			name:        "with debug level",
			opts:        []Option{WithLevel("debug")},
			wantErr:     false,
			wantErrType: nil,
			wantErrMsg:  "",
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil logger")
				}
				logger.Debug("test message", nil)
			},
		},
		{
			name:        "with info level",
			opts:        []Option{WithLevel("info")},
			wantErr:     false,
			wantErrType: nil,
			wantErrMsg:  "",
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil logger")
				}
				logger.Info("test message", nil)
			},
		},
		{
			name:        "with warn level",
			opts:        []Option{WithLevel("warn")},
			wantErr:     false,
			wantErrType: nil,
			wantErrMsg:  "",
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil logger")
				}
				logger.Warn("test message", nil)
			},
		},
		{
			name:        "with error level",
			opts:        []Option{WithLevel("error")},
			wantErr:     false,
			wantErrType: nil,
			wantErrMsg:  "",
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil logger")
				}
				logger.Error("test message", nil)
			},
		},
		{
			name:        "with fatal level",
			opts:        []Option{WithLevel("fatal")},
			wantErr:     false,
			wantErrType: nil,
			wantErrMsg:  "",
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil logger")
				}
			},
		},
		{
			name:        "with invalid log level",
			opts:        []Option{WithLevel("invalid")},
			wantErr:     true,
			wantErrType: ErrInvalidLogLevel,
			wantErrMsg:  "",
			checkFunc:   nil,
		},
		{
			name:        "with empty level, default to info",
			opts:        []Option{WithLevel("")},
			wantErr:     false,
			wantErrType: nil,
			wantErrMsg:  "",
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil logger")
				}
			},
		},
		{
			name:        "with multiple options last wins",
			opts:        []Option{WithLevel("debug"), WithLevel("warn")},
			wantErr:     false,
			wantErrType: nil,
			wantErrMsg:  "",
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil logger")
				}
			},
		},
		{
			name:        "with no options uses default",
			opts:        []Option{},
			wantErr:     false,
			wantErrType: nil,
			wantErrMsg:  "",
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil logger")
				}
				// Verify default level is info
				logger.Info("test", nil)
			},
		},
		{
			name:        "with nil option slice",
			opts:        nil,
			wantErr:     false,
			wantErrType: nil,
			wantErrMsg:  "",
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil logger")
				}
			},
		},
		{
			name:        "with output path",
			opts:        []Option{WithOutputPath("/tmp/test.log")},
			wantErr:     false,
			wantErrType: nil,
			wantErrMsg:  "",
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil logger")
				}
				logger.Info("TestLogger_Registry_NewLogger/with output path: should be info and written to file", nil)
				// check if the log file is created
				_, err := os.Stat("/tmp/test.log")
				assert.NoError(t, err)
				assert.FileExists(t, "/tmp/test.log")
				defer os.Remove("/tmp/test.log") // clean up the log file
				// read the log file
				content, err := os.ReadFile("/tmp/test.log")
				assert.NoError(t, err)
				assert.NotEmpty(t, content)
				// read the first line and must be json format
				var logEntry map[string]interface{}
				err = json.Unmarshal(content, &logEntry)
				assert.NoError(t, err)
				assert.NotEmpty(t, logEntry)
				assert.Contains(t, logEntry, "level")
				assert.Contains(t, logEntry, "ts")
				assert.Contains(t, logEntry, "caller")
				assert.Contains(t, logEntry, "msg")
				assert.Equal(t, "info", logEntry["level"], "level should be info")
				assert.NotEmpty(t, logEntry["ts"], "ts should not be empty")
				_, err = time.Parse("2006-01-02T15:04:05.000-0700", logEntry["ts"].(string))
				assert.NoError(t, err, "ts should be timestamp in the correct format")
				assert.Contains(t, logEntry["caller"], "logger/registry_test.go")
				assert.Equal(t, "TestLogger_Registry_NewLogger/with output path: should be info and written to file", logEntry["msg"])
			},
		},
		{
			name:        "with output path and level",
			opts:        []Option{WithLevel("debug"), WithOutputPath("/tmp/test.log")},
			wantErr:     false,
			wantErrType: nil,
			wantErrMsg:  "",
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil logger")
				}
				logger.Debug("TestLogger_Registry_NewLogger/with output path and level : should be debug and written to file", nil)
				// check if the log file is created
				_, err := os.Stat("/tmp/test.log")
				assert.NoError(t, err)
				assert.FileExists(t, "/tmp/test.log")
				defer os.Remove("/tmp/test.log") // clean up the log file
				// read the log file
				content, err := os.ReadFile("/tmp/test.log")
				assert.NoError(t, err)
				assert.NotEmpty(t, content)
				// read the first line and must be json format
				var logEntry map[string]interface{}
				err = json.Unmarshal(content, &logEntry)
				assert.NoError(t, err)
				assert.NotEmpty(t, logEntry)
				assert.Contains(t, logEntry, "level")
				assert.Contains(t, logEntry, "ts")
				assert.Contains(t, logEntry, "caller")
				assert.Contains(t, logEntry, "msg")
				assert.Equal(t, "debug", logEntry["level"], "level should be debug")
				assert.NotEmpty(t, logEntry["ts"], "ts should not be empty")
				_, err = time.Parse("2006-01-02T15:04:05.000-0700", logEntry["ts"].(string))
				assert.NoError(t, err, "ts should be timestamp in the correct format")
				assert.Contains(t, logEntry["caller"], "logger/registry_test.go")
				assert.Equal(t, "TestLogger_Registry_NewLogger/with output path and level : should be debug and written to file", logEntry["msg"])
			},
		},
		{
			name:        "with empty output path",
			opts:        []Option{WithOutputPath("")},
			wantErr:     false,
			wantErrType: nil,
			wantErrMsg:  "",
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil logger")
				}
				// Empty output path should default to stdout
				logger.Info("TestLogger_Registry_NewLogger/with empty output path : should be info and written to stdout", nil)
			},
		},
		{
			name:        "with relative output path",
			opts:        []Option{WithOutputPath("./test.log")},
			wantErr:     false,
			wantErrType: nil,
			wantErrMsg:  "",
			checkFunc: func(t *testing.T, logger Logger) {
				if logger == nil {
					t.Error("NewLogger() returned nil logger")
				}
				logger.Info("TestLogger_Registry_NewLogger/with relative output path : should be info and written to file", nil)
				// check if the log file is created
				_, err := os.Stat("./test.log")
				assert.NoError(t, err)
				assert.FileExists(t, "./test.log")
				defer os.Remove("./test.log") // clean up the log file
				// read the log file
				content, err := os.ReadFile("./test.log")
				assert.NoError(t, err)
				assert.NotEmpty(t, content)
				// read the first line and must be json format
				var logEntry map[string]interface{}
				err = json.Unmarshal(content, &logEntry)
				assert.NoError(t, err)
				assert.NotEmpty(t, logEntry)
				assert.Contains(t, logEntry, "level")
				assert.Contains(t, logEntry, "ts")
				assert.Contains(t, logEntry, "caller")
				assert.Contains(t, logEntry, "msg")
				assert.Equal(t, "info", logEntry["level"], "level should be info")
				assert.NotEmpty(t, logEntry["ts"], "ts should not be empty")
				_, err = time.Parse("2006-01-02T15:04:05.000-0700", logEntry["ts"].(string))
				assert.NoError(t, err, "ts should be timestamp in the correct format")
				assert.Contains(t, logEntry["caller"], "logger/registry_test.go")
				assert.Equal(t, "TestLogger_Registry_NewLogger/with relative output path : should be info and written to file", logEntry["msg"])
			},
		},
		{
			name:        "with unexisting output path",
			opts:        []Option{WithOutputPath("./this/path/does/not/exist/log.json")},
			wantErr:     true,
			wantErrType: nil,
			wantErrMsg:  "failed to build logger: couldn't open sink \"./this/path/does/not/exist/log.json\": open ./this/path/does/not/exist/log.json: no such file or directory",
			checkFunc:   nil,
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
				assert.Error(t, err)
				assert.Nil(t, loggerInstance)
				if tt.wantErrType != nil {
					assert.True(t, errors.Is(err, tt.wantErrType))
				}
				if tt.wantErrMsg != "" {
					assert.Equal(t, tt.wantErrMsg, err.Error())
				}
			} else {
				if loggerInstance == nil {
					t.Errorf("NewLogger() returned nil logger")
					return
				}
				if tt.checkFunc != nil {
					tt.checkFunc(t, loggerInstance)
				}
			}
		})
	}
}
