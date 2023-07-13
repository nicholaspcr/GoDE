package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const appname = "decli"

func init() {
	// config filename and type
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// config search path
	viper.AddConfigPath("/etc/decli/")
	viper.AddConfigPath("$HOME/.decli")
	viper.AddConfigPath(".")

	// fetch config file contents
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

type (
	// Config is a set of values that are necessary to execute an Differential
	// Evolutionary algorithm.
	Config struct {
		PopulationSize int        `name:"population_size" json:"population_size" yaml:"population_size"`
		Generations    int        `name:"generations" json:"generations" yaml:"generations"`
		SaveEachGen    bool       `name:"save_each_gen" json:"save_each_gen" yaml:"save_each_gen"`
		Executions     int        `name:"executions" json:"executions" yaml:"executions"`
		Dimensions     Dimensions `name:"dimensions" json:"dimensions" yaml:"dimensions"`
		Constants      Constants  `name:"constants" json:"constants" yaml:"constants"`
		Problem        string     `name:"problem" json:"problem" yaml:"problem"`
		Variant        string     `name:"variant" json:"variant" yaml:"variant"`
	}

	// Dimensions is a set of values to define the behaviour that happens in
	// each dimension of the DE.
	Dimensions struct {
		Size   int       `name:"size"   json:"size"   yaml:"size"`
		Floors []float64 `name:"floors" json:"floors" yaml:"floors"`
		Ceils  []float64 `name:"ceils"  json:"ceils"  yaml:"ceils"`
	}

	// Constants is a set of values to define the behaviour of a DE.
	Constants struct {
		M  int     `name:"m"  json:"m"  yaml:"m"`
		CR float64 `name:"cr" json:"cr" yaml:"cr"`
		F  float64 `name:"f"  json:"f"  yaml:"f"`
		P  float64 `name:"p"  json:"p"  yaml:"p"`
	}
)

var (
	DefaultConfig = Config{
		PopulationSize: 50,
		Generations:    100,
		SaveEachGen:    false,
		Executions:     1,
		Dimensions: Dimensions{
			Size:   7,
			Floors: []float64{0, 0, 0, 0, 0, 0, 0},
			Ceils:  []float64{1, 1, 1, 1, 1, 1, 1},
		},
		Constants: Constants{
			M:  int(3),
			CR: float64(0.9),
			F:  float64(0.5),
			P:  float64(0.2),
		},
		Problem: "dtlz1",
		Variant: "rand1",
	}
)

// Unmarshal reads the configuration from the config file and unmarshals it into
// the given pointer.
func Unmarshal(ptr any, opts ...viper.DecoderConfigOption) error {
	return viper.Unmarshal(ptr, opts...)
}

// bindFlags binds the flags to the configuration.
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := f.Name
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

// InitializeRoot initializes the configuration for the root command.
func InitializeRoot(cmd *cobra.Command) error {
	v := viper.New()

	// config filename and type
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// config search path
	v.AddConfigPath("/etc/decli/")
	v.AddConfigPath("$HOME/.decli")
	v.AddConfigPath(".")

	// fetch config file contents
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	v.AutomaticEnv()
	bindFlags(cmd, v)

	return nil
}
