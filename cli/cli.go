package cli

import (
	"flag"
	"fmt"
	"os"

	"gitlab.com/nicholaspcr/gode/mo"
	"gitlab.com/nicholaspcr/gode/so"
)

// Start the processing of the flags
func Start() {
	// commands
	sodeCommand := flag.NewFlagSet("(sode | so)\tsubcommand", flag.ExitOnError)
	modeCommand := flag.NewFlagSet("(mode | mo)\tsubcommand", flag.ExitOnError)
	plotCommand := flag.NewFlagSet("(plot | plt)\tsubcommand", flag.ExitOnError)

	_, _, _ = sodeCommand, modeCommand, plotCommand

	// single objective flags
	soNP := sodeCommand.Int("np", 100, "number of elements in the population.")
	soDim := sodeCommand.Int("dim", 5, "number of dimentions of each element.")
	soGen := sodeCommand.Int("gen", 300, "number of generations of the DE.")
	soExecs := sodeCommand.Int("execs", 1, "number of executions.")
	soFloor := sodeCommand.Float64("floor", 0, "floor for the random generated values.")
	soCeil := sodeCommand.Float64("ceil", 1, "ceil for the random generated values.")
	soCR := sodeCommand.Float64("cr", 0.9, "CR -> value used for the DE.")
	soF := sodeCommand.Float64("f", 0.5, "F -> value used for the DE.")
	soP := sodeCommand.Float64("p", 0.5, "P -> value used for the DE.")

	// multiple objetives flags
	moNP := modeCommand.Int("np", 100, "number of elements in the population.")
	moM := modeCommand.Int("m", 3, "value used in the ZDTLs.")
	moDim := modeCommand.Int("dim", 7, "number of dimentions of each element.")
	moGen := modeCommand.Int("gen", 500, "number of generations of the DE.")
	moExecs := modeCommand.Int("execs", 1, "number of executions.")
	moFloor := modeCommand.Float64("floor", 0, "floor for the random generated values.")
	moCeil := modeCommand.Float64("ceil", 1, "ceil for the random generated values.")
	moCR := modeCommand.Float64("cr", 0.9, "CR -> value used for the DE.")
	moF := modeCommand.Float64("f", 0.5, "F -> value used for the DE.")

	// image plotting flags

	if len(os.Args) < 2 {
		fmt.Println("Two commands:\n",
			"\t1: sode or so -> single-objetive DE.\n",
			"\t2: mode or mo -> multi-objetive DE.\n",
			"For more information type the command followed by -h or --help.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "sode":
		sodeCommand.Parse(os.Args[2:])
	case "so":
		sodeCommand.Parse(os.Args[2:])
	case "mode":
		modeCommand.Parse(os.Args[2:])
	case "mo":
		modeCommand.Parse(os.Args[2:])
	case "plot":
		plotCommand.Parse(os.Args[2:])
	case "plt":
		plotCommand.Parse(os.Args[2:])
	default:
		fmt.Println("Two commands:\n",
			"\t1: sode or so -> single-objetive DE.\n",
			"\t2: mode or mo -> multi-objetive DE.\n",
			"For more information type the command followed by -h or --help.")
		os.Exit(1)
	}

	// process parsed subcommand
	if sodeCommand.Parsed() {
		// todo process values
		// e, v := so.Rastrigin, so.Rand1
		// so.DE(*soNP, *soDim, *soGen, *soExecs, *soFloor, *soCeil, *soCR, *soF, *soP, e, v)
		params := so.Params{
			NP:    *soNP,
			DIM:   *soDim,
			GEN:   *soGen,
			EXECS: *soExecs,
			FLOOR: *soFloor,
			CEIL:  *soCeil,
			CR:    *soCR,
			F:     *soF,
			P:     *soP,
		}
		so.Run(params)
	}
	if modeCommand.Parsed() {
		params := mo.Params{
			NP:    *moNP,
			M:     *moM,
			DIM:   *moDim,
			FLOOR: *moFloor,
			CEIL:  *moCeil,
			CR:    *moCR,
			F:     *moF,
			GEN:   *moGen,
		}
		mo.MultiExecutions(*moExecs, params)
	}
	if plotCommand.Parsed() {

	}
}
