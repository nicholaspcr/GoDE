package multi

import "github.com/nicholaspcr/GoDE/pkg/problems"

//nolint:revive // Factory functions have unused parameters matching registry interface
func init() {
	// Register ZDT problems
	problems.DefaultRegistry.Register("zdt1", func(dim, _ int) (problems.Interface, error) {
		return Zdt1(), nil
	}, problems.ProblemMetadata{
		Description: "ZDT1 - Convex Pareto front",
		MinDim:      2,
		MaxDim:      1000,
		NumObjs:     2,
		Category:    "multi",
	})

	problems.DefaultRegistry.Register("zdt2", func(dim, _ int) (problems.Interface, error) {
		return Zdt2(), nil
	}, problems.ProblemMetadata{
		Description: "ZDT2 - Non-convex Pareto front",
		MinDim:      2,
		MaxDim:      1000,
		NumObjs:     2,
		Category:    "multi",
	})

	problems.DefaultRegistry.Register("zdt3", func(dim, _ int) (problems.Interface, error) {
		return Zdt3(), nil
	}, problems.ProblemMetadata{
		Description: "ZDT3 - Disconnected Pareto front",
		MinDim:      2,
		MaxDim:      1000,
		NumObjs:     2,
		Category:    "multi",
	})

	problems.DefaultRegistry.Register("zdt4", func(dim, _ int) (problems.Interface, error) {
		return Zdt4(), nil
	}, problems.ProblemMetadata{
		Description: "ZDT4 - Many local Pareto fronts",
		MinDim:      2,
		MaxDim:      1000,
		NumObjs:     2,
		Category:    "multi",
	})

	problems.DefaultRegistry.Register("zdt6", func(dim, _ int) (problems.Interface, error) {
		return Zdt6(), nil
	}, problems.ProblemMetadata{
		Description: "ZDT6 - Non-uniform Pareto front",
		MinDim:      2,
		MaxDim:      1000,
		NumObjs:     2,
		Category:    "multi",
	})

	// Register VNT problem
	problems.DefaultRegistry.Register("vnt1", func(dim, _ int) (problems.Interface, error) {
		return Vnt1(), nil
	}, problems.ProblemMetadata{
		Description: "VNT1 - Viennet problem",
		MinDim:      2,
		MaxDim:      2,
		NumObjs:     3,
		Category:    "multi",
	})
}
