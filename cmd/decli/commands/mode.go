package commands

import (
	"github.com/nicholaspcr/GoDE/internal/errors"
	"github.com/nicholaspcr/GoDE/internal/log"
	"github.com/nicholaspcr/GoDE/pkg/de"
	"github.com/nicholaspcr/GoDE/pkg/de/gde3"
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/spf13/cobra"
)

// local flags
var variantName string

// modeCmd represents the de command
var modeCmd = &cobra.Command{
	Use:   "multi",
	Short: "Multi-objective implementation of DE",
	Long: `
An implementation that allows the processing of multiple objective functions,
these are a bit more complex and time consuming overall.`,

	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		logger := log.FromContext(ctx)

		logger.Debug("Fetching problem")
		problem := getProblemByName(problemName)
		if problem == nil || problem.Name() == "" {
			return errors.DefineProblem("Problem %v not supported", problemName)
		}

		logger.Debug("Fetching variant")
		variant := getVariantByName(variantName)
		if variant == nil || variant.Name() == "" {
			return errors.DefineProblem("Variant %v not supported", variantName)
		}

		differentialEvolution := de.New(
			de.WithProblem(problem),
			de.WithVariant(variant),
			// de.WithStore() // TODO: Add store
			de.WithAlgorithm(gde3.New()),

			de.WithPopulation(models.Population{
				Vectors:        make([]models.Vector, np),
				DimensionsSize: dim,
				ObjectivesSize: mConst,
				FloorSlice:     floor,
				CeilSlice:      ceil,
			}),

			de.WithExecs(execs),
			de.WithDimensions(dim),
			de.WithGenerations(gen),
			de.WithObjFuncAmount(mConst),

			de.WithFConstant(fConst),
			de.WithPConstant(pConst),
			de.WithCRConstant(crConst),
		)

		pareto := make(chan models.Population)
		maxObjs := make(chan []float64)
		if err := differentialEvolution.Execute(ctx, pareto, maxObjs); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	modeCmd.Flags().StringVar(
		&variantName,
		"vr",
		"rand1",
		"name fo the variant to be used",
	)
}
