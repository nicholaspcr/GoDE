package commands

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/nicholaspcr/gde3/pkg/mode"
	"github.com/nicholaspcr/gde3/pkg/models"
	"github.com/nicholaspcr/gde3/pkg/problems"
	"github.com/nicholaspcr/gde3/pkg/variants"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var scriptCmd = &cobra.Command{
	Use:   "script",
	Short: "runs all the variants for the problem specified",
	Long: `
script is the subcommand responsible for running the gde algorithm into the
specified problem, it will test for all variants and each of them will start
with the same initial population.`,
	Run: func(cmd *cobra.Command, args []string) {
		problem := problems.GetProblemByName(functionName)
		if problem.Name() == "" {
			fmt.Println("invalid problem")
		}

		var params models.AlgorithmParams
		if filename != "" {
			data, err := os.ReadFile(filename)
			if err != nil {
				log.Fatalln("failed to open file")
			}

			yaml.Unmarshal(data, &params)
		} else {
			params = models.AlgorithmParams{
				NP:          np,
				M:           mConst,
				DIM:         dim,
				GEN:         gen,
				EXECS:       execs,
				FLOOR:       floor,
				CEIL:        ceil,
				CR:          crConst,
				F:           fConst,
				P:           pConst,
				DisablePlot: disablePlot,
			}
		}

		// checking for the ceil and floor slices
		if len(params.CEIL) != params.DIM ||
			len(params.FLOOR) != params.DIM {
			log.Fatalln(
				"floor and ceil vector should have the same size as DIM",
			)
		}

		allVariants := variants.GetAllVariants()
		defaultPValues := variants.GetStandardPValues()

		rand.Seed(time.Now().UnixNano())
		initialPopulation := mode.GeneratePopulation(params)

		for _, variant := range allVariants {
			if variant.Name() == "pbest" {
				for _, pvalue := range defaultPValues {
					params.P = pvalue
					mode.MultiExecutions(
						params,
						problem,
						variant,
						initialPopulation,
					)
				}
			} else {
				mode.MultiExecutions(params, problem, variant, initialPopulation)
			}
		}
	},
}
