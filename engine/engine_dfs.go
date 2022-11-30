package engine

import (
	"graph"
	"sort"
)

func dfs(graph *graph.EdgeBasedGraph, visit []int32, counter *int32, u int32) {
	*counter += 1
	visit[u] = *counter

	for _, adjEdge := range graph.GetForwardEdges(u) {
		v := graph.GetEdgeData(adjEdge).Target

		if visit[v] == 0 {
			dfs(graph, visit, counter, v)
		}
	}
}

func (engine *RoutingEngine) Dfs(u int32) []int32 {
	maxNode := engine.edge_based_graph.GetNumberOfNodes()

	counter := int32(0)
	visit := make([]int32, maxNode)

	dfs(engine.edge_based_graph, visit, &counter, u)

	pathOrder := make([]int32, 0, counter)
	for u := int32(0); u < maxNode; u++ {
		if visit[u] > 0 {
			pathOrder = append(pathOrder, u)
		}
	}

	sort.Slice(pathOrder, func(l, r int) bool {
		u, v := pathOrder[l], pathOrder[r]
		return visit[u] < visit[v]
	})

	return pathOrder
}
