package decmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"
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
		defer func() {
			if cerr := conn.Close(); cerr != nil {
				slog.Warn("Failed to close connection", slog.String("error", cerr.Error()))
			}
		}()

		res, err := client.ListSupportedAlgorithms(ctx, &emptypb.Empty{})
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
