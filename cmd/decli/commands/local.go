package commands

import (
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/utils"
	"github.com/nicholaspcr/GoDE/pkg/de"
	"github.com/nicholaspcr/GoDE/pkg/de/gde3"
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/spf13/cobra"
)

// localCmd represents the de command
var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Local operations related to Differential Evolutionary algorithm",
	RunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.Help()
	},
}

// localRunCmd represents the run command for local operations.
var localRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a local Differential Evolutionary algorithm",
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		f := de.New(
			de.WithExecutions(cfg.Executions),
			de.WithAlgorithm(
				gde3.New(
					gde3.WithPopulationParams(models.PopulationParams{
						DimensionSize:  cfg.Dimensions.Size,
						PopulationSize: cfg.PopulationSize,
						ObjectivesSize: cfg.Constants.M,
						FloorRange:     cfg.Dimensions.Floors,
						CeilRange:      cfg.Dimensions.Ceils,
					}),
					gde3.WithProblem(utils.GetProblemByName(cfg.Problem)),
					gde3.WithVariant(utils.GetVariantByName(cfg.Variant)),
				),
			),
		)

		if err := f.Execute(ctx); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	localCmd.AddCommand(localRunCmd)
}
