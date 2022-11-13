package config

type (
	// Config is a set of values that are necessary to execute an Differential
	// Evolutionary algorithm.
	Config struct {
		PopulationSize *int        `name:"population_size" json:"population_size" yaml:"population_size"`
		Generations    *int        `name:"generations"     json:"generations"     yaml:"generations"`
		SaveEachGen    *bool       `name:"save_each_gen"   json:"save_each_gen"   yaml:"save_each_gen"`
		Executions     *int        `name:"executions"      json:"executions"      yaml:"executions"`
		Dimensions     *Dimensions `name:"dimensions"      json:"dimensions"      yaml:"dimensions"`
		Constants      *Constants  `name:"constants"       json:"constants"       yaml:"constants"`
	}

	// Dimensions is a set of values to define the behaviour that happens in
	// each dimension of the DE.
	Dimensions struct {
		Size   *int       `name:"size"   json:"size"   yaml:"size"`
		Floors *[]float64 `name:"floors" json:"floors" yaml:"floors"`
		Ceils  *[]float64 `name:"ceils"  json:"ceils"  yaml:"ceils"`
	}

	// Constants is a set of values to define the behaviour of a DE.
	Constants struct {
		M  *int     `name:"m"  json:"m"  yaml:"m"`
		CR *float64 `name:"cr" json:"cr" yaml:"cr"`
		F  *float64 `name:"f"  json:"f"  yaml:"f"`
		P  *float64 `name:"p"  json:"p"  yaml:"p"`
	}
)
