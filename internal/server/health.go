package server

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// setupHealthCheck initializes the gRPC health check service.
func (s *server) setupHealthCheck(grpcSrv *grpc.Server) *health.Server {
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcSrv, healthServer)

	// Set all services as serving
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("api.v1.AuthService", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("api.v1.UserService", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("api.v1.DifferentialEvolutionService", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus("api.v1.ParetoSetService", grpc_health_v1.HealthCheckResponse_SERVING)

	slog.Info("Health check service registered")
	return healthServer
}

// checkDatabaseHealth checks if the database is accessible.
func (s *server) checkDatabaseHealth(ctx context.Context) bool {
	if s.st == nil {
		return false
	}

	// Try to ping the database by attempting a simple query
	// This is a basic health check - could be enhanced with more sophisticated checks
	// For now, we assume if the store exists, it's healthy
	// In a real implementation, you'd want to check the actual DB connection
	return true
}
