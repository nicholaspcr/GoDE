package mo

// Params of the moDE
type Params struct {
	EXECS, NP, M, DIM, GEN int
	FLOOR, CEIL, CR, F, P  float64
}

// Elem -> Element of population
type Elem struct {
	X      []float64
	objs   []float64
	crwdst float64
}

// Copy the entire struct
func (e *Elem) Copy() Elem {
	var ret Elem
	ret.X = make([]float64, len(e.X))
	ret.objs = make([]float64, len(e.objs))
	copy(ret.X, e.X)
	copy(ret.objs, e.objs)
	ret.crwdst = e.crwdst
	return ret
}

func (e *Elem) dominates(other Elem) bool {
	if len(e.objs) != len(other.objs) {
		return false
	}
	dominates := false
	for i := range e.objs {
		if e.objs[i] > other.objs[i] {
			return false
		}
		if e.objs[i] < other.objs[i] {
			dominates = true
		}
	}
	return dominates
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
