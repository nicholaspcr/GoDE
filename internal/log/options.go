package log

import (
	"io"
	"log/slog"
	"os"
)

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

type Option func(*Config)

func WithWriter(w io.Writer) Option {
	return func(c *Config) { c.writer = w }
}

func WithType(t string) Option {
	return func(c *Config) { c.Type = t }
}

func WithLevel(l slog.Level) Option {
	return func(c *Config) { c.Level = l }
}

func WithPrettyConfig(prettyCfg *PrettyConfig) Option {
	return func(c *Config) { c.Pretty = prettyCfg }
}

func WithHandlerOptions(opts *slog.HandlerOptions) Option {
	return func(c *Config) { c.handlerOptions = opts }
}
