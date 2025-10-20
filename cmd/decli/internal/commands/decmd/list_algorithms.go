package decmd

import (
	"fmt"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/spf13/cobra"
)

// listAlgorithmsCmd list all available algorithms.
var listAlgorithmsCmd = &cobra.Command{
	Use:   "list-algorithms",
	Short: "List available algorithms",
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, client, conn, err := getClientAndContext(cmd.Context())
		if err != nil {
			return err
		}
		defer conn.Close()

		res, err := client.ListSupportedAlgorithms(ctx, &api.Empty{})
		if err != nil {
			return err
		}

		fmt.Println(res.GetAlgorithms())
		return nil
	},
}

func init() {
	deCmd.AddCommand(listAlgorithmsCmd)
}
