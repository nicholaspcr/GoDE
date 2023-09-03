package commands

import "github.com/nicholaspcr/GoDE/internal/log"

// logOptionsFromConfig returns a slice of log Options after processing the log values in the CLI configuration.
func logOptionsFromConfig(cfg *log.Config) []log.Option {
	opts := make([]log.Option, 0, 10)
	if cfg == nil {
		return nil
	}
	if cfg.Writer != nil {
		opts = append(opts, log.WithWriter(cfg.Writer))
	}

	if cfg.Type != "" {
		opts = append(opts, log.WithType(cfg.Type))
	}

	if cfg.Level != 0 {
		opts = append(opts, log.WithLevel(cfg.Level))
	}

	if cfg.HandlerOptions != nil {
		opts = append(opts, log.WithHandlerOptions(cfg.HandlerOptions))
	}

	if cfg.Pretty != nil {
		opts = append(opts, log.WithPrettyConfig(cfg.Pretty))
	}
	return opts
}
