package graph

import "math"

type Geometry struct {
	Nodes     []int32
	Distances []int32
}

func NewGeometry(count int) *Geometry {
	ret := &Geometry{}

	ret.Nodes = make([]int32, 0, count)
	ret.Distances = make([]int32, 1, count)
	ret.Distances[0] = math.MaxInt32

	return ret
}
