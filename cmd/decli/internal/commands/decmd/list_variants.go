package decmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
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

		res, err := client.ListSupportedVariants(ctx, &emptypb.Empty{})
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintln(w, "Name\tDescription")
		fmt.Fprintln(w, "----\t----------- ")
		for _, v := range res.GetVariants() {
			fmt.Fprintf(w, "%s\t%s\n", v.Name, v.Description)
		}
		return w.Flush()
	},
}

func init() {
	deCmd.AddCommand(listVariantsCmd)
}
