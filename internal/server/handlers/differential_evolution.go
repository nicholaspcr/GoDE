package handlers

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nicholaspcr/GoDE/internal/executor"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	_ "github.com/nicholaspcr/GoDE/pkg/de/gde3"                    // Register GDE3 algorithm
	_ "github.com/nicholaspcr/GoDE/pkg/problems/many/dtlz"         // Register DTLZ problems
	_ "github.com/nicholaspcr/GoDE/pkg/problems/many/wfg"          // Register WFG problems
	_ "github.com/nicholaspcr/GoDE/pkg/problems/multi"             // Register multi-objective problems
	_ "github.com/nicholaspcr/GoDE/pkg/variants/best"              // Register best variants
	_ "github.com/nicholaspcr/GoDE/pkg/variants/current-to-best"   // Register current-to-best variant
	_ "github.com/nicholaspcr/GoDE/pkg/variants/pbest"             // Register pbest variant
	_ "github.com/nicholaspcr/GoDE/pkg/variants/rand"              // Register rand variants
	"google.golang.org/grpc"
)

// deHandler is responsible for the de service operations.
type deHandler struct {
	api.UnimplementedDifferentialEvolutionServiceServer
	Store    deStore
	executor *executor.Executor
}

// NewDEHandler returns a handler that implements
// DifferentialEvolutionServiceServer.
func NewDEHandler(st deStore, exec *executor.Executor) Handler {
	return &deHandler{Store: st, executor: exec}
}

// RegisterService adds DifferentialEvolutionService to the RPC server.
func (deh *deHandler) RegisterService(srv *grpc.Server) {
	api.RegisterDifferentialEvolutionServiceServer(srv, deh)
}

// RegisterHTTPHandler adds DifferentialEvolutionService to the grpc-gateway.
func (deh *deHandler) RegisterHTTPHandler(
	ctx context.Context,
	mux *runtime.ServeMux,
	lisAddr string,
	dialOpts []grpc.DialOption,
) error {
	return api.RegisterDifferentialEvolutionServiceHandlerFromEndpoint(
		ctx, mux, lisAddr, dialOpts,
	)
}
