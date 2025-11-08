package decmd

import (
	"fmt"
	"log/slog"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
)

// listProblemsCmd list all available problems.
var listProblemsCmd = &cobra.Command{
	Use:   "list-problems",
	Short: "List available problems",
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

		res, err := client.ListSupportedProblems(ctx, &emptypb.Empty{})
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		if _, err := fmt.Fprintln(w, "Name\tDescription"); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, "----\t----------- "); err != nil {
			return err
		}
		for _, p := range res.GetProblems() {
			if _, err := fmt.Fprintf(w, "%s\t%s\n", p.Name, p.Description); err != nil {
				return err
			}
		}
		return w.Flush()
	},
}

func init() {
	deCmd.AddCommand(listProblemsCmd)
}
