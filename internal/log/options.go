package log

import (
	"io"
	"log/slog"
	"os"
)

// Config contains structured logging configuration options.
type Config struct {
	writer         io.Writer
	Pretty         *PrettyConfig `json:"pretty" yaml:"pretty"`
	handlerOptions *slog.HandlerOptions
	Type           string     `json:"type" yaml:"type"`
	Level          slog.Level `json:"level" yaml:"level"`
}

// PrettyConfig contain configurations regarding formatting for JSON logs.
type PrettyConfig struct {
	TimeFormat string `json:"time-format" yaml:"time-format"`
	Indent     string `json:"indent" yaml:"indent"`
	Enable     bool   `json:"enable" yaml:"enable"`
	Color      bool   `json:"color" yaml:"color"`
}

var defaultConfig = Config{
	writer: os.Stdout,
	Type:   "json",
	Level:  slog.LevelInfo,
	Pretty: &PrettyConfig{
		Enable:     true,
		Color:      true,
		TimeFormat: "[15:05:05.000]",
	},

	handlerOptions: &slog.HandlerOptions{},
}

// DefaultConfig returns the default configuration for the logger.
func DefaultConfig() Config { return defaultConfig }

// Option is a functional option for configuring the logger.
type Option func(*Config)

// WithWriter sets the output writer for the logger.
func WithWriter(w io.Writer) Option {
	return func(c *Config) { c.writer = w }
}

// WithType sets the log format type (json, text, etc).
func WithType(t string) Option {
	return func(c *Config) { c.Type = t }
}

// WithLevel sets the minimum log level.
func WithLevel(l slog.Level) Option {
	return func(c *Config) { c.Level = l }
}

// WithPrettyConfig sets pretty-printing options for logs.
func WithPrettyConfig(prettyCfg *PrettyConfig) Option {
	return func(c *Config) { c.Pretty = prettyCfg }
}

// WithHandlerOptions sets custom slog handler options.
func WithHandlerOptions(opts *slog.HandlerOptions) Option {
	return func(c *Config) { c.handlerOptions = opts }
}
