package models

type Mutate func(elems, rankZero Population, p VariantParams) (Vector, error)

// Definition definition of the test case functions
type Variant struct {
	Fn          func(elems, rankZero Population, p VariantParams) (Vector, error)
	VariantName string
}

func (v *Variant) Name() string {
	return v.VariantName
}

func (v *Variant) Mutate(
	elems, rankZero Population,
	p VariantParams,
) (Vector, error) {

	return v.Fn(elems, rankZero, p)
}
