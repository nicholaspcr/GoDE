package decmd

import (
	"fmt"
	"log/slog"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/spf13/cobra"
)

var (
	deleteExecutionID string
	deleteForce       bool
)

// deleteCmd deletes an execution record.
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an execution record",
	Long: `Delete an execution record from the server.
This will remove the execution metadata and any associated results.
Running executions must be cancelled first, or use --force to cancel and delete.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		if deleteExecutionID == "" {
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

		// If force flag is set, try to cancel first
		if deleteForce {
			slog.Info("Force delete requested, attempting to cancel execution first", "execution_id", deleteExecutionID)

			// Check status first
			statusResp, err := client.GetExecutionStatus(ctx, &api.GetExecutionStatusRequest{
				ExecutionId: deleteExecutionID,
			})
			if err == nil && (statusResp.Execution.Status == api.ExecutionStatus_EXECUTION_STATUS_RUNNING ||
				statusResp.Execution.Status == api.ExecutionStatus_EXECUTION_STATUS_PENDING) {

				// Cancel the execution
				_, err = client.CancelExecution(ctx, &api.CancelExecutionRequest{
					ExecutionId: deleteExecutionID,
				})
				if err != nil {
					slog.Warn("Failed to cancel execution before deletion", "error", err.Error())
				} else {
					fmt.Printf("Execution cancelled before deletion\n")
				}
			}
		}

		// Delete the execution
		slog.Info("Deleting execution", "execution_id", deleteExecutionID)

		_, err = client.DeleteExecution(ctx, &api.DeleteExecutionRequest{
			ExecutionId: deleteExecutionID,
		})
		if err != nil {
			return fmt.Errorf("failed to delete execution: %w", err)
		}

		fmt.Printf("\nExecution deleted: %s\n", deleteExecutionID)

		return nil
	},
}

func init() {
	deCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&deleteExecutionID, "execution-id", "", "execution ID to delete")
	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "cancel execution before deletion if still running")
}
