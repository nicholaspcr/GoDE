package decmd

import (
	"fmt"
	"log/slog"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/spf13/cobra"
)

var (
	listStatus string
)

// listCmd lists all executions for the current user.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all executions for the current user",
	Long: `Retrieve a list of all executions submitted by the current user.
Optionally filter by status (pending, running, completed, failed, cancelled).`,
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

		req := &api.ListExecutionsRequest{}

		// Parse and set status filter if provided
		if listStatus != "" {
			statusValue, ok := api.ExecutionStatus_value[fmt.Sprintf("EXECUTION_STATUS_%s", listStatus)]
			if !ok {
				return fmt.Errorf("invalid status: %s (valid: pending, running, completed, failed, cancelled)", listStatus)
			}
			req.Status = api.ExecutionStatus(statusValue)
		}

		resp, err := client.ListExecutions(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to list executions: %w", err)
		}

		if len(resp.Executions) == 0 {
			fmt.Println("No executions found.")
			return nil
		}

		fmt.Printf("\nFound %d execution(s):\n\n", len(resp.Executions))

		for i, exec := range resp.Executions {
			fmt.Printf("%d. Execution ID: %s\n", i+1, exec.Id)
			fmt.Printf("   Status: %s\n", exec.Status.String())
			if exec.Config != nil {
				fmt.Printf("   Generations: %d, Population: %d\n", exec.Config.Generations, exec.Config.PopulationSize)
			}
			fmt.Printf("   Created: %s\n", exec.CreatedAt.AsTime().Format("2006-01-02 15:04:05"))

			if exec.CompletedAt != nil {
				fmt.Printf("   Completed: %s\n", exec.CompletedAt.AsTime().Format("2006-01-02 15:04:05"))
			}

			if exec.Error != "" {
				fmt.Printf("   Error: %s\n", exec.Error)
			}

			fmt.Println()
		}

		return nil
	},
}

func init() {
	deCmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&listStatus, "status", "", "filter by status (pending, running, completed, failed, cancelled)")
}
