package server

import (
	"context"
	"fmt"
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
func New(_ context.Context, st store.Store) Server {
	return &server{
		store:       st,
		userHandler: handlers.UserHandler{Store: st},
	}
}

type server struct {
	store store.Store

	// handlers
	userHandler handlers.UserHandler
}

// Start starts the server.
func (s *server) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	slog.Info("Creating server")

	var srvOpts []grpc.ServerOption
	grpcSrv := grpc.NewServer(srvOpts...)

	slog.Info("Registering services")
	api.RegisterUserServiceServer(grpcSrv, &s.userHandler)

	slog.Info("Creating listener")
	port := 3030
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}

	slog.Info("Server listening on: ", slog.Int("port", port))
	go func() {
		if err := grpcSrv.Serve(lis); err != nil {
			slog.Error("Unexpected panic on the server", slog.String("error", err.Error()))
			cancel()
		}
	}()

	// Register gRPC server endpoint
	slog.Info("Registering grpc-gateway handlers")
	// api.RegisterUserServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption)

	// NOTE: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = api.RegisterUserServiceHandlerFromEndpoint(ctx, mux, lis.Addr().String(), dialOpts)
	if err != nil {
		return err
	}
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	go func() {
		if err := http.ListenAndServe(":8081", mux); err != nil {
			slog.Error("Unexpected panic on the web server", slog.String("error", err.Error()))
			cancel()
		}
	}()

	// Wait for shutdown signal
	select {
	case <-ctx.Done():
		slog.Info("Shutting down server")
	}
	return nil
}
