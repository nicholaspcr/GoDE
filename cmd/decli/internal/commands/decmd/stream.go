package decmd

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/spf13/cobra"
)

var (
	streamExecutionID string
)

// streamCmd streams real-time progress updates for a running execution.
var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Stream real-time progress updates for an execution",
	Long: `Stream real-time progress updates for a running execution.
Updates will be displayed as they are received until the execution completes or is cancelled.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		if streamExecutionID == "" {
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

		// Create a cancellable context for streaming
		streamCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		slog.Info("Starting progress stream", "execution_id", streamExecutionID)
		fmt.Printf("\nStreaming progress for execution: %s\n", streamExecutionID)
		fmt.Printf("Press Ctrl+C to stop streaming\n\n")

		stream, err := client.StreamProgress(streamCtx, &api.StreamProgressRequest{
			ExecutionId: streamExecutionID,
		})
		if err != nil {
			return fmt.Errorf("failed to start stream: %w", err)
		}

		// Receive and display progress updates
		for {
			progress, err := stream.Recv()
			if err == io.EOF {
				fmt.Printf("\nStream ended.\n")
				break
			}
			if err != nil {
				return fmt.Errorf("stream error: %w", err)
			}

			displayProgress(progress)
		}

		return nil
	},
}

func displayProgress(progress *api.StreamProgressResponse) {
	fmt.Printf("\r[Generation %d/%d] [Execution %d/%d]",
		progress.GetCurrentGeneration(),
		progress.GetTotalGenerations(),
		progress.GetCompletedExecutions(),
		progress.GetTotalExecutions())

	// Add newline if execution is complete
	if progress.GetCompletedExecutions() == progress.GetTotalExecutions() &&
		progress.GetCurrentGeneration() == progress.GetTotalGenerations() {
		fmt.Println()
		fmt.Println("\nExecution completed!")
	}
}

func init() {
	deCmd.AddCommand(streamCmd)
	streamCmd.Flags().StringVar(&streamExecutionID, "execution-id", "", "execution ID to stream progress for")
}
