package server

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nicholaspcr/GoDE/cmd/deserver/internal/handlers"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Server is the interface that wraps the server methods.
type Server interface {
	// Start starts the server.
	Start(context.Context) error
}

// New returns a new server instance
func New(_ context.Context, opts ...serverOpts) (Server, error) {
	srv := &server{
		cfg: defaultConfig,
		handlers: []handlers.Handler{
			handlers.NewUserHandler(),
		},
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
	st       store.Store
	cfg      Config
	handlers []handlers.Handler
}

// Start starts the server.
func (s *server) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	slog.Info("Creating server")

	var srvOpts []grpc.ServerOption
	grpcSrv := grpc.NewServer(srvOpts...)

	slog.Info("Registering services")
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
	slog.Info("Registering grpc-gateway handlers")
	mux := runtime.NewServeMux()
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	err = api.RegisterUserServiceHandlerFromEndpoint(
		ctx, mux, lisAddr, dialOpts,
	)
	if err != nil {
		return err
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
