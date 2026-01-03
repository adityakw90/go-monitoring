package logger

type Options struct {
	Level      string // Level is the minimum log level to output. Valid values: "debug", "info", "warn", "error", "fatal".
	OutputPath string // OutputPath is the file path where logs will be written. If empty, logs will be written to stdout.
}

type Option func(*Options)

func WithLevel(level string) Option {
	return func(o *Options) {
		o.Level = level
	}
}

func WithOutputPath(path string) Option {
	return func(o *Options) {
		o.OutputPath = path
	}
}
