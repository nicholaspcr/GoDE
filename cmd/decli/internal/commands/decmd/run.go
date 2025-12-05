package decmd

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/nicholaspcr/GoDE/cmd/decli/internal/config"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/spf13/cobra"
)

var (
	run config.RunConfig
)

// runCmd submits an async execution and polls until completion.
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Submit DE execution and wait for results (sync wrapper around async API)",
	Long: `Submit a Differential Evolution execution to the server and wait for it to complete.
This is a synchronous wrapper around the async API that polls for completion.
For async operations, use 'run-async' instead.`,
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
		slog.Info("Submitting execution request...")
		asyncResp, err := client.RunAsync(ctx, &api.RunRequest{
			Algorithm: run.Algorithm,
			Variant:   run.Variant,
			Problem:   run.Problem,
			DeConfig: &api.DEConfig{
				Executions:     run.DeConfig.Executions,
				Generations:    run.DeConfig.Generations,
				PopulationSize: run.DeConfig.PopulationSize,
				DimensionsSize: run.DeConfig.DimensionsSize,
				ObjectivesSize:  run.DeConfig.ObjectivesSize,
				FloorLimiter:   run.DeConfig.FloorLimiter,
				CeilLimiter:    run.DeConfig.CeilLimiter,
				AlgorithmConfig: &api.DEConfig_Gde3{Gde3: &api.GDE3Config{
					Cr: run.DeConfig.GDE3.CR,
					F:  run.DeConfig.GDE3.F,
					P:  run.DeConfig.GDE3.P,
				}},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to submit execution: %w", err)
		}

		executionID := asyncResp.ExecutionId
		slog.Info("Execution submitted", "execution_id", executionID)

		// Poll for completion
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				statusResp, err := client.GetExecutionStatus(ctx, &api.GetExecutionStatusRequest{
					ExecutionId: executionID,
				})
				if err != nil {
					return fmt.Errorf("failed to get execution status: %w", err)
				}

				status := statusResp.Execution.Status
				slog.Info("Execution status", "status", status.String())

				switch status {
				case api.ExecutionStatus_EXECUTION_STATUS_COMPLETED:
					// Get results
					resultsResp, err := client.GetExecutionResults(ctx, &api.GetExecutionResultsRequest{
						ExecutionId: executionID,
					})
					if err != nil {
						return fmt.Errorf("failed to get results: %w", err)
					}

					slog.With("pareto", resultsResp.Pareto).Info("Execution completed successfully")
					return nil

				case api.ExecutionStatus_EXECUTION_STATUS_FAILED:
					return fmt.Errorf("execution failed: %s", statusResp.Execution.Error)

				case api.ExecutionStatus_EXECUTION_STATUS_CANCELLED:
					return fmt.Errorf("execution was cancelled")

				case api.ExecutionStatus_EXECUTION_STATUS_PENDING, api.ExecutionStatus_EXECUTION_STATUS_RUNNING:
					// Continue polling
					if statusResp.Progress != nil {
						slog.Info("Progress update",
							"generation", fmt.Sprintf("%d/%d", statusResp.Progress.CurrentGeneration, statusResp.Progress.TotalGenerations),
							"execution", fmt.Sprintf("%d/%d", statusResp.Progress.CompletedExecutions, statusResp.Progress.TotalExecutions))
					}
				}
			}
		}
	},
}

func init() {
	deCmd.AddCommand(runCmd)
	fs := runCmd.Flags()

	fs.StringVar(&run.Algorithm, "algorithm", "", "algorithm name")
	fs.StringVar(&run.Variant, "variant", "", "variant name")
	fs.StringVar(&run.Problem, "problem", "", "problem name")

	fs.Int64Var(&run.DeConfig.Executions, "executions", 1, "amount of executions")
	fs.Int64Var(&run.DeConfig.Generations, "generations", 100, "amount of generations")
	fs.Int64Var(&run.DeConfig.PopulationSize, "population-size", 100, "size of the initial population")
	fs.Int64Var(&run.DeConfig.DimensionsSize, "dimensions-size", 30, "amount of elements in a Vector")
	fs.Int64Var(&run.DeConfig.ObjectivesSize, "objectives-size", 2, "amount of objectives in a Vector")
	fs.Float32Var(&run.DeConfig.FloorLimiter, "floor-limiter", 0.0, "minimum value for a Vector's element")
	fs.Float32Var(&run.DeConfig.CeilLimiter, "ceil-limiter", 1.0, "maximum value for a Vector's element")

	fs.Float32Var(&run.DeConfig.GDE3.CR, "cr", 0.5, "value of the CR constant")
	fs.Float32Var(&run.DeConfig.GDE3.F, "f", 0.5, "value of the F constant")
	fs.Float32Var(&run.DeConfig.GDE3.P, "p", 0.5, "value of the P constant")
}
