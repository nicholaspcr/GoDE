package models

// ProblemFn definition of the test case functions
type ProblemFn struct {
	Fn   func(e *Elem, M int) error
	Name string
}
