package decmd

import (
	"log/slog"

	"github.com/nicholaspcr/GoDE/cmd/decli/internal/config"
	"github.com/nicholaspcr/GoDE/pkg/api/v1"
	"github.com/spf13/cobra"
)

var (
	run config.RunConfig
)

// runCmd list all available variants.
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Send run request to server",
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

		res, err := client.Run(ctx, &api.RunRequest{
			Algorithm: run.Algorithm,
			Variant:   run.Variant,
			Problem:   run.Problem,
			DeConfig: &api.DEConfig{
				Executions:     run.DeConfig.Executions,
				Generations:    run.DeConfig.Generations,
				PopulationSize: run.DeConfig.PopulationSize,
				DimensionsSize: run.DeConfig.DimensionsSize,
				ObjetivesSize:  run.DeConfig.ObjectivesSize,
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
	fs.StringVar(&run.Variant, "variant", "", "variant name")
	fs.StringVar(&run.Problem, "problem", "", "problem name")

	fs.Int64Var(&run.DeConfig.Executions, "executions", 1, "amount of executions")
	fs.Int64Var(&run.DeConfig.Generations, "generations", 100, "amount of generations")
	fs.Int64Var(&run.DeConfig.PopulationSize, "population-size", 100, "size of the initial population")
	fs.Int64Var(&run.DeConfig.DimensionsSize, "dimensions-size", 30, "amount of elements in a Vector")
	fs.Int64Var(&run.DeConfig.ObjectivesSize, "objectives-size", 2, "amount of objectives in a Vector")
	fs.Float32Var(&run.DeConfig.FloorLimiter, "floor-limiter", 0.0, "minimum value for a Vector's element")
	fs.Float32Var(&run.DeConfig.CeilLimiter, "ceil-limiter", 1.0, "maximum value for a Vector's element")

	fs.Float32Var(&run.DeConfig.GDE3.CR, "cr", 0.5, "value of the CR constant")
	fs.Float32Var(&run.DeConfig.GDE3.F, "f", 0.5, "value of the F constant")
	fs.Float32Var(&run.DeConfig.GDE3.P, "p", 0.5, "value of the P constant")
}
