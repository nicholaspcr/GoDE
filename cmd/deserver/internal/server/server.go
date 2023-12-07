package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api"
	"google.golang.org/grpc"
)

// Server is the interface that wraps the server methods.
type Server interface {
	// Start starts the server.
	Start(context.Context) error

	// API routes defined in the API.
	api.UserServicesServer
}

// New returns a new server instance
func New(_ context.Context, st store.Store) Server {
	return &server{
		store:      st,
		userServer: userServer{Store: st},
	}
}

type server struct {
	store store.Store
	userServer
}

// Start starts the server.
func (s *server) Start(ctx context.Context) error {
	slog.Info("Creating server")
	var opts []grpc.ServerOption
	grpcSrv := grpc.NewServer(opts...)

	slog.Info("Registering services")
	api.RegisterUserServicesServer(grpcSrv, &s.userServer)

	slog.Info("Creating listener")
	port := 3030
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}

	slog.Info("Server listening on: ", slog.Int("port", port))
	return grpcSrv.Serve(lis)
}
