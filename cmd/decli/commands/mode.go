package commands

import (
	"github.com/nicholaspcr/GoDE/internal/errors"
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
		problem := getProblemByName(functionName)
		variant := getVariantByName(variantName)

		if problem.Name() == "" {
			return errors.DefineConfig("Invalid problem")
		}

		if variant.Name() == "" {
			return errors.DefineConfig("Invalid variant.")
		}

		// TODO: Validate values passed, or leave it to the server?

		//// checking for the ceil and floor slices
		//if len(params.CEIL) != params.DIM ||
		//	len(params.FLOOR) != params.DIM {
		//	fmt.Println(
		//		"floor and ceil vector should have the same size as DIM",
		//	)
		//	fmt.Println("ceil = ", params.CEIL)
		//	fmt.Println("floor  = ", params.FLOOR)
		//	fmt.Println("dim = ", params.DIM)
		//	return
		//}

		ctx := cmd.Context()
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
