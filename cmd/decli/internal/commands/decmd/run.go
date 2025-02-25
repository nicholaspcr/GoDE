package decmd

import (
	"fmt"
	"log/slog"

	"github.com/nicholaspcr/GoDE/cmd/decli/internal/config"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	run config.RunConfig
)

// runCmd list all available variants.
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Send run request to server",
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
		res, err := client.Run(ctx, &api.RunRequest{
			Algorithm: run.Algorithm,
			Variant:   run.Algorithm,
			Problem:   run.Problem,
			DeConfig: &api.DEConfig{
				Executions:     run.DeConfig.Executions,
				Generations:    run.DeConfig.Generations,
				PopulationSize: run.DeConfig.PopulationSize,
				DimensionsSize: run.DeConfig.DimensionsSize,
				ObjetivesSize:  run.DeConfig.ObjetivesSize,
				FloorLimiter:   run.DeConfig.FloorLimiter,
				CeilLimiter:    run.DeConfig.CeilLimiter,
				AlgorithmConfig: &api.DEConfig_Gde3{Gde3: &api.GDE3Config{
					Cr: run.DeConfig.GDE3.CR,
					F:  run.DeConfig.GDE3.F,
					P:  run.DeConfig.GDE3.P,
				}},
			},
		})
		if err != nil {
			return err
		}

		slog.With("pareto", res.Pareto).Info("Finished the call")
		return nil
	},
}

func init() {
	deCmd.AddCommand(runCmd)
	fs := runCmd.Flags()

	fs.StringVar(&run.Algorithm, "algorithm", "", "algorithm name")
	fs.StringVar(&run.Variant, "variant", "", "variants password")
	fs.StringVar(&run.Problem, "problem", "", "problem name")

	fs.StringVar(&run.Problem, "executions", "", "amount of executions")
	fs.StringVar(&run.Problem, "generations", "", "amount of generations")
	fs.StringVar(&run.Problem, "population-size", "", "size of the initial population")
	fs.StringVar(&run.Problem, "dimensions-size", "", "amount of elements in a Vector")
	fs.StringVar(&run.Problem, "objectives-size", "", "amount of objectives in a Vector")
	fs.StringVar(&run.Problem, "floor-limiter", "", "minimum value for a Vector's element")
	fs.StringVar(&run.Problem, "ceil-limiter", "", "maximum value for a Vector's element")

	fs.StringVar(&run.Problem, "CR", "", "value of the CR constant")
	fs.StringVar(&run.Problem, "F", "", "value of the F constant")
	fs.StringVar(&run.Problem, "P", "", "value of the P constant")
}
