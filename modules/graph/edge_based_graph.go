package graph

import (
	"sort"
)

// GenerateEdgeExpandedNodes()

// GenerateEdgeExpandedEdges()

//

type EdgeBasedNode struct {
	GeometryId   uint32
	AnnotationId int32
}

type EdgeBasedNodeSegment struct {
	Forward_id  int32 // EdgeBasedNode id
	Backward_id int32 // EdgeBasedNode id

	U, V int32 // node id, for location
	Pos  int32 // seg pos in this node
}

type EdgeBasedEdge struct {
	Source int32 // EdgeBasedNode id
	Target int32 // EdgeNasedNode id

	Distance int32

	Forward  bool
	Backward bool
}

type EdgeBasedGraph struct {
	number_of_nodes int32
	number_of_edges int32

	edges []EdgeBasedEdge

	// forward edge list
	// [node][0]: forward first edge id
	// [node][1]: edge count
	fwd_adjs [][2]int32

	// reverse edge list
	// [node][0]: reverse first edge id
	// [node][1]: edge count
	rev_adjs [][2]int32
}

func NewEdgeBasedGraph(nodes []EdgeBasedNode, edges []EdgeBasedEdge) *EdgeBasedGraph {
	if nodes == nil || edges == nil {
		return nil
	}

	maxNode := int32(len(nodes))

	revEdges := make([]EdgeBasedEdge, len(edges))
	if copy(revEdges, edges) == len(edges) {
		// copy success
		// swap source, target
		// swap forward, backward
		// for reverse edge list.
		for i := range revEdges {
			edge := &revEdges[i]
			edge.Source, edge.Target = edge.Target, edge.Source
			edge.Forward, edge.Backward = edge.Backward, edge.Forward
		}
	} else {
		panic("Copy fail")
	}

	// sorted by source, target node id
	edgeSort := func(dirEdges *[]EdgeBasedEdge) {
		sort.Slice(*dirEdges, func(l, r int) bool {
			left, right := (*dirEdges)[l], (*dirEdges)[r]
			if left.Source == right.Source {
				return left.Target < right.Target
			}
			return left.Source < right.Source
		})
	}

	makeGraph := func(dirEdges []EdgeBasedEdge, node_array [][2]int32, forward bool) {
		edge, position, number_of_edges := int32(0), int32(0), int32(len(dirEdges))
		base := int32(0)
		if !forward {
			base = number_of_edges
		}

		for node := int32(0); node < maxNode; node++ {
			last_edge := edge
			for edge < number_of_edges && dirEdges[edge].Source == node {
				edge++
			}

			node_array[node][0] = position + base  // first edge id, with bias value
			node_array[node][1] = edge - last_edge // edge count
			position += node_array[node][1]        // move first edge id for next node id
		}
	}

	edgeSort(&edges)
	edgeSort(&revEdges)
	//            <-forward-> | <-   reverse  ->
	// totalEdge: [0][1]...[N][N+1][N+2]...[2N-1]
	totalEdge := append(edges, revEdges...)
	if len(totalEdge) < len(edges)*2 {
		panic("Edge is incomplete")
	}

	fwd_adjs := make([][2]int32, maxNode)
	makeGraph(edges, fwd_adjs, true)

	rev_adjs := make([][2]int32, maxNode)
	makeGraph(revEdges, rev_adjs, false)

	edgeBasedGraph := &EdgeBasedGraph{
		number_of_nodes: maxNode,
		number_of_edges: int32(len(totalEdge)),
		edges:           totalEdge,
		fwd_adjs:        fwd_adjs,
		rev_adjs:        rev_adjs,
	}

	return edgeBasedGraph
}

func (g *EdgeBasedGraph) getEdges(u int32, forward bool) []int32 {
	edgeRange := g.getEdgeRange(u, forward)
	start, count := edgeRange[0], edgeRange[1]

	edges := make([]int32, 0, count)
	for edge := start; edge < count; edge++ {
		edges = append(edges, edge)
	}
	return edges
}

func (g *EdgeBasedGraph) getEdgeRange(u int32, forward bool) [2]int32 {
	start, count := g.fwd_adjs[u][0], g.fwd_adjs[u][1]
	if !forward {
		start, count = g.rev_adjs[u][0], g.rev_adjs[u][1]
	}
	return [2]int32{start, start + count}
}

func (g *EdgeBasedGraph) GetForwardEdgeRange(u int32) [2]int32 {
	return g.getEdgeRange(u, true)
}

func (g *EdgeBasedGraph) GetBackwardEdgeRange(u int32) [2]int32 {
	return g.getEdgeRange(u, false)
}

func (g *EdgeBasedGraph) GetForwardEdges(u int32) []int32 {
	return g.getEdges(u, true)
}

func (g *EdgeBasedGraph) GetBackwardEdges(u int32) []int32 {
	return g.getEdges(u, false)
}

func (g *EdgeBasedGraph) GetNumberOfEdges() int32 {
	return g.number_of_edges
}

func (g *EdgeBasedGraph) GetEdgeData(edge_id int32) EdgeBasedEdge {
	return g.edges[edge_id]
}

func (g *EdgeBasedGraph) GetNumberOfNodes() int32 {
	return g.number_of_nodes
}
