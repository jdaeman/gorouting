package extract

import (
	"modules/files"
	"modules/graph"
	"sort"
)

type NodeBasedGraphFactory struct {
	nbg *graph.NodeBasedGraph

	restrictions []graph.InternalRestriction
	geometries   []graph.Geometry
}

func NewNodeBasedGraphFactory(datapath string) *NodeBasedGraphFactory {
	dataReader := files.NewReader(datapath)

	nodes := dataReader.LoadGeoNodes()
	edges := dataReader.LoadEdges()
	restrictions := dataReader.LoadRestrictions()
	geometries := dataReader.LoadEdgeGeometries()

	nodeCount := int32(len(nodes))
	nbg := graph.NewNodeBasedGraph(nodeCount, edges)

	ret := &NodeBasedGraphFactory{nbg: nbg, restrictions: restrictions, geometries: geometries}
	return ret
}

func (factory *NodeBasedGraphFactory) GetGraph() *graph.NodeBasedGraph {
	return factory.nbg
}

func (factory *NodeBasedGraphFactory) GetIncomingEdges(via int32) []int32 {
	nbg := factory.nbg

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

func (factory *NodeBasedGraphFactory) GetOutgoingEdges(via int32) []int32 {
	nbg := factory.nbg

	curEdge, lastEdge := nbg.BeginEdges(via), nbg.EndEdges(via)
	edges := make([]int32, 0)

	for ; curEdge < lastEdge; curEdge++ {
		if nbg.GetEdgeData(curEdge).Reverse {
			continue
		}
		edges = append(edges, curEdge)
	}

	return edges
}

func (factory *NodeBasedGraphFactory) IsTurnallowed(from, via, to int32) bool {
	restrictions := factory.restrictions

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

		// find from, via
		if restriction.Only {
			return restriction.To == to
		} else {
			// no restriction
			if restriction.To == to {
				return false
			}
		}
	}

	return true
}

func (factory *NodeBasedGraphFactory) GetGeometry(geoId uint32) []int32 {
	fwdGeoId := geoId & ((1 << 31) - 1)
	fwdGeos := factory.geometries[fwdGeoId].Nodes
	nodes := make([]int32, len(fwdGeos))
	copy(nodes, fwdGeos)

	if (geoId & (1 << 31)) > 0 {
		// reverse geometry
		head, tail := 0, len(nodes)-1
		for head < tail {
			nodes[head], nodes[tail] = nodes[tail], nodes[head]
			head++
			tail--
		}
	}

	return nodes
}
