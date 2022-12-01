package engine

import (
	"container/heap"
	"fmt"
)

// fwdHeap pop
// check at revHeap
// newWeight = revHeapNode->Weight + heapNode.weight
// newWeight < upper_bound
// |-> change middle node, upper bound
//
// heapNode.weight + min_edge_offset > upper_bound => delete all -> return
// => min_edge_offset: min(0, fwdheap.minKey()) # 0?
// relax outgoing edges

func (engine *RoutingEngine) relaxOutgoingEdges(fwdHeap *MinHeap, parents []int32, fwdTable []int32, revTable []int32, upper_bound int32, middle_node int32, forward bool) (int32, int32) {

	curNode := heap.Pop(fwdHeap).(HeapNode)
	u, weight := curNode.u, curNode.weight

	if revTable[u] != MAX_WEIGHT {
		newWeight := revTable[u] + curNode.weight

		if newWeight < upper_bound {
			upper_bound = newWeight
			middle_node = curNode.u
		}
	}

	if curNode.weight > upper_bound {
		// delete all elem from heap
		for fwdHeap.Len() > 0 {
			heap.Pop(fwdHeap)
		}
		return upper_bound, middle_node
	}

	ebg := engine.edge_based_graph
	adjEdges := ebg.GetForwardEdges(u)
	if !forward {
		adjEdges = ebg.GetBackwardEdges(u)
	}

	for _, edgeId := range adjEdges {
		edge := ebg.GetEdgeData(edgeId)
		v, nextWeight := edge.Target, edge.Distance

		if weight+nextWeight < fwdTable[v] {
			fwdTable[v] = weight + nextWeight
			heap.Push(fwdHeap, HeapNode{u: v, weight: fwdTable[v]})
			parents[v] = u
		}
	}

	return upper_bound, middle_node
}

// bidirection dijkstra
// func (engine *RoutingEngine) ShortestPathSearch(s, g int32) {
// 	ebg := engine.edge_based_graph
// 	maxNode := ebg.GetNumberOfNodes()

// 	fwdHeap, revHeap := &MinHeap{}, &MinHeap{}
// 	heap.Init(fwdHeap)
// 	heap.Init(revHeap)

// 	fwdCosts, fwdParents := make([]int32, maxNode), make([]int32, maxNode)
// 	revCosts, revParents := make([]int32, maxNode), make([]int32, maxNode)

// 	upper_bound, middle_node := MAX_WEIGHT, int32(-1)

// 	for i := range fwdCosts {
// 		fwdCosts[i], fwdParents[i] = MAX_WEIGHT, -1
// 		revCosts[i], revParents[i] = MAX_WEIGHT, -1
// 	}

// 	fwdCosts[s], fwdParents[s] = 0, s
// 	revCosts[g], revParents[g] = 0, g

// 	heap.Push(fwdHeap, HeapNode{u: s, weight: 0})
// 	heap.Push(revHeap, HeapNode{u: g, weight: 0})

// 	for fwdHeap.Len() > 0 || revHeap.Len() > 0 {

// 		// forward search
// 		if fwdHeap.Len() > 0 {
// 			upper_bound, middle_node =
// 				engine.relaxOutgoingEdges(
// 					fwdHeap, fwdParents,
// 					fwdCosts, revCosts,
// 					upper_bound, middle_node,
// 					true)
// 		}

// 		// reverse search
// 		if revHeap.Len() > 0 {
// 			upper_bound, middle_node =
// 				engine.relaxOutgoingEdges(
// 					revHeap, revParents,
// 					revCosts, fwdCosts,
// 					upper_bound, middle_node,
// 					false)
// 		}
// 	}

// 	path1 := makePathFwd(fwdParents, middle_node)
// 	path2 := makePathRev(revParents, middle_node)
// 	fmt.Println(upper_bound, middle_node)
// 	fmt.Println(path1)
// 	fmt.Println(path2)
// }

func (engine *RoutingEngine) SearchRoute(coords [2][2]float64) {
	phantoms := engine.rtree.GetPhantomNodes([][2]float64{coords[0], coords[1]})

	if len(phantoms) != len(coords) {
		fmt.Println("Fail map matching")
	}

}

// bidirection dijkstra
func (engine *RoutingEngine) ShortestPathSearch(sources, goals []int32) []int32 {
	ebg := engine.edge_based_graph
	maxNode := ebg.GetNumberOfNodes()

	fwdHeap, revHeap := &MinHeap{}, &MinHeap{}
	heap.Init(fwdHeap)
	heap.Init(revHeap)

	fwdCosts, fwdParents := make([]int32, maxNode), make([]int32, maxNode)
	revCosts, revParents := make([]int32, maxNode), make([]int32, maxNode)

	upper_bound, middle_node := MAX_WEIGHT, int32(-1)

	for i := range fwdCosts {
		fwdCosts[i], fwdParents[i] = MAX_WEIGHT, -1
		revCosts[i], revParents[i] = MAX_WEIGHT, -1
	}

	for _, source := range sources {
		fwdCosts[source], fwdParents[source] = 0, source
		heap.Push(fwdHeap, HeapNode{u: source, weight: 0})
	}

	for _, goal := range goals {
		revCosts[goal], revParents[goal] = 0, goal
		heap.Push(revHeap, HeapNode{u: goal, weight: 0})
	}

	for fwdHeap.Len() > 0 || revHeap.Len() > 0 {

		// forward search
		if fwdHeap.Len() > 0 {
			upper_bound, middle_node =
				engine.relaxOutgoingEdges(
					fwdHeap, fwdParents,
					fwdCosts, revCosts,
					upper_bound, middle_node,
					true)
		}

		// reverse search
		if revHeap.Len() > 0 {
			upper_bound, middle_node =
				engine.relaxOutgoingEdges(
					revHeap, revParents,
					revCosts, fwdCosts,
					upper_bound, middle_node,
					false)
		}
	}

	if middle_node == -1 || upper_bound == MAX_WEIGHT {
		return nil
	}

	path1 := makePathFwd(fwdParents, middle_node)
	path2 := makePathRev(revParents, middle_node)

	path1 = append(path1, path2[1:]...)

	return path1
}
