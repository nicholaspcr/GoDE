package utils

import (
	"github.com/nicholaspcr/GoDE/internal/log"
)

// LogOptionsFromConfig returns a slice of log Options after processing the log
// values in the CLI configuration.
func LogOptionsFromConfig(cfg log.Config) []log.Option {
	opts := make([]log.Option, 0, 10)
	if cfg.Writer != nil {
		opts = append(opts, log.WithWriter(cfg.Writer))
	}

	if cfg.Type != "" {
		opts = append(opts, log.WithType(cfg.Type))
	}

	if cfg.Level != 0 {
		opts = append(opts, log.WithLevel(cfg.Level))
	}

	if cfg.Pretty != nil {
		opts = append(opts, log.WithPrettyConfig(cfg.Pretty))
	}

	return opts
}
