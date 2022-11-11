// Package log contains the methods to configure the logger and to set and
// get the log mechanism from a context.
package log

import (
	"go.uber.org/zap"
)

// Logger provides a zap logger with extra wrap methods to facilitate operations
// involving the binaries related to Differential Evolution.
type Logger struct {
	*zap.SugaredLogger
}

// New returns a default logger.
func New() *Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return &Logger{
		SugaredLogger: logger.Sugar(),
	}
}
