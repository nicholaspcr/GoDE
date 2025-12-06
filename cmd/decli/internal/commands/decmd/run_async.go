package decmd

import (
	"fmt"
	"log/slog"

	"github.com/nicholaspcr/GoDE/cmd/decli/internal/config"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/spf13/cobra"
)

var (
	runAsync config.RunConfig
)

// runAsyncCmd submits an async execution and returns immediately with execution ID.
var runAsyncCmd = &cobra.Command{
	Use:   "run-async",
	Short: "Submit DE execution and return immediately with execution ID",
	Long: `Submit a Differential Evolution execution to the server asynchronously.
Returns immediately with an execution ID that can be used to check status, stream progress, or retrieve results.
For synchronous operation (submit + wait), use 'run' instead.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, client, conn, err := getClientAndContext(cmd.Context())
		if err != nil {
			return err
		}
		defer func() {
			if cerr := conn.Close(); cerr != nil {
				slog.Warn("Failed to close connection", slog.String("error", cerr.Error()))
			}
		}()

		// Submit async execution
		slog.Info("Submitting async execution request...")
		resp, err := client.RunAsync(ctx, &api.RunAsyncRequest{
			Algorithm: runAsync.Algorithm,
			Variant:   runAsync.Variant,
			Problem:   runAsync.Problem,
			DeConfig: &api.DEConfig{
				Executions:     runAsync.DeConfig.Executions,
				Generations:    runAsync.DeConfig.Generations,
				PopulationSize: runAsync.DeConfig.PopulationSize,
				DimensionsSize: runAsync.DeConfig.DimensionsSize,
				ObjectivesSize:  runAsync.DeConfig.ObjectivesSize,
				FloorLimiter:   runAsync.DeConfig.FloorLimiter,
				CeilLimiter:    runAsync.DeConfig.CeilLimiter,
				AlgorithmConfig: &api.DEConfig_Gde3{Gde3: &api.GDE3Config{
					Cr: runAsync.DeConfig.GDE3.CR,
					F:  runAsync.DeConfig.GDE3.F,
					P:  runAsync.DeConfig.GDE3.P,
				}},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to submit execution: %w", err)
		}

		slog.Info("Execution submitted successfully",
			"execution_id", resp.ExecutionId)

		fmt.Printf("\nExecution ID: %s\n", resp.ExecutionId)
		fmt.Printf("\nUse 'decli de status --execution-id %s' to check progress\n", resp.ExecutionId)
		fmt.Printf("Use 'decli de stream --execution-id %s' to stream real-time updates\n", resp.ExecutionId)

		return nil
	},
}

func init() {
	deCmd.AddCommand(runAsyncCmd)
	fs := runAsyncCmd.Flags()

	fs.StringVar(&runAsync.Algorithm, "algorithm", "", "algorithm name (required)")
	fs.StringVar(&runAsync.Variant, "variant", "", "variant name (required)")
	fs.StringVar(&runAsync.Problem, "problem", "", "problem name (required)")

	_ = runAsyncCmd.MarkFlagRequired("algorithm")
	_ = runAsyncCmd.MarkFlagRequired("variant")
	_ = runAsyncCmd.MarkFlagRequired("problem")

	fs.Int64Var(&runAsync.DeConfig.Executions, "executions", 1, "amount of executions")
	fs.Int64Var(&runAsync.DeConfig.Generations, "generations", 100, "amount of generations")
	fs.Int64Var(&runAsync.DeConfig.PopulationSize, "population-size", 100, "size of the initial population")
	fs.Int64Var(&runAsync.DeConfig.DimensionsSize, "dimensions-size", 30, "amount of elements in a Vector")
	fs.Int64Var(&runAsync.DeConfig.ObjectivesSize, "objectives-size", 2, "amount of objectives in a Vector")
	fs.Float32Var(&runAsync.DeConfig.FloorLimiter, "floor-limiter", 0.0, "minimum value for a Vector's element")
	fs.Float32Var(&runAsync.DeConfig.CeilLimiter, "ceil-limiter", 1.0, "maximum value for a Vector's element")

	fs.Float32Var(&runAsync.DeConfig.GDE3.CR, "cr", 0.5, "value of the CR constant")
	fs.Float32Var(&runAsync.DeConfig.GDE3.F, "f", 0.5, "value of the F constant")
	fs.Float32Var(&runAsync.DeConfig.GDE3.P, "p", 0.5, "value of the P constant")
}
