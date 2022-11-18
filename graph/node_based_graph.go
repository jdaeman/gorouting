package graph

import (
	"fmt"
	"sort"
)

type NodeBasedNode struct {
	first_edge int32
	edges      int32
}

type NodeBasedEdge struct {
	From     int32
	To       int32
	Distance int32

	Forward  bool
	Backward bool
	Reverse  bool

	AnnotationId int32
	// msb is clear => forward
	// msb is set   => backward
	GeometryId uint32
}

type NodeBasedGraph struct {
	number_of_nodes int32
	number_of_edges int32

	nodes []NodeBasedNode
	edges []NodeBasedEdge
}

// split NodeBasedEdge to both direction.
func directedEdges(edges []InternalEdge) []NodeBasedEdge {
	directedEdges := make([]NodeBasedEdge, 0, len(edges))

	for _, edge := range edges {
		if !edge.Forward {
			panic("invalid node_based_edge")
		}

		newEdge := NodeBasedEdge{
			From: edge.From, To: edge.To,
			Forward: edge.Forward, Backward: edge.Backward,
			Distance:     edge.Distance,
			AnnotationId: edge.AnnotationId,
			GeometryId:   edge.GeometryId,
			Reverse:      false,
		}

		directedEdges = append(directedEdges, newEdge)
		// split this edge.
		if !edge.Split {
			// normaly,
			// 1. bidirection way
			// 2. oneway,
			// not, below shape
			//        u  ~> v
			//          ^--/
			newEdge.From, newEdge.To = newEdge.To, newEdge.From
			newEdge.GeometryId |= (1 << 31) // msb is set, backward geometry.
			newEdge.Reverse = !newEdge.Backward
			directedEdges = append(directedEdges, newEdge)
		}
	}

	// sorted by from node, to node.
	sort.Slice(directedEdges, func(l, r int) bool {
		edge1, edge2 := directedEdges[l], directedEdges[r]
		if edge1.From == edge2.From {
			return edge1.To < edge2.To
		}
		return edge1.From < edge2.From
	})

	return directedEdges
}

func NewNodeBasedGraph(nodes int32, internalEdges []InternalEdge) *NodeBasedGraph {
	if nodes == 0 || len(internalEdges) == 0 {
		return nil
	}
	edges := directedEdges(internalEdges)

	// edges must be sorted.
	number_of_edges := int32(len(edges))
	node_array := make([]NodeBasedNode, nodes)

	edge, position := int32(0), int32(0)
	for node := int32(0); node < nodes; node++ {
		last_edge := edge
		for edge < number_of_edges && edges[edge].From == node {
			edge++
		}

		node_array[node].edges = edge - last_edge
		node_array[node].first_edge = position
		position += node_array[node].edges
	}

	nbg := &NodeBasedGraph{
		number_of_nodes: nodes, number_of_edges: number_of_edges,
		nodes: node_array, edges: edges}

	return nbg
}

func (graph *NodeBasedGraph) GetNumberOfNodes() int32 {
	return graph.number_of_nodes
}

func (graph *NodeBasedGraph) GetNumberOfEdges() int32 {
	return graph.number_of_edges
}

func (graph *NodeBasedGraph) GetOutDegree(node int32) int32 {
	return graph.nodes[node].edges
}

func (graph *NodeBasedGraph) BeginEdges(node int32) int32 {
	return graph.nodes[node].first_edge
}

func (graph *NodeBasedGraph) EndEdges(node int32) int32 {
	return graph.nodes[node].first_edge + graph.nodes[node].edges
}

func (graph *NodeBasedGraph) GetEdgePtr(edge_id int32) *NodeBasedEdge {
	return &graph.edges[edge_id]
}

func (graph *NodeBasedGraph) GetEdgeData(edge_id int32) NodeBasedEdge {
	return *graph.GetEdgePtr(edge_id)
}

func (graph *NodeBasedGraph) GetTarget(edge_id int32) int32 {
	return graph.edges[edge_id].To
}

func (graph *NodeBasedGraph) GetSource(edge_id int32) int32 {
	return graph.edges[edge_id].From
}

func (graph *NodeBasedGraph) FindEdge(u, v int32) int32 {
	first := graph.nodes[u].first_edge
	last := first + graph.nodes[u].edges

	for ; first < last; first++ {
		if graph.edges[first].To == v {
			return first
		}
	}
	panic(fmt.Sprintf("Could not find edge %d -> %v", u, v))
}

// func (graph *NodeBasedGraph) FindEdgeEitherDirection(u, v int32) int32 {
// 	edge1 := graph.FindEdge(u, v)
// 	if edge1 == -1 {
// 		return graph.FindEdge(v, u)
// 	}
// 	return edge1
// }

// func (graph *NodeBasedGraph) DeleteEdge(u int32, edge_id int32) {
// 	node := &graph.nodes[u]

// 	graph.number_of_edges-- // decrease total edge count
// 	node.edges--            // decrease node out degree

// 	if node.edges < 0 {
// 		panic("do not delete edge")
// 	}

// 	last := node.first_edge + node.edges
// 	graph.edges[edge_id] = graph.edges[last]
// }

// func (graph *NodeBasedGraph) SetTarget(edge_id int32, v int32) {
// 	graph.edges[edge_id].To = v
// }
