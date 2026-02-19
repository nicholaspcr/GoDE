// Package server contains all logic related to the DE API.
package server

import (
	"context"
	"fmt"
	"log/slog"
	_ "net/http/pprof" // Register pprof HTTP handlers

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/nicholaspcr/GoDE/internal/executor"
	"github.com/nicholaspcr/GoDE/internal/server/auth"
	"github.com/nicholaspcr/GoDE/internal/server/handlers"
	"github.com/nicholaspcr/GoDE/internal/slo"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/internal/telemetry"
	"github.com/nicholaspcr/GoDE/pkg/problems"
	"github.com/nicholaspcr/GoDE/pkg/variants"
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
		cfg:        cfg,
		jwtService: jwtService,
	}

	for _, opt := range opts {
		opt(srv)
	}

	// Create executor with the store
	if srv.st == nil {
		return nil, fmt.Errorf("store must be provided via WithStore option")
	}

	srv.executor = executor.New(executor.Config{
		Store:                srv.st,
		MaxWorkers:           cfg.Executor.MaxWorkers,
		MaxVectorsInProgress: cfg.Executor.MaxVectorsInProgress,
		ExecutionTTL:         cfg.Executor.ExecutionTTL,
		ResultTTL:            cfg.Executor.ResultTTL,
		ProgressTTL:          cfg.Executor.ProgressTTL,
		Metrics:              srv.metrics,
	})

	// Register all problems and variants
	problemMetas := problems.DefaultRegistry.ListMetadata()
	for _, meta := range problemMetas {
		// Create problem instance with default dimensions (will be overridden per execution)
		prob, err := problems.DefaultRegistry.Create(meta.Name, 10, 2)
		if err == nil {
			srv.executor.RegisterProblem(meta.Name, prob)
		}
	}

	variantMetas := variants.DefaultRegistry.ListMetadata()
	for _, meta := range variantMetas {
		variant, err := variants.DefaultRegistry.Create(meta.Name)
		if err == nil {
			srv.executor.RegisterVariant(meta.Name, variant)
		}
	}

	// Create handlers with dependencies
	srv.handlers = []handlers.Handler{
		handlers.NewAuthHandler(srv.st, jwtService, cfg.JWTExpiry),
		handlers.NewUserHandler(srv.st),
		handlers.NewParetoHandler(srv.st),
		handlers.NewDEHandler(srv.st, srv.executor),
	}

	return srv, nil
}

type server struct {
	st         store.Store
	jwtService auth.JWTService
	handlers   []handlers.Handler
	cfg        Config
	metrics    *telemetry.Metrics
	sloTracker *slo.Tracker
	executor   *executor.Executor
}

// Start starts the server using a lifecycle-based approach.
func (s *server) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	lifecycle := newLifecycle(s.cfg, s, s.executor)

	if err := lifecycle.setup(ctx); err != nil {
		cancel()
		return err
	}

	if err := lifecycle.run(ctx); err != nil {
		cancel()
		return err
	}

	// Cancel context to signal cleanup goroutine before waiting for it in shutdown
	cancel()
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
