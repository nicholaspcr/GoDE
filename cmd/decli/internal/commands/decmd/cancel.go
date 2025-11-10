package decmd

import (
	"fmt"
	"log/slog"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/spf13/cobra"
)

var (
	cancelExecutionID string
)

// cancelCmd cancels a running execution.
var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel a running execution",
	Long: `Request cancellation of a running execution.
The execution will be marked for cancellation and will stop as soon as possible.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		if cancelExecutionID == "" {
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

		slog.Info("Requesting cancellation", "execution_id", cancelExecutionID)

		_, err = client.CancelExecution(ctx, &api.CancelExecutionRequest{
			ExecutionId: cancelExecutionID,
		})
		if err != nil {
			return fmt.Errorf("failed to cancel execution: %w", err)
		}

		fmt.Printf("\nCancellation requested for execution: %s\n", cancelExecutionID)
		fmt.Printf("The execution will stop as soon as possible.\n")
		fmt.Printf("\nUse 'decli de status --execution-id %s' to verify cancellation\n", cancelExecutionID)

		return nil
	},
}

func init() {
	deCmd.AddCommand(cancelCmd)
	cancelCmd.Flags().StringVar(&cancelExecutionID, "execution-id", "", "execution ID to cancel")
}
