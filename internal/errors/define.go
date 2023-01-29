package errors

func DefineConfig(format string, args ...string) *definition {
	return define(configuration, format, args...)
}
