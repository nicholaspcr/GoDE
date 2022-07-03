package commands

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/nicholaspcr/GoDE/pkg/de"
	"github.com/nicholaspcr/GoDE/pkg/models"
	"github.com/nicholaspcr/GoDE/pkg/variants"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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

	Run: func(cmd *cobra.Command, args []string) {
		problem := getProblemByName(functionName)
		variant := variants.GetVariantByName(variantName)

		if problem.Name() == "" {
			fmt.Println("Invalid problem")
			return
		}

		if variant.Name() == "" {
			fmt.Println("Invalid variant.")
			return
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
			fmt.Println(
				"floor and ceil vector should have the same size as DIM",
			)
			fmt.Println("ceil = ", params.CEIL)
			fmt.Println("floor  = ", params.FLOOR)
			fmt.Println("dim = ", params.DIM)
			return
		}
		startTimer := time.Now() // time spent on script

		rand.Seed(time.Now().UnixNano())
		// generating shared initial population
		initialPopulation := de.GeneratePopulation(params)

		de.MultiExecutions(params, problem, variant, initialPopulation)

		timeSpent := time.Since(startTimer)
		fmt.Println("Time spend on the script: ", timeSpent)
	},
}
