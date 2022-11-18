package extract

import (
	"files"
	"graph"
	"sort"
)

type NodeBasedGraphFactory struct {
	nbg *graph.NodeBasedGraph

	restrictions []graph.InternalRestriction
	geometries   []graph.Geometry
}

func NewNodeBasedGraphFactory(datapath string) *NodeBasedGraphFactory {
	nodes := files.LoadGeoNodes(files.ToDataPath(datapath, files.GEONODE))
	edges := files.LoadEdges(files.ToDataPath(datapath, files.NBGEDGE))
	restrictions := files.LoadRestrictions(files.ToDataPath(datapath, files.RESTRICTION))
	geometries := files.LoadEdgeGeometries(files.ToDataPath(datapath, files.GEOMETRY))

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
		if !nbg.GetEdgeData(curEdge).Reverse {
			edges = append(edges, curEdge)
		}
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

func (factory *NodeBasedGraphFactory) GetGeometry(geoId int32) []int32 {
	geoId &= ((1 << 31) - 1)
	return factory.geometries[geoId].Nodes
}
