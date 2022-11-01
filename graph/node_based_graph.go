package graph

type NodeBasedNode struct {
	ID int32
}

type NodeBasedEdge struct {
	From   int32
	To     int32
	Weight int32
}

type NodeBasedGraph struct {
	adjList [][]NodeBasedEdge // fwd graph, outgoing edges
}

func NewNodeBasedGraph(nodes int32, edges []NodeBasedEdge) *NodeBasedGraph {
	if nodes == 0 || len(edges) == 0 {
		return nil
	}

	graph := make([][]NodeBasedEdge, nodes)
	//revGraph := make([][]NodeBasedEdge, nodes)
	for i := range graph {
		graph[i] = make([]NodeBasedEdge, 0)
		//revGraph[i] = make([]NodeBasedEdge, 0)
	}

	for _, edge := range edges {
		if edge.From >= nodes || edge.To >= nodes {
			panic("edge has invalid node")
		}

		graph[edge.From] = append(graph[edge.From], edge)
		//revGraph[edge.To] = append(revGraph[edge.To], edge)
	}

	nbg := &NodeBasedGraph{adjList: graph}
	return nbg
}
