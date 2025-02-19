package decmd

import (
	"fmt"

	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// listVariantsCmd list all available variants.
var listVariantsCmd = &cobra.Command{
	Use:   "list-variants",
	Short: "List available variants",
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()

		authToken, err := db.GetAuthToken()
		if err != nil {
			return err
		}
		ctx = metadata.NewOutgoingContext(ctx, metadata.MD{
			"authorization": []string{fmt.Sprintf("Basic %s", authToken)},
		})

		conn, err := grpc.NewClient(
			cfg.Server.GRPCAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return err
		}
		defer conn.Close()

		client := api.NewDifferentialEvolutionServiceClient(conn)
		res, err := client.ListSupportedVariants(ctx, api.Empty)
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
