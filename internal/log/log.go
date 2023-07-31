// Package log contains the methods to configure the logger and to set and
// get the log mechanism from a context.
package log

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger provides a zap logger with extra wrap methods to facilitate operations
// involving the binaries related to Differential Evolution.
type Logger struct {
	*zap.Logger
}

// New returns a default logger.
func New() *Logger {
	return &Logger{
		Logger: zap.New(
			zapcore.NewCore(
				getEncoder(),
				zapcore.Lock(os.Stdout),
				getLogLevel(),
			),
		),
	}
}

func getEncoder() zapcore.Encoder {
	enconder := strings.ToLower(os.Getenv("GODE_LOG_ENCODER"))
	switch enconder {
	case "console":
		return zapcore.NewConsoleEncoder(getEncoderConfig())
	default:
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
}
func getEncoderConfig() zapcore.EncoderConfig {
	config := strings.ToLower(os.Getenv("GODE_LOG_ENCODER_CONFIG"))
	var cfg zapcore.EncoderConfig
	switch config {
	case "development":
		cfg = zap.NewDevelopmentEncoderConfig()
	default:
		cfg = zap.NewProductionEncoderConfig()
	}
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return cfg
}

func getLogLevel() zapcore.Level {
	logLevel := strings.ToLower(os.Getenv("GODE_LOG_LEVEL"))
	switch logLevel {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
