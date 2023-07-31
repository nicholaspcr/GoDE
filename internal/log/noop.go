package log

import "go.uber.org/zap"

func noopLogger() *Logger {
	return &Logger{
		Logger: zap.NewNop(),
	}
}
