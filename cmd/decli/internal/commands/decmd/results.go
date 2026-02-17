package decmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/spf13/cobra"
)

var (
	resultsExecutionID string
	resultsOutputFile  string
	resultsFormat      string
)

// resultsCmd retrieves the Pareto set results for a completed execution.
var resultsCmd = &cobra.Command{
	Use:   "results",
	Short: "Retrieve results for a completed execution",
	Long: `Retrieve the Pareto set results for a completed execution.
Results can be displayed as JSON or saved to a file.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		if resultsExecutionID == "" {
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

		resp, err := client.GetExecutionResults(ctx, &api.GetExecutionResultsRequest{
			ExecutionId: resultsExecutionID,
		})
		if err != nil {
			return fmt.Errorf("failed to get results: %w", err)
		}

		if resp.Pareto == nil {
			return fmt.Errorf("no results available for execution %s", resultsExecutionID)
		}

		// Format output based on requested format
		var output string
		switch resultsFormat {
		case "json":
			data, err := json.MarshalIndent(resp.Pareto, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal results to JSON: %w", err)
			}
			output = string(data)

		case "summary":
			output = formatResultsSummary(resp.Pareto)

		default:
			return fmt.Errorf("invalid format: %s (valid: json, summary)", resultsFormat)
		}

		// Write to file or stdout
		if resultsOutputFile != "" {
			if err := os.WriteFile(resultsOutputFile, []byte(output), 0600); err != nil {
				return fmt.Errorf("failed to write results to file: %w", err)
			}
			slog.Info("Results written to file", "file", resultsOutputFile)
			fmt.Printf("Results saved to: %s\n", resultsOutputFile)
		} else {
			fmt.Println(output)
		}

		return nil
	},
}

func formatResultsSummary(pareto *api.Pareto) string {
	var summary strings.Builder
	summary.WriteString("\nPareto Set Results\n")
	summary.WriteString("==================\n\n")
	summary.WriteString(fmt.Sprintf("Total vectors: %d\n", len(pareto.Vectors)))

	if len(pareto.MaxObjs) > 0 {
		summary.WriteString(fmt.Sprintf("\nMax Objectives: %v\n", pareto.MaxObjs))
	}

	summary.WriteString("\nPareto Front Vectors (first 10):\n\n")
	limit := min(len(pareto.Vectors), 10)

	for i := 0; i < limit; i++ {
		v := pareto.Vectors[i]
		summary.WriteString(fmt.Sprintf("Vector %d:\n", i+1))
		summary.WriteString(fmt.Sprintf("  Objectives: %v\n", v.Objectives))
		summary.WriteString(fmt.Sprintf("  Elements: %v\n", v.Elements))
		summary.WriteString("\n")
	}

	if len(pareto.Vectors) > 10 {
		summary.WriteString(fmt.Sprintf("... and %d more vectors\n", len(pareto.Vectors)-10))
		summary.WriteString("\nUse --format json --output results.json to save all results\n")
	}

	return summary.String()
}

func init() {
	deCmd.AddCommand(resultsCmd)
	resultsCmd.Flags().StringVar(&resultsExecutionID, "execution-id", "", "execution ID to get results for")
	resultsCmd.Flags().StringVar(&resultsOutputFile, "output", "", "output file path (default: stdout)")
	resultsCmd.Flags().StringVar(&resultsFormat, "format", "summary", "output format (json, summary)")
}
