package log

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"

	"github.com/fatih/color"
)

// New creates a new logger with the given options. If no option is given then
// it will use the default configuration.
func New(opts ...Option) *slog.Logger {
	cfg := new(Config)
	*cfg = *defaultConfig

	for _, opt := range opts {
		opt(cfg)
	}
	// Manually pass the level to the handler options
	cfg.HandlerOptions.Level = cfg.Level

	var h slog.Handler
	switch cfg.Type {
	case "text":
		h = slog.NewTextHandler(cfg.Writer, cfg.HandlerOptions)
	default:
		h = slog.NewJSONHandler(cfg.Writer, cfg.HandlerOptions)
	}

	if cfg.Pretty.Enable {
		h = &prettyHandler{
			Handler: h,
			Logger:  log.New(cfg.Writer, "", 0),
			cfg:     cfg.Pretty,
		}
	}

	return slog.New(h)
}

type prettyHandler struct {
	slog.Handler
	*log.Logger
	cfg *PrettyConfig
}

func (h *prettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String()
	if h.cfg.Color {
		switch r.Level {
		case slog.LevelDebug:
			level = color.MagentaString(level) + ":"
		case slog.LevelInfo:
			level = color.BlueString(level) + ":"
		case slog.LevelWarn:
			level = color.YellowString(level) + ":"
		case slog.LevelError:
			level = color.RedString(level) + ":"
		}
	}

	fields := make(map[string]any, r.NumAttrs())

	// Local fields of the logger
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	var b []byte
	var err error
	if h.cfg.Indent != "" {
		b, err = json.MarshalIndent(fields, "", h.cfg.Indent)
		if err != nil {
			return err
		}
	} else {
		b, err = json.Marshal(fields)
		if err != nil {
			return err
		}
	}

	timeStr := r.Time.Format(h.cfg.TimeFormat)
	msg := color.CyanString(r.Message)

	h.Println(timeStr, level, msg, color.WhiteString(string(b)))
	return nil
}
