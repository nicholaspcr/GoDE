package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nicholaspcr/GoDE/internal/server/handlers"
	"github.com/nicholaspcr/GoDE/internal/server/middleware"
	"github.com/nicholaspcr/GoDE/internal/server/session"
	"github.com/nicholaspcr/GoDE/internal/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Server is the interface that wraps the server methods.
type Server interface {
	// Start starts the server.
	Start(context.Context) error
}

// New returns a new server instance
func New(ctx context.Context, cfg Config, opts ...serverOpts) (Server, error) {
	sessionStore := session.NewInMemoryStore()

	srv := &server{
		cfg: cfg,
		handlers: []handlers.Handler{
			handlers.NewAuthHandler(sessionStore),
			handlers.NewUserHandler(),
			handlers.NewDEHandler(cfg.DE),
		},
		sessionStore: sessionStore,
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
	st           store.Store
	cfg          Config
	handlers     []handlers.Handler
	sessionStore session.Store
}

// Start starts the server.
func (s *server) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	slog.Info("Creating server")

	logger := slog.Default()

	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(logger)),
			middleware.UnaryAuthMiddleware(s.sessionStore),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(InterceptorLogger(logger)),
		),
	)

	slog.Info("Registering RPC services")
	for _, handler := range s.handlers {
		handler.RegisterService(grpcSrv)
	}

	slog.Info("Creating listener")
	lis, err := net.Listen("tcp", s.cfg.LisAddr)
	if err != nil {
		return err
	}
	lisAddr := lis.Addr().String()

	slog.Info("RPC server listening on: ", slog.String("addr", lisAddr))
	go func() {
		if err := grpcSrv.Serve(lis); err != nil {
			slog.Error(
				"Unexpected panic on the server",
				slog.String("error", err.Error()),
			)
			cancel()
		}
	}()

	// NOTE: Make sure the gRPC server is running properly and accessible.
	mux := runtime.NewServeMux()
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	slog.Info("Registering grpc-gateway handlers")
	for _, handler := range s.handlers {
		handler.RegisterHTTPHandler(ctx, mux, lisAddr, dialOpts)
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	slog.Info("HTTP server listening on: ", slog.String("port", s.cfg.HTTPPort))
	go func() {
		if err := http.ListenAndServe(s.cfg.HTTPPort, mux); err != nil {
			slog.Error(
				"Unexpected panic on the web server",
				slog.String("error", err.Error()),
			)
			cancel()
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	slog.Info("Shutting down server")
	return nil
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
