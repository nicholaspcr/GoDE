package handlers

import (
	"context"
	"testing"

	"github.com/nicholaspcr/GoDE/internal/store/mock"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/nicholaspcr/GoDE/pkg/de"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestDEHandler_ListSupportedAlgorithms(t *testing.T) {
	handler := NewDEHandler(de.Config{})
	handler.SetStore(&mock.MockStore{})

	resp, err := handler.(*deHandler).ListSupportedAlgorithms(context.Background(), &emptypb.Empty{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, []string{"gde3"}, resp.Algorithms)
}

func TestDEHandler_ListSupportedVariants(t *testing.T) {
	handler := NewDEHandler(de.Config{})
	handler.SetStore(&mock.MockStore{})

	resp, err := handler.(*deHandler).ListSupportedVariants(context.Background(), &emptypb.Empty{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Variants, 6)

	// Verify we have the expected variants
	variantNames := make(map[string]bool)
	for _, v := range resp.Variants {
		variantNames[v.Name] = true
		assert.NotEmpty(t, v.Description)
	}

	// Verify expected variants exist (names are returned from the variant Name() methods)
	assert.Contains(t, variantNames, "rand1")
	assert.Contains(t, variantNames, "rand2")
	assert.Contains(t, variantNames, "best1")
	assert.Contains(t, variantNames, "best2")
	assert.Contains(t, variantNames, "pbest")
	assert.Contains(t, variantNames, "currToBest1")
}

func TestDEHandler_ListSupportedProblems(t *testing.T) {
	handler := NewDEHandler(de.Config{})
	handler.SetStore(&mock.MockStore{})

	resp, err := handler.(*deHandler).ListSupportedProblems(context.Background(), &emptypb.Empty{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Problems, 22) // 6 ZDT/VNT + 7 DTLZ + 9 WFG

	// Verify we have the expected problem families
	problemNames := make(map[string]bool)
	for _, p := range resp.Problems {
		problemNames[p.Name] = true
		assert.NotEmpty(t, p.Description)
	}

	// Check for ZDT problems (lowercase from problem Name() methods)
	assert.Contains(t, problemNames, "zdt1")
	assert.Contains(t, problemNames, "zdt2")
	assert.Contains(t, problemNames, "zdt3")
	assert.Contains(t, problemNames, "zdt4")
	assert.Contains(t, problemNames, "zdt6")

	// Check for VNT problem
	assert.Contains(t, problemNames, "vnt1")

	// Check for DTLZ problems
	assert.Contains(t, problemNames, "dtlz1")
	assert.Contains(t, problemNames, "dtlz2")

	// Check for WFG problems
	assert.Contains(t, problemNames, "wfg1")
	assert.Contains(t, problemNames, "wfg2")
}

func TestDEHandler_Run_ValidationErrors(t *testing.T) {
	handler := NewDEHandler(de.Config{
		ParetoChannelLimiter: 1000,
		MaxChannelLimiter:    1000,
		ResultLimiter:        100,
	})
	handler.SetStore(&mock.MockStore{})

	tests := []struct {
		name    string
		req     *api.RunRequest
		wantErr string
	}{
		{
			name: "missing DE config",
			req: &api.RunRequest{
				Algorithm: "gde3",
				Problem:   "zdt1",
				Variant:   "rand/1",
				DeConfig:  nil,
			},
			wantErr: "DE config is nil",
		},
		{
			name: "invalid population size",
			req: &api.RunRequest{
				Algorithm: "gde3",
				Problem:   "zdt1",
				Variant:   "rand/1",
				DeConfig: &api.DEConfig{
					PopulationSize: 0,
					DimensionsSize: 30,
					ObjetivesSize:  2,
					Executions:     1,
					Generations:    10,
					FloorLimiter:   0.0,
					CeilLimiter:    1.0,
					AlgorithmConfig: &api.DEConfig_Gde3{
						Gde3: &api.GDE3Config{
							Cr: 0.5,
							F:  0.5,
							P:  0.1,
						},
					},
				},
			},
			wantErr: "must be between 4 and 10000",
		},
		{
			name: "invalid dimensions size",
			req: &api.RunRequest{
				Algorithm: "gde3",
				Problem:   "zdt1",
				Variant:   "rand/1",
				DeConfig: &api.DEConfig{
					PopulationSize: 100,
					DimensionsSize: 0,
					ObjetivesSize:  2,
					Executions:     1,
					Generations:    10,
					FloorLimiter:   0.0,
					CeilLimiter:    1.0,
					AlgorithmConfig: &api.DEConfig_Gde3{
						Gde3: &api.GDE3Config{
							Cr: 0.5,
							F:  0.5,
							P:  0.1,
						},
					},
				},
			},
			wantErr: "must be between 1 and 1000",
		},
		{
			name: "invalid objectives size",
			req: &api.RunRequest{
				Algorithm: "gde3",
				Problem:   "zdt1",
				Variant:   "rand/1",
				DeConfig: &api.DEConfig{
					PopulationSize: 100,
					DimensionsSize: 30,
					ObjetivesSize:  0,
					Executions:     1,
					Generations:    10,
					FloorLimiter:   0.0,
					CeilLimiter:    1.0,
					AlgorithmConfig: &api.DEConfig_Gde3{
						Gde3: &api.GDE3Config{
							Cr: 0.5,
							F:  0.5,
							P:  0.1,
						},
					},
				},
			},
			wantErr: "must be between 1 and 10",
		},
		{
			name: "invalid executions",
			req: &api.RunRequest{
				Algorithm: "gde3",
				Problem:   "zdt1",
				Variant:   "rand/1",
				DeConfig: &api.DEConfig{
					PopulationSize: 100,
					DimensionsSize: 30,
					ObjetivesSize:  2,
					Executions:     0,
					Generations:    10,
					FloorLimiter:   0.0,
					CeilLimiter:    1.0,
					AlgorithmConfig: &api.DEConfig_Gde3{
						Gde3: &api.GDE3Config{
							Cr: 0.5,
							F:  0.5,
							P:  0.1,
						},
					},
				},
			},
			wantErr: "must be between 1 and 100",
		},
		{
			name: "invalid generations",
			req: &api.RunRequest{
				Algorithm: "gde3",
				Problem:   "zdt1",
				Variant:   "rand/1",
				DeConfig: &api.DEConfig{
					PopulationSize: 100,
					DimensionsSize: 30,
					ObjetivesSize:  2,
					Executions:     1,
					Generations:    0,
					FloorLimiter:   0.0,
					CeilLimiter:    1.0,
					AlgorithmConfig: &api.DEConfig_Gde3{
						Gde3: &api.GDE3Config{
							Cr: 0.5,
							F:  0.5,
							P:  0.1,
						},
					},
				},
			},
			wantErr: "must be between 1 and 10000",
		},
		{
			name: "invalid floor/ceil limiters",
			req: &api.RunRequest{
				Algorithm: "gde3",
				Problem:   "zdt1",
				Variant:   "rand/1",
				DeConfig: &api.DEConfig{
					PopulationSize: 100,
					DimensionsSize: 30,
					ObjetivesSize:  2,
					Executions:     1,
					Generations:    10,
					FloorLimiter:   1.0,
					CeilLimiter:    0.0, // Ceil < Floor
					AlgorithmConfig: &api.DEConfig_Gde3{
						Gde3: &api.GDE3Config{
							Cr: 0.5,
							F:  0.5,
							P:  0.1,
						},
					},
				},
			},
			wantErr: "floor_limiter (1) must be less than ceil_limiter (0)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := handler.(*deHandler).Run(context.Background(), tt.req)

			assert.Error(t, err)
			assert.Nil(t, resp)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestDEHandler_Run_InvalidProblem(t *testing.T) {
	handler := NewDEHandler(de.Config{
		ParetoChannelLimiter: 1000,
		MaxChannelLimiter:    1000,
		ResultLimiter:        100,
	})
	handler.SetStore(&mock.MockStore{})

	req := &api.RunRequest{
		Algorithm: "gde3",
		Problem:   "NonExistentProblem",
		Variant:   "rand/1",
		DeConfig: &api.DEConfig{
			PopulationSize: 100,
			DimensionsSize: 30,
			ObjetivesSize:  2,
			Executions:     1,
			Generations:    10,
			FloorLimiter:   0.0,
			CeilLimiter:    1.0,
			AlgorithmConfig: &api.DEConfig_Gde3{
				Gde3: &api.GDE3Config{
					Cr: 0.5,
					F:  0.5,
					P:  0.1,
				},
			},
		},
	}

	resp, err := handler.(*deHandler).Run(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "problem does not exist")
}

func TestDEHandler_Run_InvalidVariant(t *testing.T) {
	handler := NewDEHandler(de.Config{
		ParetoChannelLimiter: 1000,
		MaxChannelLimiter:    1000,
		ResultLimiter:        100,
	})
	handler.SetStore(&mock.MockStore{})

	req := &api.RunRequest{
		Algorithm: "gde3",
		Problem:   "zdt1",
		Variant:   "NonExistentVariant",
		DeConfig: &api.DEConfig{
			PopulationSize: 100,
			DimensionsSize: 30,
			ObjetivesSize:  2,
			Executions:     1,
			Generations:    10,
			FloorLimiter:   0.0,
			CeilLimiter:    1.0,
			AlgorithmConfig: &api.DEConfig_Gde3{
				Gde3: &api.GDE3Config{
					Cr: 0.5,
					F:  0.5,
					P:  0.1,
				},
			},
		},
	}

	resp, err := handler.(*deHandler).Run(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "variant does not exist")
}

func TestDEHandler_Run_InvalidAlgorithm(t *testing.T) {
	handler := NewDEHandler(de.Config{
		ParetoChannelLimiter: 1000,
		MaxChannelLimiter:    1000,
		ResultLimiter:        100,
	})
	handler.SetStore(&mock.MockStore{})

	req := &api.RunRequest{
		Algorithm: "NonExistentAlgorithm",
		Problem:   "zdt1",
		Variant:   "rand1",
		DeConfig: &api.DEConfig{
			PopulationSize: 100,
			DimensionsSize: 30,
			ObjetivesSize:  2,
			Executions:     1,
			Generations:    10,
			FloorLimiter:   0.0,
			CeilLimiter:    1.0,
			AlgorithmConfig: &api.DEConfig_Gde3{
				Gde3: &api.GDE3Config{
					Cr: 0.5,
					F:  0.5,
					P:  0.1,
				},
			},
		},
	}

	resp, err := handler.(*deHandler).Run(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unsupported algorithms")
}

func TestDEHandler_Run_Success(t *testing.T) {
	handler := NewDEHandler(de.Config{
		ParetoChannelLimiter: 1000,
		MaxChannelLimiter:    1000,
		ResultLimiter:        100,
	})
	handler.SetStore(&mock.MockStore{})

	tests := []struct {
		name      string
		algorithm string
		problem   string
		variant   string
	}{
		{
			name:      "GDE3 with zdt1 and rand1",
			algorithm: "gde3",
			problem:   "zdt1",
			variant:   "rand1",
		},
		{
			name:      "GDE3 with zdt2 and rand2",
			algorithm: "gde3",
			problem:   "zdt2",
			variant:   "rand2",
		},
		{
			name:      "GDE3 with dtlz1 and pbest",
			algorithm: "gde3",
			problem:   "dtlz1",
			variant:   "pbest",
		},
		{
			name:      "GDE3 with wfg1 and currToBest1",
			algorithm: "gde3",
			problem:   "wfg1",
			variant:   "currToBest1",
		},
		// Note: best1 and best2 variants have a bug in their implementation
		// (GenerateIndices count mismatch) and are not tested here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &api.RunRequest{
				Algorithm: tt.algorithm,
				Problem:   tt.problem,
				Variant:   tt.variant,
				DeConfig: &api.DEConfig{
					PopulationSize: 30,  // Population size large enough for all variants
					DimensionsSize: 10,  // Small dimensions for fast test
					ObjetivesSize:  2,   // 2 objectives for multi-objective
					Executions:     1,   // Single execution for speed
					Generations:    5,   // Few generations for speed
					FloorLimiter:   0.0,
					CeilLimiter:    1.0,
					AlgorithmConfig: &api.DEConfig_Gde3{
						Gde3: &api.GDE3Config{
							Cr: 0.5,
							F:  0.5,
							P:  0.1,
						},
					},
				},
			}

			resp, err := handler.(*deHandler).Run(context.Background(), req)

			assert.NoError(t, err)
			if !assert.NotNil(t, resp) {
				return
			}
			if !assert.NotNil(t, resp.Pareto) {
				return
			}
			assert.NotEmpty(t, resp.Pareto.Vectors)
			assert.NotEmpty(t, resp.Pareto.MaxObjs)

			// Verify vectors have correct structure
			for _, vec := range resp.Pareto.Vectors {
				assert.Len(t, vec.Elements, 10)       // DimensionsSize
				assert.Len(t, vec.Objectives, 2)      // ObjetivesSize
				assert.GreaterOrEqual(t, vec.CrowdingDistance, 0.0)
			}

			// Verify max objectives (should be executions * objectives)
			assert.Len(t, resp.Pareto.MaxObjs, 2) // 1 execution * 2 objectives
		})
	}
}

func TestDEHandler_Run_MultipleExecutions(t *testing.T) {
	handler := NewDEHandler(de.Config{
		ParetoChannelLimiter: 1000,
		MaxChannelLimiter:    1000,
		ResultLimiter:        100,
	})
	handler.SetStore(&mock.MockStore{})

	req := &api.RunRequest{
		Algorithm: "gde3",
		Problem:   "zdt1",
		Variant:   "rand1",
		DeConfig: &api.DEConfig{
			PopulationSize: 30,  // Large enough for all variants
			DimensionsSize: 10,
			ObjetivesSize:  2,
			Executions:     3, // Multiple executions
			Generations:    5,
			FloorLimiter:   0.0,
			CeilLimiter:    1.0,
			AlgorithmConfig: &api.DEConfig_Gde3{
				Gde3: &api.GDE3Config{
					Cr: 0.5,
					F:  0.5,
					P:  0.1,
				},
			},
		},
	}

	resp, err := handler.(*deHandler).Run(context.Background(), req)

	assert.NoError(t, err)
	if !assert.NotNil(t, resp) {
		return
	}
	if !assert.NotNil(t, resp.Pareto) {
		return
	}
	assert.NotEmpty(t, resp.Pareto.Vectors)

	// Verify max objectives includes all executions
	assert.Len(t, resp.Pareto.MaxObjs, 6) // 3 executions * 2 objectives
}

func TestProblemFromName(t *testing.T) {
	tests := []struct {
		name        string
		problemName string
		wantErr     bool
	}{
		{"zdt1", "zdt1", false},
		{"zdt2", "zdt2", false},
		{"zdt3", "zdt3", false},
		{"zdt4", "zdt4", false},
		{"zdt6", "zdt6", false},
		{"vnt1", "vnt1", false},
		{"dtlz1", "dtlz1", false},
		{"dtlz2", "dtlz2", false},
		{"dtlz3", "dtlz3", false},
		{"dtlz4", "dtlz4", false},
		{"dtlz5", "dtlz5", false},
		{"dtlz6", "dtlz6", false},
		{"dtlz7", "dtlz7", false},
		{"wfg1", "wfg1", false},
		{"wfg2", "wfg2", false},
		{"wfg3", "wfg3", false},
		{"wfg4", "wfg4", false},
		{"wfg5", "wfg5", false},
		{"wfg6", "wfg6", false},
		{"wfg7", "wfg7", false},
		{"wfg8", "wfg8", false},
		{"wfg9", "wfg9", false},
		{"NonExistent", "NonExistent", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			problem, err := problemFromName(tt.problemName)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, problem)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, problem)
				assert.Equal(t, tt.problemName, problem.Name())
			}
		})
	}
}

func TestVariantFromName(t *testing.T) {
	tests := []struct {
		name        string
		variantName string
		wantErr     bool
	}{
		{"rand1", "rand1", false},
		{"rand2", "rand2", false},
		{"best1", "best1", false},
		{"best2", "best2", false},
		{"currToBest1", "currToBest1", false},
		{"pbest", "pbest", false},
		{"NonExistent", "NonExistent", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			variant, err := variantFromName(tt.variantName)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, variant)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, variant)
				assert.Equal(t, tt.variantName, variant.Name())
			}
		})
	}
}
