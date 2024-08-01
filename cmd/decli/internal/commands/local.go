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

		// Initialize the population shared between the executions.
		populationParams := models.PopulationParams{
			DimensionSize:  cfg.Local.Dimensions.Size,
			PopulationSize: cfg.Local.PopulationSize,
			ObjectivesSize: cfg.Local.Constants.M,
			FloorRange:     cfg.Local.Dimensions.Floors,
			CeilRange:      cfg.Local.Dimensions.Ceils,
		}
		initialPopulation, err := models.GeneratePopulation(populationParams)
		if err != nil {
			return err
		}

		problem, err := utils.GetProblemByName(cfg.Local.Problem)
		if err != nil {
			return err
		}
		variant, err := utils.GetVariantByName(cfg.Local.Variant)
		if err != nil {
			return err
		}

		f := de.New(
			de.WithExecutions(cfg.Local.Executions),
			de.WithObjFuncAmount(cfg.Local.Constants.M),
			de.WithAlgorithm(
				gde3.New(
					gde3.WithInitialPopulation(initialPopulation),
					gde3.WithPopulationParams(populationParams),
					gde3.WithConstants(de.Constants{
						F:             cfg.Local.Constants.F,
						P:             cfg.Local.Constants.P,
						CR:            cfg.Local.Constants.CR,
						ObjFuncAmount: populationParams.ObjectivesSize,
						Executions:    cfg.Local.Executions,
						Generations:   cfg.Local.Generations,
						Dimensions:    cfg.Local.Dimensions.Size,
					}),
					gde3.WithProblem(problem),
					gde3.WithVariant(variant),
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
