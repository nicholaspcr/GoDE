package decmd

import (
	"fmt"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/spf13/cobra"
)

// listVariantsCmd list all available variants.
var listVariantsCmd = &cobra.Command{
	Use:   "list-variants",
	Short: "List available variants",
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, client, conn, err := getClientAndContext(cmd.Context())
		if err != nil {
			return err
		}
		defer conn.Close()

		res, err := client.ListSupportedVariants(ctx, &api.Empty{})
		if err != nil {
			return err
		}

		fmt.Println(res.GetVariants())
		return nil
	},
}

func init() {
	deCmd.AddCommand(listVariantsCmd)
}
