package engine

import (
	"container/heap"
)

func makePathFwd(parents []int32, u int32) []int32 {
	ret := []int32{u}
	if u == parents[u] {
		// arrive source
		return ret
	}

	nextPath := makePathFwd(parents, parents[u])
	ret = append(nextPath, ret...)
	return ret
}

func makePathRev(parents []int32, u int32) []int32 {
	ret := []int32{u}
	if u == parents[u] {
		// arrive goal
		return ret
	}

	nextPath := makePathRev(parents, parents[u])
	ret = append(ret, nextPath...)
	return ret
}

func (engine *RoutingEngine) SimplePathSearch(s, g int32) []int32 {
	ebg := engine.edge_based_graph
	maxNode := ebg.GetNumberOfNodes()
	fwdHeap := &MinHeap{}

	parents := make([]int32, maxNode)
	costTable := make([]int32, maxNode)
	for i := range costTable {
		parents[i] = -1
		costTable[i] = MAX_WEIGHT
	}

	heap.Init(fwdHeap)
	heap.Push(fwdHeap, HeapNode{u: s, weight: 0})
	costTable[s], parents[s] = 0, s

	for fwdHeap.Len() > 0 {
		curNode := heap.Pop(fwdHeap).(HeapNode)
		if curNode.u == g {
			break
		}

		u, weight := curNode.u, curNode.weight
		edgeRange := ebg.GetForwardEdgeRange(u)

		for edgeId := edgeRange[0]; edgeId < edgeRange[1]; edgeId++ {
			edge := ebg.GetEdgeData(edgeId)
			v, nextWeight := edge.Target, edge.Distance

			if weight+nextWeight < costTable[v] {
				costTable[v] = weight + nextWeight
				parents[v] = u
				heap.Push(fwdHeap, HeapNode{u: v, weight: costTable[v]})
			}
		}
	}

	ret := makePathFwd(parents, g)
	return ret
}
