package gde3

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/nicholaspcr/gde3/pkg/mode"
	"github.com/nicholaspcr/gde3/pkg/problems"
	"github.com/nicholaspcr/gde3/pkg/problems/models"
	"github.com/nicholaspcr/gde3/pkg/variants"
	"github.com/spf13/cobra"
)

// local flags
var variantName string

// modeCmd represents the mode command
var modeCmd = &cobra.Command{
	Use:   "multi",
	Short: "Multi-objective implementation of DE",
	Long:  `An implementation that allows the processing of multiple objective functions, these are a bit more complex and time consuming overall.`,

	Run: func(cmd *cobra.Command, args []string) {
		problem := problems.GetProblemByName(functionName)
		variant := variants.GetVariantByName(variantName)
		if problem.Name == "" {
			fmt.Println("Invalid problem")
			return
		}
		if variant.Name == "" {
			fmt.Println("Invalid variant.")
			return
		}
		params := models.Params{
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
		if cpuprofile != "" {
			cpuF, err := os.Create(cpuprofile)
			if err != nil {
				log.Fatal("could not create CPU profile: ", err)
			}
			defer cpuF.Close() // error handling omitted for example
			if err := pprof.StartCPUProfile(cpuF); err != nil {
				log.Fatal("could not start CPU profile: ", err)
			}
			defer pprof.StopCPUProfile()

		}

		startTimer := time.Now() // time spent on script

		rand.Seed(time.Now().UnixNano())
		// generating shared initial population
		initialPopulation := mode.GeneratePopulation(params)

		mode.MultiExecutions(params, problem, variant, initialPopulation)

		timeSpent := time.Since(startTimer)
		fmt.Println("Time spend on the script: ", timeSpent)

		if memprofile != "" {
			memF, err := os.Create(memprofile)
			if err != nil {
				log.Fatal("could not create memory profile: ", err)
			}
			defer memF.Close() // error handling omitted for example
			runtime.GC()       // get up-to-date statistics
			if err := pprof.WriteHeapProfile(memF); err != nil {
				log.Fatal("could not write memory profile: ", err)
			}
		}
	},
}
