package dtlz

import "github.com/nicholaspcr/GoDE/pkg/problems"

func init() {
	// Register DTLZ problems
	problems.DefaultRegistry.Register("dtlz1", func(dim, objs int) (problems.Interface, error) {
		return Dtlz1(), nil
	}, problems.ProblemMetadata{
		Description: "DTLZ1 - Linear Pareto front with (11^k - 1) local fronts",
		MinDim:      3,
		MaxDim:      1000,
		NumObjs:     0, // Variable
		Category:    "many",
	})

	problems.DefaultRegistry.Register("dtlz2", func(dim, objs int) (problems.Interface, error) {
		return Dtlz2(), nil
	}, problems.ProblemMetadata{
		Description: "DTLZ2 - Concave Pareto front",
		MinDim:      3,
		MaxDim:      1000,
		NumObjs:     0,
		Category:    "many",
	})

	problems.DefaultRegistry.Register("dtlz3", func(dim, objs int) (problems.Interface, error) {
		return Dtlz3(), nil
	}, problems.ProblemMetadata{
		Description: "DTLZ3 - Concave front with (3^k - 1) local fronts",
		MinDim:      3,
		MaxDim:      1000,
		NumObjs:     0,
		Category:    "many",
	})

	problems.DefaultRegistry.Register("dtlz4", func(dim, objs int) (problems.Interface, error) {
		return Dtlz4(), nil
	}, problems.ProblemMetadata{
		Description: "DTLZ4 - Concave front with biased density",
		MinDim:      3,
		MaxDim:      1000,
		NumObjs:     0,
		Category:    "many",
	})

	problems.DefaultRegistry.Register("dtlz5", func(dim, objs int) (problems.Interface, error) {
		return Dtlz5(), nil
	}, problems.ProblemMetadata{
		Description: "DTLZ5 - Degenerate Pareto front",
		MinDim:      3,
		MaxDim:      1000,
		NumObjs:     0,
		Category:    "many",
	})

	problems.DefaultRegistry.Register("dtlz6", func(dim, objs int) (problems.Interface, error) {
		return Dtlz6(), nil
	}, problems.ProblemMetadata{
		Description: "DTLZ6 - Degenerate front with (3^k - 1) local fronts",
		MinDim:      3,
		MaxDim:      1000,
		NumObjs:     0,
		Category:    "many",
	})

	problems.DefaultRegistry.Register("dtlz7", func(dim, objs int) (problems.Interface, error) {
		return Dtlz7(), nil
	}, problems.ProblemMetadata{
		Description: "DTLZ7 - Disconnected Pareto regions",
		MinDim:      3,
		MaxDim:      1000,
		NumObjs:     0,
		Category:    "many",
	})
}
