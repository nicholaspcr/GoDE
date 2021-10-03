package models

// Definition definition of the test case functions
type Problem struct {
	Fn          func(e *Vector, M int) error
	ProblemName string
}

func (p *Problem) Name() string {
	return p.ProblemName
}

func (p *Problem) Evaluate(e *Vector, M int) error {
	return p.Fn(e, M)
}
