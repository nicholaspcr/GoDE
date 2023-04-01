package log

import "context"

// logkey defines an empty struct to be used as a key in the context.
type logKey struct{}

var loggerKey = &logKey{}

// SetContext creates a new context with the logger.
func SetContext(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// FromContext returns a logger from a context.
func FromContext(ctx context.Context) *Logger {
	logger := ctx.Value(loggerKey)
	if logger == nil {
		return noopLogger()
	}
	return logger.(*Logger)
}
