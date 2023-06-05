package api

import "google.golang.org/protobuf/proto"

// Vectors is a slice of Vector.
type Vectors []*Vector

// Copy returns a copy of the Vectors.
func (v *Vectors) Copy() Vectors {
	vectors := make(Vectors, len(*v))
	for i, vec := range *v {
		vec := vec
		vectors[i] = proto.Clone(vec).(*Vector)
	}
	return vectors
}

// Copy returns a copy of the Vector.
func (v *Vector) Copy() *Vector {
	return proto.Clone(v).(*Vector)
}

// Copy returns a copy of the Population.
func (v *Population) Copy() *Population {
	return proto.Clone(v).(*Population)
}
