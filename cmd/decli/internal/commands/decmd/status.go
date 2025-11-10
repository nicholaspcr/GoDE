package decmd

import (
	"fmt"
	"log/slog"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/spf13/cobra"
)

var (
	statusExecutionID string
)

// statusCmd checks the status of an async execution.
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of an async execution",
	Long:  `Retrieve the current status and progress of a previously submitted execution.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		if statusExecutionID == "" {
			return fmt.Errorf("--execution-id is required")
		}

		ctx, client, conn, err := getClientAndContext(cmd.Context())
		if err != nil {
			return err
		}
		defer func() {
			if cerr := conn.Close(); cerr != nil {
				slog.Warn("Failed to close connection", slog.String("error", cerr.Error()))
			}
		}()

		resp, err := client.GetExecutionStatus(ctx, &api.GetExecutionStatusRequest{
			ExecutionId: statusExecutionID,
		})
		if err != nil {
			return fmt.Errorf("failed to get execution status: %w", err)
		}

		execution := resp.Execution
		fmt.Printf("\nExecution ID: %s\n", execution.Id)
		fmt.Printf("Status: %s\n", execution.Status.String())
		if execution.Config != nil {
			fmt.Printf("Generations: %d\n", execution.Config.Generations)
			fmt.Printf("Population Size: %d\n", execution.Config.PopulationSize)
			fmt.Printf("Executions: %d\n", execution.Config.Executions)
		}
		fmt.Printf("Created At: %s\n", execution.CreatedAt.AsTime().Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated At: %s\n", execution.UpdatedAt.AsTime().Format("2006-01-02 15:04:05"))

		if execution.CompletedAt != nil {
			fmt.Printf("Completed At: %s\n", execution.CompletedAt.AsTime().Format("2006-01-02 15:04:05"))
		}

		if execution.Error != "" {
			fmt.Printf("Error: %s\n", execution.Error)
		}

		if resp.Progress != nil {
			fmt.Printf("\nProgress:\n")
			fmt.Printf("  Generation: %d/%d\n", resp.Progress.CurrentGeneration, resp.Progress.TotalGenerations)
			fmt.Printf("  Executions: %d/%d\n", resp.Progress.CompletedExecutions, resp.Progress.TotalExecutions)
		}

		switch execution.Status {
		case api.ExecutionStatus_EXECUTION_STATUS_PENDING:
			fmt.Printf("\nExecution is queued and waiting to start.\n")
		case api.ExecutionStatus_EXECUTION_STATUS_RUNNING:
			fmt.Printf("\nExecution is currently running.\n")
			fmt.Printf("Use 'decli de stream --execution-id %s' for real-time updates\n", statusExecutionID)
		case api.ExecutionStatus_EXECUTION_STATUS_COMPLETED:
			fmt.Printf("\nExecution completed successfully!\n")
			fmt.Printf("Use 'decli de results --execution-id %s' to retrieve results\n", statusExecutionID)
		case api.ExecutionStatus_EXECUTION_STATUS_FAILED:
			fmt.Printf("\nExecution failed.\n")
		case api.ExecutionStatus_EXECUTION_STATUS_CANCELLED:
			fmt.Printf("\nExecution was cancelled.\n")
		}

		return nil
	},
}

func init() {
	deCmd.AddCommand(statusCmd)
	statusCmd.Flags().StringVar(&statusExecutionID, "execution-id", "", "execution ID to check")
}
