package best

import "github.com/nicholaspcr/GoDE/pkg/variants"

func init() {
	variants.DefaultRegistry.Register("best/1/bin", func() variants.Interface {
		return Best1()
	}, variants.VariantMetadata{
		Description: "Best/1/Bin - Uses best vector from rank zero",
		Category:    "best",
	})

	variants.DefaultRegistry.Register("best/2/bin", func() variants.Interface {
		return Best2()
	}, variants.VariantMetadata{
		Description: "Best/2/Bin - Uses best vector with two difference vectors",
		Category:    "best",
	})
}
