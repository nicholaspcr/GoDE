package log

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// This file aggregates the zapcore field methods to be a part of the log
// package, so that the user doesn't need to import zapcore to use them.

func Error(k string, v error) zapcore.Field {
	return zap.Error(v)
}

func String(k string, v string) zapcore.Field {
	return zap.String(k, v)
}

func Int(k string, v int) zapcore.Field {
	return zap.Int(k, v)
}

func Float64(k string, v float64) zapcore.Field {
	return zap.Float64(k, v)
}

func Bool(k string, v bool) zapcore.Field {
	return zap.Bool(k, v)
}

func Duration(k string, v time.Duration) zapcore.Field {
	return zap.Duration(k, v)
}

func Time(k string, v time.Time) zapcore.Field {
	return zap.Time(k, v)
}

func Object(k string, v zapcore.ObjectMarshaler) zapcore.Field {
	return zap.Object(k, v)
}

func Array(k string, v zapcore.ArrayMarshaler) zapcore.Field {
	return zap.Array(k, v)
}

func Binary(k string, v []byte) zapcore.Field {
	return zap.Binary(k, v)
}

func Stringer(k string, v fmt.Stringer) zapcore.Field {
	return zap.Stringer(k, v)
}

func Any(k string, v any) zapcore.Field {
	return zap.Any(k, v)
}

func Namespace(k string) zapcore.Field {
	return zap.Namespace(k)
}
