package graph_test

import (
	"files"
	"fmt"
	"graph"
	"sort"
	"testing"
)

func getIncomingEdges(nbg *graph.NodeBasedGraph, via int32) []int32 {
	curEdge, lastEdge := nbg.BeginEdges(via), nbg.EndEdges(via)
	edges := make([]int32, 0)

	for ; curEdge < lastEdge; curEdge++ {
		from := nbg.GetTarget(curEdge)
		incoming := nbg.FindEdge(from, via)

		if !nbg.GetEdgeData(incoming).Reverse {
			edges = append(edges, incoming)
		}
	}
	return edges
}

func getOutgoingEdges(nbg *graph.NodeBasedGraph, via int32) []int32 {
	curEdge, lastEdge := nbg.BeginEdges(via), nbg.EndEdges(via)
	edges := make([]int32, 0)

	for ; curEdge < lastEdge; curEdge++ {
		edges = append(edges, curEdge)
	}

	return edges
}

func isTurnAllowed(restrictions []graph.InternalRestriction, from, via, to int32) bool {
	count := len(restrictions)
	i := sort.Search(count, func(i int) bool {
		return restrictions[i].From >= from
	})
	j := sort.Search(count, func(i int) bool {
		return restrictions[i].From > from
	})

	for ; i < j; i++ {
		restriction := &restrictions[i]

		if restriction.Via != via {
			continue
		}

		if restriction.Only {
			return restriction.To == to
		} else {
			if restriction.To == to {
				return false
			}
		}
	}

	return true
}

func TestNBG(t *testing.T) {
	nodes := files.LoadGeoNodes("../extract/data/map.node")
	edges := files.LoadEdges("../extract/data/map.edge")
	annos := files.LoadEdgeAnnotations("../extract/data/map.anno")
	restrictions := files.LoadRestrictions("../extract/data/map.restriction")

	nodeCount := int32(len(nodes))
	nbg := graph.NewNodeBasedGraph(nodeCount, edges)

	fmt.Println(len(restrictions))

	var via int32
	for via = 0; via < nodeCount; via++ {

		incomingEdges := getIncomingEdges(nbg, via)
		outgoingEdges := getOutgoingEdges(nbg, via)

		for _, incoming_edge := range incomingEdges {
			from := nbg.GetSource(incoming_edge)
			// is incoming  edge restriction via?

			for _, outgoing_edge := range outgoingEdges {
				to := nbg.GetTarget(outgoing_edge)
				// is turn allowed?
				inway := annos[nbg.GetEdgeData(incoming_edge).AnnotationId].Id
				outway := annos[nbg.GetEdgeData(outgoing_edge).AnnotationId].Id
				if !isTurnAllowed(restrictions, from, via, to) {
					fmt.Println(inway, "->", outway, "turn restriction")
				}
			}
		}
	}
}
