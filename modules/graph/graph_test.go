package graph_test

import (
	"fmt"
	"log"
	"modules/files"
	"modules/graph"
	"modules/view"
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

	dataReader := files.NewReader("../extract/data/map.osm")

	nodes := dataReader.LoadGeoNodes()            //files.LoadGeoNodes("../extract/data/map.node")
	edges := dataReader.LoadEdges()               //files.LoadEdges("../extract/data/map.edge")
	annos := dataReader.LoadEdgeAnnotations()     //files.LoadEdgeAnnotations("../extract/data/map.anno")
	restrictions := dataReader.LoadRestrictions() //files.LoadRestrictions("../extract/data/map.restriction")

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

func dfs(u int32, parent int32, graph *graph.EdgeBasedGraph, visit []bool, path *[]int32) {
	*path = append(*path, u)
	visit[u] = true

	edgeRange := graph.GetForwardEdgeRange(u)
	for edgeId := edgeRange[0]; edgeId < edgeRange[1]; edgeId++ {
		v := graph.GetEdgeData(edgeId).Target

		if v != parent && visit[v] == false {
			dfs(v, u, graph, visit, path)
			break
		}
	}
}

func TestEBG(t *testing.T) {
	dataReader := files.NewReader("../extract/data/map.osm")
	nodes := dataReader.LoadEdgeBasedNodes()
	edges := dataReader.LoadEdgeBasedEdges()

	log.Println("node count", len(nodes))
	log.Println("edge count", len(edges))

	graph := graph.NewEdgeBasedGraph(nodes, edges)
	log.Println(graph.GetNumberOfEdges())

	for u := int32(0); u < 4; u++ {
		t := graph.GetForwardEdgeRange(u)
		log.Println("forward", t[0], t[1])
		t = graph.GetBackwardEdgeRange(u)
		log.Println("backward", t[0], t[1])
	}

	path := make([]int32, 0, 100)
	visit := make([]bool, len(nodes))
	dfs(0, 0, graph, visit, &path)

	viewFactory := view.NewViewFactory("../extract/data/map.osm")
	viewFactory.DrawingEdges(path, -1)
}
