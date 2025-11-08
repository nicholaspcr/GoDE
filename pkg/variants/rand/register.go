package rand

import "github.com/nicholaspcr/GoDE/pkg/variants"

func init() {
	variants.DefaultRegistry.Register("rand1", func() variants.Interface {
		return Rand1()
	}, variants.VariantMetadata{
		Description: "Rand/1/Bin - Uses random base vector",
		Category:    "rand",
	})

	variants.DefaultRegistry.Register("rand2", func() variants.Interface {
		return Rand2()
	}, variants.VariantMetadata{
		Description: "Rand/2/Bin - Uses random base with two difference vectors",
		Category:    "rand",
	})
}
