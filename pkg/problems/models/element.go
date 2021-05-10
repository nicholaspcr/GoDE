package models

// Elem -> Element of population
type Elem struct {
	X      []float64
	Objs   []float64
	Crwdst float64
}

// Copy the entire struct
func (e *Elem) Copy() Elem {
	var ret = Elem{}

	ret.X = make([]float64, len(e.X))
	copy(ret.X, e.X)

	ret.Objs = make([]float64, len(e.Objs))
	copy(ret.Objs, e.Objs)

	return ret
}

// Elements is a slice of the type Elem
type Elements []Elem

// Copy of the []Elem slice
func (e Elements) Copy() Elements {
	arr := make(Elements, len(e))
	for i, v := range e {
		arr[i] = v.Copy()
	}
	return arr
}
