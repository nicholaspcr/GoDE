package de

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

// ModeOption defines a configuration method for the de struct
type ModeOptions func(*de) *de

// WithAlgorithm sets the algorithm to be used.
func WithAlgorithm(a Algorithm) ModeOptions {
	return func(m *de) *de {
		m.algorithm = a
		return m
	}
}

// WithExecutions determines the amount of times the DE will be
// executed, all executions start with the same population.
func WithExecutions(n int) ModeOptions {
	return func(m *de) *de {
		m.constants.Executions = n
		return m
	}
}

// WithDimensions determines the size of each vector in the population element.
func WithDimensions(dim int) ModeOptions {
	return func(m *de) *de {
		m.constants.Dimensions = dim
		return m
	}
}

// WithGenerations determines the amount of generations of the DE.
func WithGenerations(gen int) ModeOptions {
	return func(m *de) *de {
		m.constants.Generations = gen
		return m
	}
}

// WithObjFuncAmount determines the amount of objective functions.
func WithObjFuncAmount(n int) ModeOptions {
	return func(m *de) *de {
		m.constants.ObjFuncAmount = n
		return m
	}
}
