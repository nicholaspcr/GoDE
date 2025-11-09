package de

// ModeOptions is a functional option for configuring Differential Evolution execution.
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

// WithProgressCallback sets a callback function to receive progress updates.
func WithProgressCallback(callback ProgressCallback) ModeOptions {
	return func(m *de) *de {
		m.progressCallback = callback
		return m
	}
}
