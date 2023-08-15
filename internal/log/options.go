package log

import (
	"io"
	"log/slog"
	"os"
)

type Config struct {
	// Writer is the location where logs are written to. Defaults to os.Stdout.
	Writer io.Writer
	// Type can be either "json" or "text"
	Type string `json:"type" yaml:"type"`
	// Level is the minimum log level to output. Defaults to slog.LevelInfo.
	Level slog.Level `json:"level" yaml:"level"`
	// HandlerOptions are the options to pass to the handler.
	HandlerOptions *slog.HandlerOptions
	// Pretty contain configurations regarding formatting for JSON logs.
	Pretty *PrettyConfig `json:"pretty" yaml:"pretty"`
}

// PrettyConfig contain configurations regarding formatting for JSON logs.
type PrettyConfig struct {
	Enable bool `json:"enable"`
	Color  bool `json:"color"`
	// TimeFormat is the format for timestamps. Defaults to time.RFC3339.
	TimeFormat string `json:"time-format"`
	Indent     string `json:"indent"`
}

var defaultConfig = &Config{
	Writer:         os.Stdout,
	Type:           "json",
	Level:          slog.LevelInfo,
	HandlerOptions: &slog.HandlerOptions{},
	Pretty: &PrettyConfig{
		Enable:     true,
		Color:      true,
		TimeFormat: "[15:05:05.000]",
	},
}

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
	return func(c *Config) { c.HandlerOptions = opts }
}
