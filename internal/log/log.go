// Package log contains the methods to configure the logger and to set and
// get the log mechanism from a context.
package log

import (
	"context"

	"go.uber.org/zap"
)

// Logger provides a zap logger with extra wrap methods to facilitate operations
// involving the binaries related to Differential Evolution.
type Logger struct {
	*zap.SugaredLogger
}

// logkey defines an empty struct to be used as a key in the context.
type logKey struct{}

var loggerKey = &logKey{}

// New returns a default logger.
func New() Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return Logger{
		SugaredLogger: logger.Sugar(),
	}
}

// IsNil informs if the inner logger is instantiated or not.
func (l *Logger) IsNil() bool {
	return l.SugaredLogger == nil
}

// SetContext creates a new context with the logger.
func (l *Logger) SetContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// SetContext creates a new context with the logger.
func SetContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

// FromContext returns a logger from a context.
func FromContext(ctx context.Context) Logger {
	logger := ctx.Value(loggerKey)
	if logger == nil {
		return noopLogger()
	}
	return logger.(Logger)
}
