package local

import (
	"github.com/nicholaspcr/GoDE/cmd/decli/internal/config"
	. "github.com/nicholaspcr/GoDE/cmd/decli/internal/utils"
	"github.com/nicholaspcr/GoDE/internal/errors"
	"github.com/nicholaspcr/GoDE/internal/log"
	"github.com/nicholaspcr/GoDE/pkg/de"
	"github.com/nicholaspcr/GoDE/pkg/de/gde3"
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/spf13/cobra"
)

var (
	problemUnsupported   = errors.DefineConfig("problem_unsupported")
	variantUnsupported   = errors.DefineConfig("variant_unsupported")
	algorithmUnsupported = errors.DefineConfig("algorithm_unsupported")
)

// RunCmd represents the de command
var RunCmd = &cobra.Command{
	Use:     "run",
	Aliases: []string{"start"},
	Short:   "Runs local multi-objective implementation of DE",
	Long: `
An implementation that allows the processing of multiple objective functions,
these are a bit more complex and time consuming overall.

Specify the algorithm via argument, for example: 'decli local run gde3'
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Help()
		}

		ctx := cmd.Context()
		logger := log.FromContext(ctx)

		logger.Debug("Fetching problem")
		problem := GetProblemByName(*config.ProblemName)
		if problem == nil || problem.Name() == "" {
			return problemUnsupported.WithField("problem", *config.ProblemName)
		}

		logger.Debug("Fetching variant")
		variant := GetVariantByName(*config.VariantName)
		if variant == nil || variant.Name() == "" {
			return variantUnsupported.WithField("variant", *config.VariantName)
		}

		deOpts := []de.ModeOptions{
			// de.WithStore() // TODO: Add store
			de.WithProblem(problem),
			de.WithVariant(variant),
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
		}

		algo := args[0]
		switch algo {
		case "gde3":
			deOpts = append(deOpts, de.WithAlgorithm(gde3.New()))
		default:
			return algorithmUnsupported.WithField("algorithm", algo)
		}

		pareto := make(chan models.Population)
		maxObjs := make(chan []float64)
		if err := de.New(deOpts...).Execute(ctx, pareto, maxObjs); err != nil {
			return err
		}

		logger.Info("Maximum Objectives:")
		idx := 0
		for maxi := range maxObjs {
			logger.Info("Maximum Objectives :: ", idx, " :: ", maxi)
			idx++
		}

		logger.Info("Pareto:")
		idx = 0
		for p := range pareto {
			logger.Info("Pareto :: ", p.Vectors)
			idx++
		}

		return nil
	},
}

func init() {
	LocalCmd.AddCommand(RunCmd)
	config.ModeFlags(RunCmd.Flags())
}
