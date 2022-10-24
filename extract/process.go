package extract

import (
	"sort"

	"github.com/paulmach/osm"
)

type ExternalNode struct {
	Id int
}

type ExternalEdge struct {
	From int
	To   int

	Oneway bool
}

type Extractor struct {
	osmNodes        osm.Objects
	osmWays         osm.Objects
	osmRestrictions osm.Objects

	AllNodes []ParsingNode
	AllEdges []ParsingWay

	UsedNodes []ExternalNode
	UsedEdges []ExternalEdge

	InternalNodes []NodeBasedNode
	InternalEdges []NodeBasedEdge
}

type ExtractorInterface interface {
	//
	ProcessOSMNodes() int
	ProcessOSMWays()
	ProcessOSMRestrictions() bool

	//
	ProcessNodes()
	ProcessEdges()
	processEdge(edge *ParsingWay)

	//
	PrepareData()
	prepareNodes()
	prepareEdges()
}

func NewExtractor(nodes, ways, restrictions osm.Objects) *Extractor {
	ret := &Extractor{}

	ret.osmNodes = nodes
	ret.osmWays = ways
	ret.osmRestrictions = restrictions

	return ret
}

func (extractor *Extractor) ProcessOSMNodes() int {
	osmNodes := &extractor.osmNodes
	allNodes := &extractor.AllNodes

	*allNodes = make([]ParsingNode, 0, len(*osmNodes))

	for _, node := range *osmNodes {
		resultNode := ParseOSMNode(node)
		*allNodes = append(*allNodes, *resultNode)
	}

	extractor.osmNodes = nil
	return len(*allNodes)
}

func (extractor *Extractor) ProcessOSMWays() int {
	osmWays := &extractor.osmWays
	allEdges := &extractor.AllEdges

	*allEdges = make([]ParsingWay, 0, len(*osmWays))
	for _, way := range *osmWays {
		resultWay := ParseOSMWay(way)
		*allEdges = append(*allEdges, *resultWay)
	}

	extractor.osmWays = nil
	return len(*allEdges)
}

// func ProcessOSMRestriction(objs osm.Objects) []ParsingRestriction {
// 	allRestriction := make([]ParsingRestriction, 0, len(objs))
// 	for _, restriction := range objs {
// 		resultRestriction := ParseOSMRestriction(restriction)
// 		allRestriction = append(allRestriction, *resultRestriction)
// 	}

// 	return allRestriction
// }

///////////////////////////////////////////////////////////////////////////////

func (extractor *Extractor) ProcessNodes() {
	// do nothing
}

func (extractor *Extractor) processEdge(edge *ParsingWay) {
	if len(edge.Nodes) <= 1 {
		return
	}

	usedNodes := &extractor.UsedNodes
	usedEdges := &extractor.UsedEdges

	from, to := 0, 0
	for v := 0; v+1 < len(edge.Nodes); v++ {
		from = edge.Nodes[v]
		to = edge.Nodes[v+1]

		*usedEdges = append(*usedEdges, ExternalEdge{from, to, edge.Oneway})
		*usedNodes = append(*usedNodes, ExternalNode{from})
	}
	*usedNodes = append(*usedNodes, ExternalNode{to})
}

func (extractor *Extractor) ProcessEdges() {
	allEdges := &extractor.AllEdges
	usedNodes := &extractor.UsedNodes
	usedEdges := &extractor.UsedEdges

	*usedNodes = make([]ExternalNode, 0)
	*usedEdges = make([]ExternalEdge, 0)

	for _, edge := range *allEdges {
		extractor.processEdge(&edge)
	}
}

func (extractor *Extractor) PrepareData() {
	// node
	extractor.prepareNodes()
	// edges...
}

func (extractor *Extractor) prepareNodes() {
	allNodes := extractor.AllNodes
	usedNodes := extractor.UsedNodes

	// sorted by osm node id
	sort.Slice(usedNodes, func(l, r int) bool {
		return usedNodes[l].Id < usedNodes[r].Id
	})
	sort.Slice(allNodes, func(l, r int) bool {
		return allNodes[l].ID < allNodes[r].ID
	})

	// remove duplicated elements
	uniqueNodes := make([]ExternalNode, 0, len(usedNodes)/2)
	for _, node := range usedNodes {
		// find node id
		tail := len(uniqueNodes) - 1
		if tail == -1 || uniqueNodes[tail].Id != node.Id {
			uniqueNodes = append(uniqueNodes, node)
		}
	}

	// remove duplicated locations
	extractor.UsedNodes = nil
	uniqueGeos := make([]ParsingNode, 0, len(usedNodes)/2)
	for _, node := range uniqueNodes {
		// find node location
		idx := sort.Search(len(allNodes), func(idx int) bool {
			return allNodes[idx].ID >= node.Id
		})
		if idx < len(allNodes) && allNodes[idx].ID == node.Id {
			uniqueGeos = append(uniqueGeos, allNodes[idx])
		}
	}

	extractor.UsedNodes = uniqueNodes
	extractor.AllNodes = uniqueGeos
}

func (extractor *Extractor) prepareEdges() {

}
