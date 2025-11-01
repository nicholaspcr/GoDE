package decmd

import (
	"fmt"
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
		defer conn.Close()

		res, err := client.ListSupportedProblems(ctx, &emptypb.Empty{})
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintln(w, "Name\tDescription")
		fmt.Fprintln(w, "----\t----------- ")
		for _, p := range res.GetProblems() {
			fmt.Fprintf(w, "%s\t%s\n", p.Name, p.Description)
		}
		return w.Flush()
	},
}

func init() {
	deCmd.AddCommand(listProblemsCmd)
}
