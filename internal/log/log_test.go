package log

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_DefaultConfig(t *testing.T) {
	logger := New()
	require.NotNil(t, logger)
}

func TestNew_JSONHandler(t *testing.T) {
	var buf bytes.Buffer
	logger := New(
		WithWriter(&buf),
		WithType("json"),
		WithPrettyConfig(&PrettyConfig{Enable: false}),
	)

	logger.Info("test message", slog.String("key", "value"))

	output := buf.String()
	assert.Contains(t, output, "test message")
	assert.Contains(t, output, "key")
	assert.Contains(t, output, "value")
}

func TestNew_TextHandler(t *testing.T) {
	var buf bytes.Buffer
	logger := New(
		WithWriter(&buf),
		WithType("text"),
		WithPrettyConfig(&PrettyConfig{Enable: false}),
	)

	logger.Info("text message", slog.String("name", "test"))

	output := buf.String()
	assert.Contains(t, output, "text message")
	assert.Contains(t, output, "name=test")
}

func TestNew_LogLevels(t *testing.T) {
	tests := []struct {
		name     string
		level    slog.Level
		logLevel slog.Level
		visible  bool
	}{
		{
			name:     "debug message at info level",
			level:    slog.LevelInfo,
			logLevel: slog.LevelDebug,
			visible:  false,
		},
		{
			name:     "info message at info level",
			level:    slog.LevelInfo,
			logLevel: slog.LevelInfo,
			visible:  true,
		},
		{
			name:     "warn message at info level",
			level:    slog.LevelInfo,
			logLevel: slog.LevelWarn,
			visible:  true,
		},
		{
			name:     "error message at info level",
			level:    slog.LevelInfo,
			logLevel: slog.LevelError,
			visible:  true,
		},
		{
			name:     "debug message at debug level",
			level:    slog.LevelDebug,
			logLevel: slog.LevelDebug,
			visible:  true,
		},
		{
			name:     "info message at error level",
			level:    slog.LevelError,
			logLevel: slog.LevelInfo,
			visible:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := New(
				WithWriter(&buf),
				WithLevel(tt.level),
				WithPrettyConfig(&PrettyConfig{Enable: false}),
			)

			switch tt.logLevel {
			case slog.LevelDebug:
				logger.Debug("test")
			case slog.LevelInfo:
				logger.Info("test")
			case slog.LevelWarn:
				logger.Warn("test")
			case slog.LevelError:
				logger.Error("test")
			}

			if tt.visible {
				assert.NotEmpty(t, buf.String(), "message should be visible")
			} else {
				assert.Empty(t, buf.String(), "message should be filtered")
			}
		})
	}
}

func TestNew_PrettyHandler(t *testing.T) {
	t.Run("pretty with color", func(t *testing.T) {
		var buf bytes.Buffer
		logger := New(
			WithWriter(&buf),
			WithPrettyConfig(&PrettyConfig{
				Enable:     true,
				Color:      true,
				TimeFormat: "15:04:05",
			}),
		)

		logger.Info("pretty message", slog.String("key", "val"))
		output := buf.String()
		assert.Contains(t, output, "pretty message")
	})

	t.Run("pretty without color", func(t *testing.T) {
		var buf bytes.Buffer
		logger := New(
			WithWriter(&buf),
			WithPrettyConfig(&PrettyConfig{
				Enable:     true,
				Color:      false,
				TimeFormat: "15:04:05",
			}),
		)

		logger.Info("no color message")
		output := buf.String()
		assert.Contains(t, output, "no color message")
	})

	t.Run("pretty with indent", func(t *testing.T) {
		var buf bytes.Buffer
		logger := New(
			WithWriter(&buf),
			WithPrettyConfig(&PrettyConfig{
				Enable:     true,
				Color:      false,
				TimeFormat: "15:04:05",
				Indent:     "  ",
			}),
		)

		logger.Info("indented", slog.String("key", "val"))
		output := buf.String()
		assert.Contains(t, output, "indented")
	})

	t.Run("pretty all log levels", func(t *testing.T) {
		levels := []struct {
			level slog.Level
			name  string
		}{
			{slog.LevelDebug, "DEBUG"},
			{slog.LevelInfo, "INFO"},
			{slog.LevelWarn, "WARN"},
			{slog.LevelError, "ERROR"},
		}

		for _, l := range levels {
			var buf bytes.Buffer
			logger := New(
				WithWriter(&buf),
				WithLevel(slog.LevelDebug),
				WithPrettyConfig(&PrettyConfig{
					Enable:     true,
					Color:      true,
					TimeFormat: "15:04:05",
				}),
			)

			switch l.level {
			case slog.LevelDebug:
				logger.Debug("msg")
			case slog.LevelInfo:
				logger.Info("msg")
			case slog.LevelWarn:
				logger.Warn("msg")
			case slog.LevelError:
				logger.Error("msg")
			}

			assert.NotEmpty(t, buf.String(), "level %s should produce output", l.name)
		}
	})
}

func TestNew_WithAttributes(t *testing.T) {
	var buf bytes.Buffer
	logger := New(
		WithWriter(&buf),
		WithPrettyConfig(&PrettyConfig{Enable: false}),
	)

	logger.Info("multi attrs",
		slog.String("string", "val"),
		slog.Int("int", 42),
		slog.Bool("bool", true),
		slog.Float64("float", 3.14),
	)

	output := buf.String()
	assert.Contains(t, output, "multi attrs")
	assert.Contains(t, output, "string")
	assert.Contains(t, output, "42")
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, "json", cfg.Type)
	assert.Equal(t, slog.LevelInfo, cfg.Level)
	assert.NotNil(t, cfg.Pretty)
	assert.True(t, cfg.Pretty.Enable)
	assert.True(t, cfg.Pretty.Color)
	assert.NotEmpty(t, cfg.Pretty.TimeFormat)
}

func TestWithWriter(t *testing.T) {
	var buf bytes.Buffer
	cfg := DefaultConfig()
	WithWriter(&buf)(&cfg)
	assert.Equal(t, &buf, cfg.writer)
}

func TestWithType(t *testing.T) {
	cfg := DefaultConfig()

	WithType("text")(&cfg)
	assert.Equal(t, "text", cfg.Type)

	WithType("json")(&cfg)
	assert.Equal(t, "json", cfg.Type)
}

func TestWithLevel(t *testing.T) {
	cfg := DefaultConfig()

	WithLevel(slog.LevelDebug)(&cfg)
	assert.Equal(t, slog.LevelDebug, cfg.Level)

	WithLevel(slog.LevelError)(&cfg)
	assert.Equal(t, slog.LevelError, cfg.Level)
}

func TestWithPrettyConfig(t *testing.T) {
	cfg := DefaultConfig()

	prettyCfg := &PrettyConfig{
		Enable:     false,
		Color:      false,
		TimeFormat: "2006-01-02",
		Indent:     "\t",
	}

	WithPrettyConfig(prettyCfg)(&cfg)
	assert.Equal(t, prettyCfg, cfg.Pretty)
}

func TestWithHandlerOptions(t *testing.T) {
	cfg := DefaultConfig()

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelWarn,
	}

	WithHandlerOptions(opts)(&cfg)
	assert.Equal(t, opts, cfg.handlerOptions)
}

func TestPrettyHandler_NoAttributes(t *testing.T) {
	var buf bytes.Buffer
	logger := New(
		WithWriter(&buf),
		WithPrettyConfig(&PrettyConfig{
			Enable:     true,
			Color:      false,
			TimeFormat: "15:04:05",
		}),
	)

	logger.Info("no attrs")
	output := buf.String()
	assert.Contains(t, output, "no attrs")
	// Should have empty JSON object
	assert.True(t, strings.Contains(output, "{}"))
}
