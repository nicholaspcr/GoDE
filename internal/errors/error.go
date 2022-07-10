// Package errors define custom errors that provide more context to the failures
// found during the execution of the differential evolution procedure.
package errors

import "fmt"

type definition struct {
	code    uint32
	message string
}

func (d *definition) Error() string {
	return d.message
}

func define(code uint32, format string, args ...string) error {
	return &definition{
		code:    code,
		message: fmt.Sprintf(format, args),
	}
}
