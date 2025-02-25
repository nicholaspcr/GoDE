package de

// Config of the DE implementation
type Config struct {
	ParetoChannelLimiter int
	MaxChannelLimiter    int
	ResultLimiter        int
}

// Constants are the set of values that determine the behaviour of the Mode
// execution.
type Constants struct {
	// Executions is the amount of times the algorithm will run.
	// All executions start with the same initial population.
	Executions int `json:"executions" yaml:"executions" name:"executions"`

	// Generations of an execution.
	Generations int `json:"generations" yaml:"generations" name:"generations"`

	// Dimensions is the size of the dimensions of each vector element.
	Dimensions int `json:"dimensions" yaml:"dimensions" name:"dimensions"`

	// ObjFuncAmount represents the amount of objective functions.
	ObjFuncAmount int `json:"obj_func_amount" yaml:"obj_func_amount" name:"obj_func_amount"`
}
