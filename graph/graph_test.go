package graph_test

import (
	"files"
	"fmt"
	"graph"
	"testing"
)

func TestOsmGraph(t *testing.T) {
	geos := files.LoadGeoNodes("../extract/data/map.geos")
	uncomps := files.LoadUncompEdges("../extract/data/map.uncomp")

	nodeCount := len(geos)
	graph := graph.NewOsmGraph(int32(nodeCount), uncomps)

	var v int32
	for v = 0; v < int32(nodeCount); v++ {
		degree := graph.GetOutDegree(v)
		if degree == 2 {
			break
		}
	}

	u := graph.GetTarget(v, 0)
	w := graph.GetTarget(v, 1)

	fwd1, _ := graph.FindConstEdge(u, v)
	fwd2, _ := graph.FindConstEdge(v, w)
	rev1, _ := graph.FindConstEdge(w, v)
	rev2, _ := graph.FindConstEdge(v, u)

	graph.DelEdge(v, u)
	graph.DelEdge(v, w)

	graph.SetNewTarget(u, v, w)
	graph.SetNewTarget(w, v, u)

	fmt.Println(graph.FindEdge(u, w).Distance)
	fmt.Println(graph.FindEdge(w, u).Distance)

	graph.FindEdge(u, w).Distance = fwd1.Distance + fwd2.Distance
	graph.FindEdge(w, u).Distance = rev1.Distance + rev2.Distance

	if graph.GetOutDegree(v) != 0 {
		t.Fail()
	}

	fmt.Println(graph.FindEdge(u, w).Distance)
	fmt.Println(graph.FindEdge(w, u).Distance)

}
