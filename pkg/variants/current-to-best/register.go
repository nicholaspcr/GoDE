package currenttobest

import "github.com/nicholaspcr/GoDE/pkg/variants"

func init() {
	variants.DefaultRegistry.Register("current-to-best/1/bin", func() variants.Interface {
		return CurrToBest1()
	}, variants.VariantMetadata{
		Description: "Current-to-best/1/Bin - Combination of current and best vectors",
		Category:    "current-to-best",
	})
}
