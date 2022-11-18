package extract

import (
	"graph"
	"math"
	"sort"
)

type EdgeBasedGraphFactory struct {
	nbg_factory *NodeBasedGraphFactory

	number_of_edge_based_nodes int32
	nbe_to_ebn_mapping         []int32

	edge_based_node_segments []graph.EdgeBasedNodeSegment
	edge_based_nodes         []graph.EdgeBasedNode
	edge_based_edges         []graph.EdgeBasedEdge
}

func NewEdgeBasedGraphFactory(nbgFactory *NodeBasedGraphFactory) *EdgeBasedGraphFactory {
	ret := &EdgeBasedGraphFactory{nbg_factory: nbgFactory}
	return ret
}

// mapping.
// node_based_edge id -> edge_based_node id
func (factory *EdgeBasedGraphFactory) LabelEdgeBasedNodes() int32 {
	nbg := factory.nbg_factory.GetGraph()

	mapping := make([]int32, nbg.GetNumberOfEdges())
	for i := range mapping {
		// math.MaxInt32 is invalid id.
		mapping[i] = math.MaxInt32
	}

	edge_based_node_id := int32(0)
	for current_node := int32(0); current_node < nbg.GetNumberOfNodes(); current_node++ {
		current_edge := nbg.BeginEdges(current_node)

		for ; current_edge < nbg.EndEdges(current_node); current_edge++ {
			if nbg.GetEdgeData(current_edge).Reverse {
				continue
			}

			mapping[current_edge] = edge_based_node_id
			edge_based_node_id++
		}
	}

	factory.nbe_to_ebn_mapping = mapping
	return edge_based_node_id
}

// TBD
// EdgeBasedNodeSegment for RTree (map matching)
func (factory *EdgeBasedGraphFactory) InsertEdgeBasedNode(u, v int32) {
	nbg_factory := factory.nbg_factory
	nbg := nbg_factory.GetGraph()
	mapping := factory.nbe_to_ebn_mapping
	segments := &factory.edge_based_node_segments

	edge1 := nbg.FindEdge(u, v)
	edge2 := nbg.FindEdge(v, u)

	if mapping[edge1] == math.MaxInt32 {
		panic("Edge is reverse?")
	}

	geometryId := nbg.GetEdgeData(edge1).GeometryId
	nodes := nbg_factory.GetGeometry(int32(geometryId))
	segCount := len(nodes)

	for i := 1; i < segCount; i++ {
		coord1, coord2 := nodes[i-1], nodes[i]

		*segments = append(*segments, graph.EdgeBasedNodeSegment{
			Forward_id:  mapping[edge1],
			Backward_id: mapping[edge2],
			U:           coord1,
			V:           coord2,
		})
	}
}

// TBD
func (factory *EdgeBasedGraphFactory) GenerateEdgeExpandedNodes() {
	nbg := factory.nbg_factory.GetGraph()
	mapping := factory.nbe_to_ebn_mapping
	segments := &factory.edge_based_node_segments
	*segments = make([]graph.EdgeBasedNodeSegment, 0)

	for u := int32(0); u < nbg.GetNumberOfNodes(); u++ {
		nbg_edge_id := nbg.BeginEdges(u)

		for ; nbg_edge_id < nbg.EndEdges(u); nbg_edge_id++ {
			v := nbg.GetTarget(nbg_edge_id)

			if u >= v {
				continue
			}

			// insert edge based node
			if mapping[nbg_edge_id] == math.MaxInt32 {
				// v -> u
				factory.InsertEdgeBasedNode(v, u)
			} else {
				// u -> v
				factory.InsertEdgeBasedNode(u, v)
			}
		}
	}

	// sort segments
	// TBD. via-way turn
}

//
func (factory *EdgeBasedGraphFactory) GenerateEdgeExpandedEdges() {
	nbg_factory := factory.nbg_factory
	nbg := nbg_factory.GetGraph()
	mapping := factory.nbe_to_ebn_mapping

	//lineEdges := &factory.edge_based_edges
	lineEdges := make([]graph.EdgeBasedEdge, 0)

	generate_edge := func(
		edge_based_node_from,
		edge_based_node_to int32,
		intersection int32,
		distance int32) *graph.EdgeBasedEdge {

		return &graph.EdgeBasedEdge{
			Source:   edge_based_node_from,
			Target:   edge_based_node_to,
			Distance: distance,
			Forward:  true,
			Backward: false,
		}
	}

	nodeCount := nbg.GetNumberOfNodes()
	for via := int32(0); via < nodeCount; via++ {
		incomingEdges := nbg_factory.GetIncomingEdges(via)
		outgoingEdges := nbg_factory.GetOutgoingEdges(via)

		for _, incoming := range incomingEdges {
			fromNode := nbg.GetSource(incoming)
			distance := nbg.GetEdgeData(incoming).Distance

			for _, outgoing := range outgoingEdges {
				toNode := nbg.GetTarget(outgoing)

				// is turn allowed => continue
				if !nbg_factory.IsTurnallowed(fromNode, via, toNode) {
					continue
				}

				edge := generate_edge(mapping[incoming], mapping[outgoing], via, distance)
				lineEdges = append(lineEdges, *edge)
			}
		}
	}

	// sorted by source, target edge_based_node id
	sort.Slice(lineEdges, func(l, r int) bool {
		left, right := lineEdges[l], lineEdges[r]
		if left.Source == right.Source {
			return left.Target < right.Target
		}
		return left.Source < right.Source
	})

	// finish
	factory.edge_based_edges = lineEdges
}

func (factory *EdgeBasedGraphFactory) Run() {
	number_of_edge_based_node := factory.LabelEdgeBasedNodes()
	factory.number_of_edge_based_nodes = number_of_edge_based_node

	factory.GenerateEdgeExpandedNodes()
	factory.GenerateEdgeExpandedEdges()
}

// func getIncomingEdges(nbg *graph.NodeBasedGraph, via int32) []int32 {
// 	curEdge, lastEdge := nbg.BeginEdges(via), nbg.EndEdges(via)
// 	edges := make([]int32, 0)

// 	for ; curEdge < lastEdge; curEdge++ {
// 		from := nbg.GetTarget(curEdge)
// 		incoming := nbg.FindEdge(from, via)

// 		if !nbg.GetEdgeData(incoming).Reverse {
// 			edges = append(edges, incoming)
// 		}
// 	}
// 	return edges
// }

// func getOutgoingEdges(nbg *graph.NodeBasedGraph, via int32) []int32 {
// 	curEdge, lastEdge := nbg.BeginEdges(via), nbg.EndEdges(via)
// 	edges := make([]int32, 0)

// 	for ; curEdge < lastEdge; curEdge++ {
// 		edges = append(edges, curEdge)
// 	}

// 	return edges
// }

// func isTurnAllowed(restrictions []graph.InternalRestriction, from, via, to int32) bool {
// 	count := len(restrictions)
// 	i := sort.Search(count, func(i int) bool {
// 		return restrictions[i].From >= from
// 	})
// 	j := sort.Search(count, func(i int) bool {
// 		return restrictions[i].From > from
// 	})

// 	for ; i < j; i++ {
// 		restriction := &restrictions[i]

// 		if restriction.Via != via {
// 			continue
// 		}

// 		if restriction.Only {
// 			return restriction.To == to
// 		} else {
// 			if restriction.To == to {
// 				return false
// 			}
// 		}
// 	}

// 	return true
// }
