package handlers

import (
	"context"
	"errors"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/nicholaspcr/GoDE/pkg/de"
	"github.com/nicholaspcr/GoDE/pkg/de/gde3"
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/problems"
	"github.com/nicholaspcr/GoDE/pkg/problems/many/dtlz"
	"github.com/nicholaspcr/GoDE/pkg/problems/many/wfg"
	"github.com/nicholaspcr/GoDE/pkg/problems/multi"
	"github.com/nicholaspcr/GoDE/pkg/variants"
	"github.com/nicholaspcr/GoDE/pkg/variants/best"
	currenttobest "github.com/nicholaspcr/GoDE/pkg/variants/current-to-best"
	"github.com/nicholaspcr/GoDE/pkg/variants/pbest"
	"github.com/nicholaspcr/GoDE/pkg/variants/rand"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// deHandler is responsible for the de service operations.
type deHandler struct {
	store.Store
	api.UnimplementedDifferentialEvolutionServiceServer
}

// NewDEHandler returns a handle that implements
// DifferentialEvolutionServiceServer.
func NewDEHandler() Handler { return &deHandler{} }

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
	return &api.ListSupportedVariantsResponse{
		Variants: []string{
			rand.Rand1().Name(),
			rand.Rand2().Name(),
			best.Best1().Name(),
			best.Best2().Name(),
			pbest.Pbest().Name(),
			currenttobest.CurrToBest1().Name(),
		},
	}, nil
}

func (deh *deHandler) ListSupportedProblems(
	ctx context.Context, _ *emptypb.Empty,
) (*api.ListSupportedProblemsResponse, error) {
	return &api.ListSupportedProblemsResponse{
		Problems: []string{
			multi.Zdt1().Name(),
			multi.Zdt2().Name(),
			multi.Zdt3().Name(),
			multi.Zdt4().Name(),
			multi.Zdt6().Name(),
			multi.Vnt1().Name(),

			dtlz.Dtlz1().Name(),
			dtlz.Dtlz2().Name(),
			dtlz.Dtlz3().Name(),
			dtlz.Dtlz4().Name(),
			dtlz.Dtlz5().Name(),
			dtlz.Dtlz6().Name(),
			dtlz.Dtlz7().Name(),

			wfg.Wfg1().Name(),
			wfg.Wfg2().Name(),
			wfg.Wfg3().Name(),
			wfg.Wfg4().Name(),
			wfg.Wfg5().Name(),
			wfg.Wfg6().Name(),
			wfg.Wfg7().Name(),
			wfg.Wfg8().Name(),
			wfg.Wfg9().Name(),
		},
	}, nil
}

func (deh *deHandler) Run(
	ctx context.Context, req *api.RunRequest,
) (*api.RunResponse, error) {
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

	problem, err := problemFromName(req.Problem)
	if err != nil {
		return nil, err
	}

	variant, err := variantFromName(req.Variant)
	if err != nil {
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
			gde3.WithInitialPopulation(generatePopulation(populationParams)),
			gde3.WithPopulationParams(populationParams),
			gde3.WithProblem(problem),
			gde3.WithVariant(variant),
		)

	default:
		return nil, errors.New("unsupported algorithms")
	}

	DE := de.New(
		de.WithAlgorithm(algo),
		de.WithExecutions(int(req.DeConfig.Executions)),
		de.WithGenerations(int(req.DeConfig.Generations)),
		de.WithDimensions(int(req.DeConfig.DimensionsSize)),
		de.WithObjFuncAmount(int(req.DeConfig.ObjetivesSize)),
	)
	if err := DE.Execute(ctx); err != nil {
		return nil, err
	}

	return nil, nil
}

// problemFromName returns the problems.Interface implementation of the problem
// referenced by name.
func problemFromName(p string) (problems.Interface, error) {
	switch p {
	case "zdt1":
		return multi.Zdt1(), nil
	}
	return nil, errors.New("does not exist")
}

// variantFromName returns the variants.Interface implementation of the variant
// referenced by name.
func variantFromName(p string) (variants.Interface, error) {
	switch p {
	case "rand1":
		return rand.Rand1(), nil
	}
	return nil, errors.New("does not exist")
}
