package pbest

import "github.com/nicholaspcr/GoDE/pkg/variants"

func init() {
	variants.DefaultRegistry.Register("pbest", func() variants.Interface {
		return Pbest()
	}, variants.VariantMetadata{
		Description: "pBest/1/Bin - Uses probability-based best selection",
		Category:    "pbest",
	})
}
