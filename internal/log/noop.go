package log

import "go.uber.org/zap"

func noopLogger() *Logger {
	return &Logger{
		SugaredLogger: zap.NewNop().Sugar(),
	}
}
