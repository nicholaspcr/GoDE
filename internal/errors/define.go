package errors

func DefineConfig(format string, args ...string) error {
	return define(configuration, format, args...)
}

func DefineVariant(format string, args ...string) error {
	return define(variant, format, args...)
}

func DefineProblem(format string, args ...string) error {
	return define(problem, format, args...)
}

func DefineAlgorithm(format string, args ...string) error {
	return define(algorithm, format, args...)
}

