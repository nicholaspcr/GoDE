package wfg

import (
	"fmt"

	"github.com/nicholaspcr/GoDE/pkg/problems"
)

//nolint:revive // Factory functions have unused parameters matching registry interface
func init() {
	// Register WFG problems
	for i := 1; i <= 9; i++ {
		num := i // Capture for closure
		name := fmt.Sprintf("wfg%d", num)

		var factory problems.ProblemFactory
		switch num {
		case 1:
			factory = func(dim, objs int) (problems.Interface, error) { return Wfg1(), nil }
		case 2:
			factory = func(dim, objs int) (problems.Interface, error) { return Wfg2(), nil }
		case 3:
			factory = func(dim, objs int) (problems.Interface, error) { return Wfg3(), nil }
		case 4:
			factory = func(dim, objs int) (problems.Interface, error) { return Wfg4(), nil }
		case 5:
			factory = func(dim, objs int) (problems.Interface, error) { return Wfg5(), nil }
		case 6:
			factory = func(dim, objs int) (problems.Interface, error) { return Wfg6(), nil }
		case 7:
			factory = func(dim, objs int) (problems.Interface, error) { return Wfg7(), nil }
		case 8:
			factory = func(dim, objs int) (problems.Interface, error) { return Wfg8(), nil }
		case 9:
			factory = func(dim, objs int) (problems.Interface, error) { return Wfg9(), nil }
		}

		problems.DefaultRegistry.Register(name, factory, problems.ProblemMetadata{
			Description: fmt.Sprintf("WFG%d - Walking Fish Group test problem %d", num, num),
			MinDim:      2,
			MaxDim:      1000,
			NumObjs:     0,
			Category:    "many",
		})
	}
}
