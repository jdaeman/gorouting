package extract

import (
	"files"
	"graph"
	"log"
	"sort"
	"util"

	"github.com/paulmach/osm"
)

type Extractor struct {
	rawFilePath string

	// raw osm data
	osmNodes        osm.Objects
	osmWays         osm.Objects
	osmRestrictions osm.Objects

	// parsing raw osm data
	AllNodes        []graph.ExternalNode
	AllEdges        []graph.ResultWay
	EdgeAnnotations []graph.EdgeAnnotation

	UsedNodes []int64
	UsedEdges []graph.ExternalEdge

	InternalEdges []graph.InternalEdge
}

func NewExtractor(nodes, ways, restrictions osm.Objects, output string) *Extractor {
	ret := &Extractor{}

	ret.osmNodes = nodes
	ret.osmWays = ways
	ret.osmRestrictions = restrictions

	ret.rawFilePath = output

	return ret
}

func (extractor *Extractor) ProcessOSMNodes() int {
	osmNodes := &extractor.osmNodes
	allNodes := &extractor.AllNodes

	*allNodes = make([]graph.ExternalNode, 0, len(*osmNodes))

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

	*allEdges = make([]graph.ResultWay, 0, len(*osmWays))
	for _, way := range *osmWays {
		resultWay := ParseOSMWay(way)
		if resultWay == nil {
			continue
		}
		*allEdges = append(*allEdges, *resultWay)
	}

	extractor.osmWays = nil
	return len(*allEdges)
}

func (extractor *Extractor) ProcessOSMRestriction() int {
	return 0
}

func (extractor *Extractor) ProcessNodes() {
	// do nothing
}

func (extractor *Extractor) processEdge(edge *graph.ResultWay) {
	if len(edge.Nodes) <= 1 {
		return
	}

	usedNodes := &extractor.UsedNodes
	usedEdges := &extractor.UsedEdges
	edgeAnnotations := &extractor.EdgeAnnotations
	annoId := len(*edgeAnnotations)

	var from, to int64
	var segIndex int16

	segIndex = 0
	for v := 0; v+1 < len(edge.Nodes); v++ {
		from = edge.Nodes[v]
		to = edge.Nodes[v+1]

		edge := graph.ExternalEdge{Id: edge.Id, From: from, To: to, Oneway: edge.Oneway, AnnoId: int32(annoId), SegIndex: segIndex}
		*usedEdges = append(*usedEdges, edge)
		*usedNodes = append(*usedNodes, from)
		segIndex += 1
	}
	*usedNodes = append(*usedNodes, to)

	*edgeAnnotations = append(*edgeAnnotations, graph.EdgeAnnotation{Id: edge.Id})
}

func (extractor *Extractor) ProcessEdges() {
	allEdges := &extractor.AllEdges
	usedNodes := &extractor.UsedNodes
	usedEdges := &extractor.UsedEdges
	edgeAnnotations := &extractor.EdgeAnnotations

	*usedNodes = make([]int64, 0)
	*usedEdges = make([]graph.ExternalEdge, 0)
	*edgeAnnotations = make([]graph.EdgeAnnotation, 0)

	for _, edge := range *allEdges {
		extractor.processEdge(&edge)
	}
}

func (extractor *Extractor) writeNodes() {
	log.Println("Location node count", len(extractor.AllNodes))
	log.Println("Unique node count", len(extractor.UsedNodes))

	newFile := files.ToDataPath(extractor.rawFilePath, ".geos")
	err := files.StoreGeoNodes(newFile, extractor.AllNodes)
	if err == nil {
		log.Println("Saved to", newFile)
	} else {
		log.Println("ERROR", err)
	}
}

func (extractor *Extractor) writeEdges() {
	log.Println("Uncomp edge count", len(extractor.InternalEdges))
	log.Println("Edge annotation count", len(extractor.EdgeAnnotations))

	newFile := files.ToDataPath(extractor.rawFilePath, ".uncomp")
	err := files.StoreUncompEdges(newFile, extractor.InternalEdges)
	if err == nil {
		log.Println("Saved to", newFile)
	} else {
		log.Println("ERROR", err)
	}

	newFile = files.ToDataPath(extractor.rawFilePath, ".anno")
	err = files.StoreEdgeAnnotations(newFile, extractor.EdgeAnnotations)
	if err == nil {
		log.Println("Saved to", newFile)
	} else {
		log.Println("ERROR", err)
	}
}

func (extractor *Extractor) PrepareData() {
	// node
	extractor.prepareNodes()
	extractor.writeNodes()

	// edges...
	extractor.prepareEdges()
	extractor.writeEdges()
}

func (extractor *Extractor) prepareNodes() {
	allNodes := extractor.AllNodes
	usedNodes := extractor.UsedNodes

	// sorted by osm node id
	sort.Slice(usedNodes, func(l, r int) bool {
		return usedNodes[l] < usedNodes[r]
	})
	sort.Slice(allNodes, func(l, r int) bool {
		return allNodes[l].Id < allNodes[r].Id
	})

	// remove duplicated elements
	uniqueNodes := make([]int64, 0, len(usedNodes)/2)
	for _, node := range usedNodes {
		// find node id
		tail := len(uniqueNodes) - 1
		if tail == -1 || uniqueNodes[tail] != node {
			uniqueNodes = append(uniqueNodes, node)
		}
	}

	// remove duplicated locations
	extractor.UsedNodes = nil
	uniqueGeos := make([]graph.ExternalNode, 0, len(usedNodes)/2)
	for _, node := range uniqueNodes {
		// find node location
		idx := sort.Search(len(allNodes), func(idx int) bool {
			return allNodes[idx].Id >= node
		})
		if idx < len(allNodes) && allNodes[idx].Id == node {
			uniqueGeos = append(uniqueGeos, allNodes[idx])
		}
	}

	extractor.UsedNodes = uniqueNodes
	extractor.AllNodes = uniqueGeos
}

func (extractor *Extractor) getInternalNodeId(osmId int64) int32 {
	count := len(extractor.UsedNodes)
	newId := sort.Search(count, func(idx int) bool {
		return extractor.UsedNodes[idx] >= osmId
	})

	if newId < count && extractor.UsedNodes[newId] == osmId {
		return int32(newId)
	} else {
		return -1
	}
}

func (extractor *Extractor) prepareEdges() {

	uncompEdges := &extractor.InternalEdges
	usedEdges := extractor.UsedEdges
	geoNodes := extractor.AllNodes

	sort.Slice(usedEdges, func(l, r int) bool {
		if usedEdges[l].From == usedEdges[r].From {
			return usedEdges[l].To < usedEdges[r].To
		}
		return usedEdges[l].From < usedEdges[r].From
	})

	*uncompEdges = make([]graph.InternalEdge, 0, len(usedEdges)*2)
	for _, edge := range usedEdges {
		from, to := edge.From, edge.To
		internalFrom, internalTo := extractor.getInternalNodeId(from), extractor.getInternalNodeId(to)

		if internalFrom == -1 || internalTo == -1 {
			panic("Invalid node found from edge")
		}

		coord1 := [2]float64{geoNodes[internalFrom].X, geoNodes[internalFrom].Y}
		coord2 := [2]float64{geoNodes[internalTo].X, geoNodes[internalTo].Y}
		distance := int32(util.HaversineDistance(coord1, coord2))
		segIndex := edge.SegIndex

		fwdEdge := graph.InternalEdge{From: internalFrom, To: internalTo, Distance: distance,
			Forward: true, Reverse: false,
			AnnoId: edge.AnnoId, SegIndex: segIndex}
		revEdge := graph.InternalEdge{From: internalTo, To: internalFrom, Distance: distance,
			Forward: !edge.Oneway, Reverse: edge.Oneway,
			AnnoId: edge.AnnoId, SegIndex: segIndex}

		*uncompEdges = append(*uncompEdges, fwdEdge)
		*uncompEdges = append(*uncompEdges, revEdge)
	}

	sort.Slice(*uncompEdges, func(l, r int) bool {
		if (*uncompEdges)[l].From == (*uncompEdges)[r].From {
			return (*uncompEdges)[l].To < (*uncompEdges)[r].To
		}

		return (*uncompEdges)[l].From < (*uncompEdges)[r].From
	})

	extractor.UsedEdges = nil
}
