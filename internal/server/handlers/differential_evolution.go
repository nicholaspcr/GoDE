package handlers

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/nicholaspcr/GoDE/pkg/de"
	"github.com/nicholaspcr/GoDE/pkg/de/gde3"
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
	_ "github.com/nicholaspcr/GoDE/pkg/problems/many/dtlz" // Register DTLZ problems
	_ "github.com/nicholaspcr/GoDE/pkg/problems/many/wfg"  // Register WFG problems
	_ "github.com/nicholaspcr/GoDE/pkg/problems/multi"     // Register multi-objective problems
	"github.com/nicholaspcr/GoDE/pkg/validation"
	"github.com/nicholaspcr/GoDE/pkg/variants"
	_ "github.com/nicholaspcr/GoDE/pkg/variants/best"              // Register best variants
	_ "github.com/nicholaspcr/GoDE/pkg/variants/current-to-best"   // Register current-to-best variant
	_ "github.com/nicholaspcr/GoDE/pkg/variants/pbest"             // Register pbest variant
	_ "github.com/nicholaspcr/GoDE/pkg/variants/rand"              // Register rand variants
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// deHandler is responsible for the de service operations.
type deHandler struct {
	api.UnimplementedDifferentialEvolutionServiceServer
	store.Store
	cfg de.Config
}

// NewDEHandler returns a handle that implements
// DifferentialEvolutionServiceServer.
func NewDEHandler(deCfg de.Config) Handler { return &deHandler{cfg: deCfg} }

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

func (deh *deHandler) SetStore(st store.Store) {
	deh.Store = st
}

func (deh *deHandler) ListSupportedAlgorithms(
	ctx context.Context, _ *emptypb.Empty,
) (*api.ListSupportedAlgorithmsResponse, error) {
	return &api.ListSupportedAlgorithmsResponse{
		Algorithms: []string{"gde3"},
	}, nil
}

func (deh *deHandler) ListSupportedVariants(
	ctx context.Context, _ *emptypb.Empty,
) (*api.ListSupportedVariantsResponse, error) {
	metas := variants.DefaultRegistry.ListMetadata()
	apiVariants := make([]*api.Variant, len(metas))
	for i, meta := range metas {
		apiVariants[i] = &api.Variant{
			Name:        meta.Name,
			Description: meta.Description,
		}
	}
	return &api.ListSupportedVariantsResponse{Variants: apiVariants}, nil
}

func (deh *deHandler) ListSupportedProblems(
	ctx context.Context, _ *emptypb.Empty,
) (*api.ListSupportedProblemsResponse, error) {
	metas := problems.DefaultRegistry.ListMetadata()
	apiProblems := make([]*api.Problem, len(metas))
	for i, meta := range metas {
		apiProblems[i] = &api.Problem{
			Name:        meta.Name,
			Description: meta.Description,
		}
	}
	return &api.ListSupportedProblemsResponse{Problems: apiProblems}, nil
}

func (deh *deHandler) Run(
	ctx context.Context, req *api.RunRequest,
) (*api.RunResponse, error) {
	tracer := otel.Tracer("handlers.de")
	ctx, span := tracer.Start(ctx, "deHandler.Run")
	defer span.End()

	span.SetAttributes(
		attribute.String("algorithm", req.Algorithm),
		attribute.String("problem", req.Problem),
		attribute.String("variant", req.Variant),
	)

	// Validate DE configuration
	if err := validation.ValidateDEConfig(req.DeConfig); err != nil {
		span.RecordError(err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	span.SetAttributes(
		attribute.Int64("executions", int64(req.DeConfig.Executions)),
		attribute.Int64("generations", int64(req.DeConfig.Generations)),
		attribute.Int64("population_size", int64(req.DeConfig.PopulationSize)),
	)

	var algo de.Algorithm

	populationParams := models.PopulationParams{
		PopulationSize: int(req.DeConfig.PopulationSize),
		DimensionSize:  int(req.DeConfig.DimensionsSize),
		ObjectivesSize: int(req.DeConfig.ObjetivesSize),
		FloorRange:     make([]float64, req.DeConfig.DimensionsSize),
		CeilRange:      make([]float64, req.DeConfig.DimensionsSize),
	}

	// Setup the limiters for the population.
	for i := range populationParams.CeilRange {
		populationParams.CeilRange[i] = float64(req.DeConfig.CeilLimiter)
		populationParams.FloorRange[i] = float64(req.DeConfig.FloorLimiter)
	}

	problem, err := problemFromName(req.Problem, int(req.DeConfig.DimensionsSize), int(req.DeConfig.ObjetivesSize))
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	variant, err := variantFromName(req.Variant)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	initialPopulation, err := generatePopulation(populationParams, rand.New(rand.NewSource(time.Now().UnixNano())))
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	switch req.Algorithm {
	case "gde3":
		algo = gde3.New(
			gde3.WithConstants(gde3.Constants{
				DE: de.Constants{
					Executions:    int(req.DeConfig.Executions),
					Generations:   int(req.DeConfig.Generations),
					Dimensions:    int(req.DeConfig.DimensionsSize),
					ObjFuncAmount: int(req.DeConfig.ObjetivesSize),
				},
				CR: float64(req.DeConfig.GetGde3().Cr),
				F:  float64(req.DeConfig.GetGde3().F),
				P:  float64(req.DeConfig.GetGde3().P),
			}),
			gde3.WithInitialPopulation(initialPopulation),
			gde3.WithPopulationParams(populationParams),
			gde3.WithProblem(problem),
			gde3.WithVariant(variant),
		)

	default:
		err := errors.New("unsupported algorithms")
		span.RecordError(err)
		return nil, err
	}

	DE, err := de.New(
		deh.cfg,
		de.WithAlgorithm(algo),
		de.WithExecutions(int(req.DeConfig.Executions)),
		de.WithGenerations(int(req.DeConfig.Generations)),
		de.WithDimensions(int(req.DeConfig.DimensionsSize)),
		de.WithObjFuncAmount(int(req.DeConfig.ObjetivesSize)),
	)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	finalPareto, maxObjs, err := DE.Execute(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	span.SetAttributes(attribute.Int("pareto_result_size", len(finalPareto)))

	// Convert vectors to proto format
	vectorsProto := make([]*api.Vector, len(finalPareto))
	for i, vec := range finalPareto {
		vectorsProto[i] = &api.Vector{
			Elements:         vec.Elements,
			Objectives:       vec.Objectives,
			CrowdingDistance: vec.CrowdingDistance,
		}
	}

	// Flatten max objectives from all executions
	flatMaxObjs := make([]float64, 0, len(maxObjs)*populationParams.ObjectivesSize)
	for _, objs := range maxObjs {
		flatMaxObjs = append(flatMaxObjs, objs...)
	}

	return &api.RunResponse{
		Pareto: &api.Pareto{
			Vectors: vectorsProto,
			MaxObjs: flatMaxObjs,
		},
	}, nil
}

// problemFromName returns the problems.Interface implementation of the problem
// referenced by name using the registry pattern.
func problemFromName(name string, dim, objs int) (problems.Interface, error) {
	return problems.DefaultRegistry.Create(name, dim, objs)
}

// variantFromName returns the variants.Interface implementation of the variant
// referenced by name using the registry pattern.
func variantFromName(name string) (variants.Interface, error) {
	return variants.DefaultRegistry.Create(name)
}
