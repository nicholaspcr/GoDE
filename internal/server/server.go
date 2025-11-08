// Package server contains all logic related to the DE API.
package server

import (
	"context"
	"log/slog"
	_ "net/http/pprof" // Import pprof for profiling endpoints

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/nicholaspcr/GoDE/internal/server/auth"
	"github.com/nicholaspcr/GoDE/internal/server/handlers"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/internal/telemetry"
)

// Server is the interface that wraps the server methods.
type Server interface {
	// Start starts the server.
	Start(context.Context) error
}

// New returns a new server instance
func New(ctx context.Context, cfg Config, opts ...serverOpts) (Server, error) {
	jwtService := auth.NewJWTService(cfg.JWTSecret, cfg.JWTExpiry)

	srv := &server{
		cfg: cfg,
		handlers: []handlers.Handler{
			handlers.NewAuthHandler(jwtService),
			handlers.NewUserHandler(),
			handlers.NewParetoHandler(),
			handlers.NewDEHandler(cfg.DE),
		},
		jwtService: jwtService,
	}

	for _, opt := range opts {
		opt(srv)
	}

	// Setup handlers' store
	for _, handler := range srv.handlers {
		handler.SetStore(srv.st)
	}

	return srv, nil
}

type server struct {
	st         store.Store
	jwtService auth.JWTService
	handlers   []handlers.Handler
	cfg        Config
	metrics    *telemetry.Metrics
}

// Start starts the server using a lifecycle-based approach.
func (s *server) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	lifecycle := newLifecycle(s.cfg, s)

	if err := lifecycle.setup(ctx); err != nil {
		return err
	}

	if err := lifecycle.run(ctx); err != nil {
		return err
	}

	return lifecycle.shutdown(ctx)
}

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(
		func(
			ctx context.Context, lvl logging.Level, msg string, fields ...any,
		) {
			l.Log(ctx, slog.Level(lvl), msg, fields...)
		})
}
