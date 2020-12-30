package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/nicholaspcr/go-de/mo"
)

// local flags
var mConst int
var functionName string
var variantName string

// modeCmd represents the mode command
var modeCmd = &cobra.Command{
	Use:   "multi",
	Short: "Multi-objective implementation of DE",
	Long:  `An implementation that allows the processing of multiple objective functions, these are a bit more complex and time consuming overall.`,

	Run: func(cmd *cobra.Command, args []string) {
		problem := mo.GetProblemByName(functionName)
		variant := mo.GetVariantByName(variantName)
		if problem == nil || variant == nil {
			// todo error message and list of problems and variants
			fmt.Println("Invalid problem or variant.")
			return
		}
		params := mo.Params{
			NP:    np,
			M:     mConst,
			DIM:   dim,
			GEN:   gen,
			EXECS: execs,
			FLOOR: floor,
			CEIL:  ceil,
			CR:    crConst,
			F:     fConst,
		}
		mo.MultiExecutions(params, problem, variant)
	},
}

func init() {
	rootCmd.AddCommand(modeCmd)
	modeCmd.Flags().IntVar(&mConst,
		"M",
		3,
		"M -> DE constant")
	modeCmd.Flags().StringVar(&functionName,
		"fn",
		"DTLZ1",
		"name of the problem to be used.")
	modeCmd.Flags().StringVar(&variantName,
		"vr",
		"rand1",
		"name fo the variant to be used")
}
