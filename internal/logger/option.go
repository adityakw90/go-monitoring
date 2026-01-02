package logger

type Options struct {
	Level string // Level is the minimum log level to output. Valid values: "debug", "info", "warn", "error", "fatal".
}

type Option func(*Options)

func WithLevel(level string) Option {
	return func(o *Options) {
		o.Level = level
	}
}
