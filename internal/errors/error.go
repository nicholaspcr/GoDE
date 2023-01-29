// Package errors define custom errors that provide more context to the failures
// found during the execution of the differential evolution procedure.
package errors

import (
	"encoding/json"
	"fmt"
)

type definition struct {
	Code    uint32                 `json:"code"`
	Message string                 `json:"message"`
	Fields  map[string]interface{} `json:"fields"`
	Cause   error                  `json:"cause"`
}

// Error returns a marshalled json with all the contents of the error.
func (d *definition) Error() string {
	b, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	return string(b)
}

// WithField adds fields to the error.
func (d *definition) WithField(s string, v interface{}) *definition {
	d.Fields[s] = v
	return d
}

// WithCause saves the internal error which triggers the defined error.
func (d *definition) WithCause(err error) *definition {
	d.Cause = err
	return d
}

func define(code uint32, format string, args ...string) *definition {
	return &definition{
		Code:    code,
		Message: fmt.Sprintf(format, args),
		Fields:  make(map[string]interface{}),
	}
}
