package extract

import (
	"sort"

	"github.com/paulmach/osm"
)

type ProcNode struct {
	Id int
}

type ProcEdge struct {
	From int
	To   int

	Oneway bool
}

func ProcessNodeObjs(objs osm.Objects) []ParsingNode {
	allnode := make([]ParsingNode, 0, len(objs))

	for _, node := range objs {
		resultNode := ParseOSMNode(node)
		allnode = append(allnode, *resultNode)
	}

	return allnode
}

func ProcessWayObjs(objs osm.Objects) []ParsingWay {
	allway := make([]ParsingWay, 0, len(objs))
	for _, way := range objs {
		resultWay := ParseOSMWay(way)
		allway = append(allway, *resultWay)
	}

	return allway
}

func ProcessRestrictionObjs(objs osm.Objects) []ParsingRestriction {
	allRestriction := make([]ParsingRestriction, 0, len(objs))
	for _, restriction := range objs {
		resultRestriction := ParseOSMRestriction(restriction)
		allRestriction = append(allRestriction, *resultRestriction)
	}

	return allRestriction
}

///////////////////////////////////////////////////////////////////////////////

func ProcessNode(node ParsingNode) {
	// do nothing
}

func ProcessWay(way ParsingWay) []ProcEdge {
	if len(way.Nodes) <= 1 {
		return nil
	}

	nodes := make([]ProcNode, 0, len(way.Nodes))
	edges := make([]ProcEdge, 0, len(way.Nodes)-1)

	for v := 0; v+1 < len(way.Nodes); v++ {
		from, to := way.Nodes[v], way.Nodes[v+1]
		edges = append(edges, ProcEdge{from, to, way.Oneway})
		nodes = append(nodes, ProcNode{from})
		if v+1 == len(way.Nodes)-1 {
			nodes = append(nodes, ProcNode{to})
		}
	}

	return edges
}

func PrepareData(nodes []ProcNode, edges []ProcEdge) {
	// node
	// edges
}

func prepareNodes(nodes []ProcNode) {
	// sort by osm node id
	sort.Slice(nodes, func(l, r int) bool {
		return nodes[l].Id < nodes[r].Id
	})

	// erase duplicated nodes
	nodes = eraseDuplicatedNode(nodes)
}

func eraseDuplicatedNode(nodes []ProcNode) []ProcNode {

	index, value := 0, nodes[0]
	cpy := nodes
	for len(cpy) > 0 {
		index += 1

		idx := sort.Search(len(cpy), func(idx int) bool {
			return cpy[idx].Id > value.Id
		})

		if idx >= len(cpy) {
			cpy = nil
		} else {
			nodes[index] = cpy[idx]
			value = cpy[idx]
			cpy = cpy[idx:]
		}
	}

	return nodes[:index]
}

func prepareEdges(edges []ProcEdge) {

}
