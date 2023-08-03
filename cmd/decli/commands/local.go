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
			DimensionSize:  cfg.Dimensions.Size,
			PopulationSize: cfg.PopulationSize,
			ObjectivesSize: cfg.Constants.M,
			FloorRange:     cfg.Dimensions.Floors,
			CeilRange:      cfg.Dimensions.Ceils,
		}
		initialPopulation, err := models.GeneratePopulation(populationParams)
		if err != nil {
			return err
		}

		problem, err := utils.GetProblemByName(cfg.Problem)
		if err != nil {
			return err
		}
		variant, err := utils.GetVariantByName(cfg.Variant)
		if err != nil {
			return err
		}

		f := de.New(
			de.WithExecutions(cfg.Executions),
			de.WithObjFuncAmount(cfg.Constants.M),
			de.WithAlgorithm(
				gde3.New(
					gde3.WithInitialPopulation(initialPopulation),
					gde3.WithPopulationParams(populationParams),
					gde3.WithConstants(de.Constants{
						F:             cfg.Constants.F,
						P:             cfg.Constants.P,
						CR:            cfg.Constants.CR,
						ObjFuncAmount: populationParams.ObjectivesSize,
						Executions:    cfg.Executions,
						Generations:   cfg.Generations,
						Dimensions:    cfg.Dimensions.Size,
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
