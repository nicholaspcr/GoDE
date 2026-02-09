package decmd

import (
	"testing"

	"github.com/nicholaspcr/GoDE/cmd/decli/internal/config"
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/state"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockStateOps implements state.Operations for testing.
type mockStateOps struct {
	token string
	err   error
}

func (m *mockStateOps) GetAuthToken() (string, error) { return m.token, m.err }
func (m *mockStateOps) InvalidateAuthToken() error     { return m.err }
func (m *mockStateOps) SaveAuthToken(t string) error {
	m.token = t
	return m.err
}

func TestListAlgorithmsCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, listAlgorithmsCmd)
		assert.Equal(t, "list-algorithms", listAlgorithmsCmd.Use)
		assert.NotEmpty(t, listAlgorithmsCmd.Short)
		assert.Contains(t, listAlgorithmsCmd.Short, "algorithm")
	})

	t.Run("has RunE function", func(t *testing.T) {
		assert.NotNil(t, listAlgorithmsCmd.RunE)
	})
}

func TestListProblemsCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, listProblemsCmd)
		assert.Equal(t, "list-problems", listProblemsCmd.Use)
		assert.NotEmpty(t, listProblemsCmd.Short)
		assert.Contains(t, listProblemsCmd.Short, "problem")
	})

	t.Run("has RunE function", func(t *testing.T) {
		assert.NotNil(t, listProblemsCmd.RunE)
	})
}

func TestListVariantsCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, listVariantsCmd)
		assert.Equal(t, "list-variants", listVariantsCmd.Use)
		assert.NotEmpty(t, listVariantsCmd.Short)
		assert.Contains(t, listVariantsCmd.Short, "variant")
	})

	t.Run("has RunE function", func(t *testing.T) {
		assert.NotNil(t, listVariantsCmd.RunE)
	})
}

func TestDECommand(t *testing.T) {
	t.Run("de command exists", func(t *testing.T) {
		assert.NotNil(t, deCmd)
		assert.Equal(t, "de", deCmd.Use)
		assert.NotEmpty(t, deCmd.Short)
	})

	t.Run("has RunE that returns help", func(t *testing.T) {
		assert.NotNil(t, deCmd.RunE)
	})

	t.Run("has subcommands", func(t *testing.T) {
		commands := deCmd.Commands()
		assert.NotEmpty(t, commands)

		commandNames := make(map[string]bool)
		for _, cmd := range commands {
			commandNames[cmd.Use] = true
		}

		assert.True(t, commandNames["list-algorithms"], "list-algorithms should be registered")
		assert.True(t, commandNames["list-problems"], "list-problems should be registered")
		assert.True(t, commandNames["list-variants"], "list-variants should be registered")
		assert.True(t, commandNames["run"], "run should be registered")
		assert.True(t, commandNames["run-async"], "run-async should be registered")
		assert.True(t, commandNames["status"], "status should be registered")
		assert.True(t, commandNames["cancel"], "cancel should be registered")
		assert.True(t, commandNames["delete"], "delete should be registered")
		assert.True(t, commandNames["results"], "results should be registered")
		assert.True(t, commandNames["stream"], "stream should be registered")
		assert.True(t, commandNames["list"], "list should be registered")
	})
}

func TestRunCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, runCmd)
		assert.Equal(t, "run", runCmd.Use)
		assert.NotEmpty(t, runCmd.Short)
		assert.NotEmpty(t, runCmd.Long)
	})

	t.Run("has RunE function", func(t *testing.T) {
		assert.NotNil(t, runCmd.RunE)
	})

	t.Run("has required flags", func(t *testing.T) {
		requiredFlags := []string{"algorithm", "variant", "problem"}
		for _, name := range requiredFlags {
			flag := runCmd.Flags().Lookup(name)
			require.NotNil(t, flag, "flag %s should exist", name)
			assert.Equal(t, "", flag.DefValue, "flag %s should have empty default", name)
		}
	})

	t.Run("has DE config flags", func(t *testing.T) {
		deFlags := map[string]string{
			"executions":      "1",
			"generations":     "100",
			"population-size": "100",
			"dimensions-size": "30",
			"objectives-size": "2",
			"floor-limiter":   "0",
			"ceil-limiter":    "1",
		}
		for name, defValue := range deFlags {
			flag := runCmd.Flags().Lookup(name)
			require.NotNil(t, flag, "flag %s should exist", name)
			assert.Equal(t, defValue, flag.DefValue, "flag %s default", name)
		}
	})

	t.Run("has GDE3 config flags", func(t *testing.T) {
		gdeFlags := []string{"cr", "f", "p"}
		for _, name := range gdeFlags {
			flag := runCmd.Flags().Lookup(name)
			require.NotNil(t, flag, "flag %s should exist", name)
			assert.Equal(t, "0.5", flag.DefValue, "flag %s default", name)
		}
	})
}

func TestRunAsyncCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, runAsyncCmd)
		assert.Equal(t, "run-async", runAsyncCmd.Use)
		assert.NotEmpty(t, runAsyncCmd.Short)
		assert.NotEmpty(t, runAsyncCmd.Long)
	})

	t.Run("has RunE function", func(t *testing.T) {
		assert.NotNil(t, runAsyncCmd.RunE)
	})

	t.Run("has required flags", func(t *testing.T) {
		requiredFlags := []string{"algorithm", "variant", "problem"}
		for _, name := range requiredFlags {
			flag := runAsyncCmd.Flags().Lookup(name)
			require.NotNil(t, flag, "flag %s should exist", name)
			assert.Equal(t, "", flag.DefValue, "flag %s should have empty default", name)
		}
	})

	t.Run("has DE config flags", func(t *testing.T) {
		deFlags := map[string]string{
			"executions":      "1",
			"generations":     "100",
			"population-size": "100",
			"dimensions-size": "30",
			"objectives-size": "2",
		}
		for name, defValue := range deFlags {
			flag := runAsyncCmd.Flags().Lookup(name)
			require.NotNil(t, flag, "flag %s should exist", name)
			assert.Equal(t, defValue, flag.DefValue, "flag %s default", name)
		}
	})
}

func TestStatusCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, statusCmd)
		assert.Equal(t, "status", statusCmd.Use)
		assert.NotEmpty(t, statusCmd.Short)
		assert.NotEmpty(t, statusCmd.Long)
	})

	t.Run("has RunE function", func(t *testing.T) {
		assert.NotNil(t, statusCmd.RunE)
	})

	t.Run("has execution-id flag", func(t *testing.T) {
		flag := statusCmd.Flags().Lookup("execution-id")
		require.NotNil(t, flag)
		assert.Equal(t, "", flag.DefValue)
	})

	t.Run("requires execution-id", func(t *testing.T) {
		statusExecutionID = ""
		err := statusCmd.RunE(statusCmd, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "--execution-id is required")
	})
}

func TestCancelCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, cancelCmd)
		assert.Equal(t, "cancel", cancelCmd.Use)
		assert.NotEmpty(t, cancelCmd.Short)
		assert.NotEmpty(t, cancelCmd.Long)
	})

	t.Run("has RunE function", func(t *testing.T) {
		assert.NotNil(t, cancelCmd.RunE)
	})

	t.Run("has execution-id flag", func(t *testing.T) {
		flag := cancelCmd.Flags().Lookup("execution-id")
		require.NotNil(t, flag)
		assert.Equal(t, "", flag.DefValue)
	})

	t.Run("requires execution-id", func(t *testing.T) {
		cancelExecutionID = ""
		err := cancelCmd.RunE(cancelCmd, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "--execution-id is required")
	})
}

func TestDeleteCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, deleteCmd)
		assert.Equal(t, "delete", deleteCmd.Use)
		assert.NotEmpty(t, deleteCmd.Short)
		assert.NotEmpty(t, deleteCmd.Long)
	})

	t.Run("has RunE function", func(t *testing.T) {
		assert.NotNil(t, deleteCmd.RunE)
	})

	t.Run("has execution-id flag", func(t *testing.T) {
		flag := deleteCmd.Flags().Lookup("execution-id")
		require.NotNil(t, flag)
		assert.Equal(t, "", flag.DefValue)
	})

	t.Run("has force flag", func(t *testing.T) {
		flag := deleteCmd.Flags().Lookup("force")
		require.NotNil(t, flag)
		assert.Equal(t, "false", flag.DefValue)
	})

	t.Run("requires execution-id", func(t *testing.T) {
		deleteExecutionID = ""
		err := deleteCmd.RunE(deleteCmd, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "--execution-id is required")
	})
}

func TestResultsCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, resultsCmd)
		assert.Equal(t, "results", resultsCmd.Use)
		assert.NotEmpty(t, resultsCmd.Short)
		assert.NotEmpty(t, resultsCmd.Long)
	})

	t.Run("has RunE function", func(t *testing.T) {
		assert.NotNil(t, resultsCmd.RunE)
	})

	t.Run("has execution-id flag", func(t *testing.T) {
		flag := resultsCmd.Flags().Lookup("execution-id")
		require.NotNil(t, flag)
		assert.Equal(t, "", flag.DefValue)
	})

	t.Run("has output flag", func(t *testing.T) {
		flag := resultsCmd.Flags().Lookup("output")
		require.NotNil(t, flag)
		assert.Equal(t, "", flag.DefValue)
	})

	t.Run("has format flag", func(t *testing.T) {
		flag := resultsCmd.Flags().Lookup("format")
		require.NotNil(t, flag)
		assert.Equal(t, "summary", flag.DefValue)
	})

	t.Run("requires execution-id", func(t *testing.T) {
		resultsExecutionID = ""
		err := resultsCmd.RunE(resultsCmd, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "--execution-id is required")
	})
}

func TestStreamCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, streamCmd)
		assert.Equal(t, "stream", streamCmd.Use)
		assert.NotEmpty(t, streamCmd.Short)
		assert.NotEmpty(t, streamCmd.Long)
	})

	t.Run("has RunE function", func(t *testing.T) {
		assert.NotNil(t, streamCmd.RunE)
	})

	t.Run("has execution-id flag", func(t *testing.T) {
		flag := streamCmd.Flags().Lookup("execution-id")
		require.NotNil(t, flag)
		assert.Equal(t, "", flag.DefValue)
	})

	t.Run("requires execution-id", func(t *testing.T) {
		streamExecutionID = ""
		err := streamCmd.RunE(streamCmd, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "--execution-id is required")
	})
}

func TestListCommand(t *testing.T) {
	t.Run("command exists", func(t *testing.T) {
		assert.NotNil(t, listCmd)
		assert.Equal(t, "list", listCmd.Use)
		assert.NotEmpty(t, listCmd.Short)
		assert.NotEmpty(t, listCmd.Long)
	})

	t.Run("has RunE function", func(t *testing.T) {
		assert.NotNil(t, listCmd.RunE)
	})

	t.Run("has status flag", func(t *testing.T) {
		flag := listCmd.Flags().Lookup("status")
		require.NotNil(t, flag)
		assert.Equal(t, "", flag.DefValue)
	})
}

func TestFormatResultsSummary(t *testing.T) {
	t.Run("empty pareto set", func(t *testing.T) {
		pareto := &api.Pareto{
			Vectors: []*api.Vector{},
		}
		result := formatResultsSummary(pareto)
		assert.Contains(t, result, "Pareto Set Results")
		assert.Contains(t, result, "Total vectors: 0")
	})

	t.Run("pareto set with vectors", func(t *testing.T) {
		pareto := &api.Pareto{
			Vectors: []*api.Vector{
				{
					Objectives: []float64{0.1, 0.2},
					Elements:   []float64{0.5, 0.6, 0.7},
				},
				{
					Objectives: []float64{0.3, 0.4},
					Elements:   []float64{0.8, 0.9, 1.0},
				},
			},
			MaxObjs: []float64{1.0, 1.0},
		}
		result := formatResultsSummary(pareto)
		assert.Contains(t, result, "Total vectors: 2")
		assert.Contains(t, result, "Max Objectives:")
		assert.Contains(t, result, "Vector 1:")
		assert.Contains(t, result, "Vector 2:")
		assert.Contains(t, result, "[0.1 0.2]")
		assert.Contains(t, result, "[0.5 0.6 0.7]")
	})

	t.Run("pareto set with more than 10 vectors", func(t *testing.T) {
		vectors := make([]*api.Vector, 15)
		for i := range 15 {
			vectors[i] = &api.Vector{
				Objectives: []float64{float64(i)},
				Elements:   []float64{float64(i)},
			}
		}
		pareto := &api.Pareto{Vectors: vectors}
		result := formatResultsSummary(pareto)
		assert.Contains(t, result, "Total vectors: 15")
		assert.Contains(t, result, "and 5 more vectors")
		assert.Contains(t, result, "Vector 10:")
		assert.NotContains(t, result, "Vector 11:")
	})

	t.Run("pareto set without max objectives", func(t *testing.T) {
		pareto := &api.Pareto{
			Vectors: []*api.Vector{
				{Objectives: []float64{0.5}, Elements: []float64{0.5}},
			},
		}
		result := formatResultsSummary(pareto)
		assert.NotContains(t, result, "Max Objectives:")
	})
}

func TestRegisterCommands(t *testing.T) {
	root := &cobra.Command{Use: "test"}
	RegisterCommands(root)

	commands := root.Commands()
	assert.Len(t, commands, 1)
	assert.Equal(t, "de", commands[0].Use)
}

func TestSetupConfig(t *testing.T) {
	testCfg := config.Default()
	SetupConfig(testCfg)
	assert.Equal(t, testCfg, cfg)
}

func TestSetupStateHandler(t *testing.T) {
	mock := &mockStateOps{}
	SetupStateHandler(mock)
	assert.Equal(t, state.Operations(mock), db)
}

func TestDisplayProgress(t *testing.T) {
	t.Run("in-progress update", func(t *testing.T) {
		progress := &api.StreamProgressResponse{
			CurrentGeneration:   50,
			TotalGenerations:    100,
			CompletedExecutions: 1,
			TotalExecutions:     3,
		}
		// Should not panic
		displayProgress(progress)
	})

	t.Run("completed update", func(t *testing.T) {
		progress := &api.StreamProgressResponse{
			CurrentGeneration:   100,
			TotalGenerations:    100,
			CompletedExecutions: 3,
			TotalExecutions:     3,
		}
		// Should not panic
		displayProgress(progress)
	})

	t.Run("zero values", func(t *testing.T) {
		progress := &api.StreamProgressResponse{}
		// Should not panic with defaults
		displayProgress(progress)
	})
}

func TestGetClientAndContext(t *testing.T) {
	t.Run("returns error when state has no token", func(t *testing.T) {
		db = &mockStateOps{err: assert.AnError}
		cfg = config.Default()

		ctx, client, conn, err := getClientAndContext(t.Context())
		assert.Error(t, err)
		assert.Nil(t, ctx)
		assert.Nil(t, client)
		assert.Nil(t, conn)
	})

	t.Run("creates client with valid token", func(t *testing.T) {
		db = &mockStateOps{token: "test-token"}
		cfg = config.Default()

		ctx, client, conn, err := getClientAndContext(t.Context())
		require.NoError(t, err)
		assert.NotNil(t, ctx)
		assert.NotNil(t, client)
		assert.NotNil(t, conn)
		defer func() { _ = conn.Close() }()
	})
}
