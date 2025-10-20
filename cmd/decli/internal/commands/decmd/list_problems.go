package decmd

import (
	"fmt"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/spf13/cobra"
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

		res, err := client.ListSupportedProblems(ctx, &api.Empty{})
		if err != nil {
			return err
		}

		fmt.Println(res.GetProblems())
		return nil
	},
}

func init() {
	deCmd.AddCommand(listProblemsCmd)
}
