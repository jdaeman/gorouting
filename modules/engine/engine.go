package engine

import (
	"fmt"
	"log"
	"modules/files"
	"modules/graph"
)

const (
	MAX_WEIGHT = int32(987654321)
)

type EngineConfig struct {
	dataPath string
}

func NewEngineConfig(datapath string) *EngineConfig {
	ret := &EngineConfig{
		dataPath: datapath,
	}

	return ret
}

type RoutingEngine struct {
	rtree            *RTree
	edge_based_graph *graph.EdgeBasedGraph
}

func NewRoutingEngine(config EngineConfig) *RoutingEngine {
	dataReader := files.NewReader(config.dataPath)

	edge_based_nodes := dataReader.LoadEdgeBasedNodes()
	edge_based_edges := dataReader.LoadEdgeBasedEdges()

	fmt.Println("Node / Edge", len(edge_based_nodes), len(edge_based_edges))
	ebg := graph.NewEdgeBasedGraph(edge_based_nodes, edge_based_edges)

	ret := &RoutingEngine{
		edge_based_graph: ebg,
	}

	return ret
}

func NewRoutingEngineByData(
	nodes []graph.EdgeBasedNode,
	edges []graph.EdgeBasedEdge,
	geoNodes []graph.ResultNode,
	segs []graph.EdgeBasedNodeSegment) *RoutingEngine {

	log.Println("Make edge_based_graph...")
	ebg := graph.NewEdgeBasedGraph(nodes, edges)
	log.Println("Success!")
	rtree := NewRTree(geoNodes, segs)
	ret := &RoutingEngine{
		rtree:            rtree,
		edge_based_graph: ebg,
	}

	return ret
}
