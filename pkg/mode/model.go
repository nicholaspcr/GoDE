package mo

// Params of the moDE
type Params struct {
	EXECS, NP, M, DIM, GEN int
	FLOOR, CEIL, CR, F, P  float64
	MemProf, CPUProf       string
}

// Elem -> Element of population
type Elem struct {
	X      []float64
	Objs   []float64
	Crwdst float64
}

// Copy the entire struct
func (e *Elem) Copy() Elem {
	var ret Elem
	ret.X = make([]float64, len(e.X))
	ret.Objs = make([]float64, len(e.Objs))
	copy(ret.X, e.X)
	copy(ret.Objs, e.Objs)
	ret.Crwdst = e.Crwdst
	return ret
}

func (e *Elem) dominates(other Elem) bool {
	if len(e.Objs) != len(other.Objs) {
		return false
	}
	dominates := false
	for i := range e.Objs {
		if e.Objs[i] > other.Objs[i] {
			return false
		}
		if e.Objs[i] < other.Objs[i] {
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
