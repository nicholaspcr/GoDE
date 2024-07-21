package log

import (
	"io"
	"log/slog"
	"os"
)

type Config struct {
	// Writer is the location where logs are written to. Defaults to os.Stdout.
	Writer io.Writer `json:",omitempty" yaml:",omitempty"`
	// Type can be either "json" or "text"
	Type string `json:"type" yaml:"type"`
	// Level is the minimum log level to output. Defaults to slog.LevelInfo.
	Level slog.Level `json:"level" yaml:"level"`
	// Pretty contain configurations regarding formatting for JSON logs.
	Pretty *PrettyConfig `json:"pretty" yaml:"pretty"`

	// handlerOptions are the options to pass to the handler.
	handlerOptions *slog.HandlerOptions `json:",omitempty" yaml:",omitempty"`
}

// PrettyConfig contain configurations regarding formatting for JSON logs.
type PrettyConfig struct {
	Enable bool `json:"enable" yaml:"enable"`
	Color  bool `json:"color" yaml:"color"`
	// TimeFormat is the format for timestamps. Defaults to time.RFC3339.
	TimeFormat string `json:"time-format" yaml:"time-format"`
	Indent     string `json:"indent" yaml:"indent"`
}

var defaultConfig = &Config{
	Writer:         os.Stdout,
	Type:           "json",
	Level:          slog.LevelInfo,
	handlerOptions: &slog.HandlerOptions{},
	Pretty: &PrettyConfig{
		Enable:     true,
		Color:      true,
		TimeFormat: "[15:05:05.000]",
	},
}

// DefaultConfig returns the default configuration for the logger.
func DefaultConfig() *Config { return defaultConfig }

type Option func(*Config)

func WithWriter(w io.Writer) Option {
	return func(c *Config) { c.Writer = w }
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
