package commands

import (
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/config"
	"github.com/nicholaspcr/GoDE/internal/errors"
	"github.com/nicholaspcr/GoDE/internal/log"
	"github.com/nicholaspcr/GoDE/pkg/de"
	"github.com/nicholaspcr/GoDE/pkg/de/gde3"
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/spf13/cobra"
)

// modeCmd represents the de command
var modeCmd = &cobra.Command{
	Use:   "multi",
	Short: "Multi-objective implementation of DE",
	Long: `
An implementation that allows the processing of multiple objective functions,
these are a bit more complex and time consuming overall.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		config.ModeLocalFlags()
	},

	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		logger := log.FromContext(ctx)

		logger.Debug("Fetching problem")
		problem := getProblemByName(*config.ProblemName)
		if problem == nil || problem.Name() == "" {
			return errors.DefineProblem(
				"Problem %v not supported",
				*config.ProblemName,
			)
		}

		logger.Debug("Fetching variant")
		variant := getVariantByName(*config.VariantName)
		if variant == nil || variant.Name() == "" {
			return errors.DefineProblem(
				"Variant %v not supported",
				*config.VariantName,
			)
		}

		differentialEvolution := de.New(
			de.WithProblem(problem),
			de.WithVariant(variant),
			// de.WithStore() // TODO: Add store
			de.WithAlgorithm(gde3.New()),

			de.WithPopulation(models.Population{
				Vectors: make(
					[]models.Vector,
					*config.CLI.PopulationSize,
				),
				DimensionsSize: *config.CLI.Dimensions.Size,
				ObjectivesSize: *config.CLI.Constants.M,
				FloorSlice:     *config.CLI.Dimensions.Floors,
				CeilSlice:      *config.CLI.Dimensions.Ceils,
			}),

			de.WithExecs(*config.CLI.Executions),
			de.WithDimensions(*config.CLI.Dimensions.Size),
			de.WithGenerations(*config.CLI.Generations),
			de.WithObjFuncAmount(*config.CLI.Constants.M),

			de.WithFConstant(*config.CLI.Constants.F),
			de.WithPConstant(*config.CLI.Constants.P),
			de.WithCRConstant(*config.CLI.Constants.CR),
		)

		pareto := make(chan models.Population)
		maxObjs := make(chan []float64)
		if err := differentialEvolution.Execute(ctx, pareto, maxObjs); err != nil {
			return err
		}

		return nil
	},
}
