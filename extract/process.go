package extract

import (
	"files"
	"graph"
	"log"
	"math"
	"sort"
	"util"

	"github.com/paulmach/osm"
)

type Extractor struct {
	rawFilePath string

	// raw osm data
	osmNodes     osm.Objects
	osmWays      osm.Objects
	osmRelations osm.Objects

	// parsing raw osm data
	AllNodes        []graph.ResultNode
	AllEdges        []graph.ResultWay
	AllRestrictions []graph.ResultRestriction

	//
	UsedNodes []int64

	//
	NodeBasedEdges  []graph.NodeBasedEdge // deprecate
	InternalEdges   []graph.InternalEdge
	EdgeAnnotations []graph.EdgeAnnotation
	Geometries      []graph.Geometry

	//
	Restrictions []graph.InternalRestriction
}

func NewExtractor(nodes, ways, restrictions osm.Objects, output string) *Extractor {
	ret := &Extractor{}

	ret.osmNodes = nodes
	ret.osmWays = ways
	ret.osmRelations = restrictions

	ret.rawFilePath = output

	return ret
}

// Parse osm node objests.
// osm node id.
// longitude, latitude.
func (extractor *Extractor) ProcessOSMNodes() int {
	osmNodes := &extractor.osmNodes
	allNodes := &extractor.AllNodes

	*allNodes = make([]graph.ResultNode, 0, len(*osmNodes))

	for _, node := range *osmNodes {
		resultNode := ParseOSMNode(node)
		*allNodes = append(*allNodes, *resultNode)
	}

	// sorted by osm node id.
	sort.Slice(*allNodes, func(l, r int) bool {
		return (*allNodes)[l].Id < (*allNodes)[r].Id
	})

	// dealloc unused memory.
	extractor.osmNodes = nil
	return len(*allNodes)
}

// Parse osm way objects.
// osm way id.
// nodes.
// attriubtes.
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

	// sorted by osm way id.
	sort.Slice(*allEdges, func(l, r int) bool {
		return (*allEdges)[l].Id < (*allEdges)[r].Id
	})

	// dealloc unused memory.
	extractor.osmWays = nil
	return len(*allEdges)
}

// Parse osm restriction objects.
// from osm way id.
// via osm node id.
// to osm way id.
func (extractor *Extractor) ProcessOSMRestrictions() int {
	osmRelations := &extractor.osmRelations
	allRestrictions := &extractor.AllRestrictions

	*allRestrictions = make([]graph.ResultRestriction, 0, len(*osmRelations))
	for _, relation := range *osmRelations {
		resultRestriction := ParseOSMRelation(relation)
		if resultRestriction == nil {
			continue
		}
		*allRestrictions = append(*allRestrictions, *resultRestriction)
	}

	sort.Slice(*allRestrictions, func(l, r int) bool {
		return (*allRestrictions)[l].From < (*allRestrictions)[r].From
	})

	// dealloc unused memory.
	extractor.osmRelations = nil
	return len(*allRestrictions)
}

func (extractor *Extractor) ProcessNodes() {
	extractor.collectUsedNodes()
}

func (extractor *Extractor) splitCrossWays() {
	allWays := &extractor.AllEdges

	// default value is zero.
	// some node affect several ways.
	sections := make([]int, len(extractor.UsedNodes))
	for _, way := range *allWays {
		if len(way.Nodes) <= 1 {
			panic("way node count is less than 2")
		}
		//          ^
		//          |
		//          |
		//   <----- u -----> way: A
		//          |
		//          |
		//          v  way: B
		//
		// node 'u' is in way 'A' and way 'B'
		// in this case, split way 'A' and way 'B'
		// [way 'A1'], [way 'A2'], [way 'B1'], [way 'B2']
		for _, osmNode := range way.Nodes {
			id := extractor.getInternalNodeId(osmNode)
			sections[id] += 1
		}
	}

	ways := make([]graph.ResultWay, 0, len(*allWays))
	for _, way := range *allWays {
		wayId, oneway := way.Id, way.Oneway
		nodes := []int64{way.Nodes[0]}

		var i int
		// source ~ [shape vertex...] ~ target
		for i = 1; i < len(way.Nodes)-1; i++ {
			osmNode := way.Nodes[i]
			nodes = append(nodes, osmNode)
			id := extractor.getInternalNodeId(osmNode)

			if sections[id] > 1 {
				// split way
				ways = append(ways, graph.ResultWay{Id: wayId, Nodes: nodes, Oneway: oneway})
				nodes = []int64{osmNode}
			}
		}

		nodes = append(nodes, way.Nodes[i])
		ways = append(ways, graph.ResultWay{Id: wayId, Nodes: nodes, Oneway: oneway})
	}

	// sorted by way id.
	sort.Slice(ways, func(l, r int) bool {
		return ways[l].Id < ways[r].Id
	})

	fromCount, toCount := len(*allWays), len(ways)
	*allWays = ways

	log.Println("Split ways", fromCount, "increase to", toCount)
}

func (extractor *Extractor) collectUsedNodes() {
	allEdges := &extractor.AllEdges
	usedNodes := &extractor.UsedNodes

	// node id that is construct way.
	*usedNodes = make([]int64, 0)

	for _, edge := range *allEdges {
		if len(edge.Nodes) <= 1 {
			panic("way node count is less than 2")
		}
		for _, u := range edge.Nodes {
			*usedNodes = append(*usedNodes, u)
		}
	}

	sort.Slice(*usedNodes, func(l, r int) bool {
		return (*usedNodes)[l] < (*usedNodes)[r]
	})

	// There will be duplicated node.
	// do not erase duplicated node in this step.
}

func (extractor *Extractor) ProcessEdges() {
	extractor.splitCrossWays()
}

func (extractor *Extractor) ProcessRestrictions() {
	allEdges := extractor.AllEdges
	allRestrictions := extractor.AllRestrictions

	// Get osm way id index.
	// [from, to)
	getWays := func(wayId int64) [2]int {
		l1 := sort.Search(len(allEdges), func(i int) bool {
			return allEdges[i].Id >= wayId
		})
		if l1 >= len(allEdges) || allEdges[l1].Id != wayId {
			return [2]int{l1, l1}
		}
		l2 := sort.Search(len(allEdges), func(i int) bool {
			return allEdges[i].Id > wayId
		})
		return [2]int{l1, l2}
	}

	// Get source, target osm node id.
	getNodes := func(wayIndex int) [2]int64 {
		way := allEdges[wayIndex]
		nodes := way.Nodes
		return [2]int64{nodes[0], nodes[len(nodes)-1]}
	}

	// from way -> via node -> to way
	// (u ~> v)        v       (v ~> w)
	//
	// from node -> via node -> to node
	//     u            v          w
	for pos, restriction := range allRestrictions {
		var fromNode, toNode int64 // osm node id
		viaNode := restriction.Via
		fromNode, toNode = -1, -1 // invalid node id

		fromWayRange := getWays(restriction.From)
		toWayRange := getWays(restriction.To)

		// search
		for from := fromWayRange[0]; from < fromWayRange[1]; from++ {

			fromWayNodes := getNodes(from)
			if fromWayNodes[0] != viaNode && fromWayNodes[1] != viaNode {
				continue
			}

			if viaNode == fromWayNodes[0] {
				fromNode = fromWayNodes[1]
			} else {
				fromNode = fromWayNodes[0]
			}

			for to := toWayRange[0]; to < toWayRange[1]; to++ {
				toWayNodes := getNodes(to)
				if toWayNodes[0] != viaNode && toWayNodes[1] != viaNode {
					continue
				}

				if viaNode == toWayNodes[0] {
					toNode = toWayNodes[1]
				} else {
					toNode = toWayNodes[0]
				}
			}

			if fromNode != -1 && toNode != -1 {
				break
			}
		}

		if fromNode == -1 || toNode == -1 {
			fromNode, toNode = math.MaxInt64, math.MaxInt64
		}

		allRestrictions[pos].From = fromNode
		allRestrictions[pos].To = toNode
	}
}

func (extractor *Extractor) PrepareData() {
	// node
	extractor.prepareNodes()
	extractor.writeNodes()

	// edges...
	extractor.prepareEdges()
	extractor.writeEdges()

	// restrictions...
	extractor.prepareRestrictions()
	extractor.writeRestrictions()
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
	// in this step,
	// used node id must be unique.
	uniqueNodes := make([]int64, 0, len(usedNodes))
	for _, node := range usedNodes {
		tail := len(uniqueNodes) - 1
		if tail == -1 || uniqueNodes[tail] != node {
			uniqueNodes = append(uniqueNodes, node)
		}
	}
	extractor.UsedNodes = nil

	// remove duplicated locations
	uniqueGeos := make([]graph.ResultNode, 0, len(uniqueNodes))
	for _, node := range uniqueNodes {
		// find node location
		idx := sort.Search(len(allNodes), func(idx int) bool {
			return allNodes[idx].Id >= node
		})
		// valid node location
		if idx < len(allNodes) && allNodes[idx].Id == node {
			uniqueGeos = append(uniqueGeos, allNodes[idx])
		}
	}

	// both nodes are sorted.
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
		panic("Could not find node id")
	}
}

func (extractor *Extractor) prepareEdges() {
	geoNodes := extractor.AllNodes

	_fwdEdges := make([]graph.InternalEdge, 0)
	geometries := make([]graph.Geometry, 0, len(extractor.AllEdges))
	annotations := make([]graph.EdgeAnnotation, 0, len(extractor.AllEdges))

	// convert osm Way to edge.
	for _, edge := range extractor.AllEdges {
		edgeId, oneway := edge.Id, edge.Oneway
		osmNodes := edge.Nodes
		nodes := make([]int32, 0, len(osmNodes))
		distance := int32(0)

		geometry := graph.NewGeometry(len(osmNodes))
		annotation := graph.NewEdgeAnnotation(edgeId)

		reverse := false
		forward, backward := true, !oneway

		source, target := osmNodes[0], osmNodes[len(osmNodes)-1]
		if source > target {
			// program rule.
			// always source is less than target.
			// except oneway edge.
			head, tail := 0, len(osmNodes)-1
			// reverse geometry if not oneway.
			for !edge.Oneway && head <= tail {
				osmNodes[head], osmNodes[tail] = osmNodes[tail], osmNodes[head]
				head += 1
				tail -= 1
			}
			// reverse direction
			// if oneway, after that fix this.
			reverse = oneway
		}

		// to internal node id
		for _, osmNode := range osmNodes {
			internalId := extractor.getInternalNodeId(osmNode)
			nodes = append(nodes, internalId)
			geometry.Nodes = append(geometry.Nodes, internalId)
		}

		// edge length, used weight.
		for i := 1; i < len(nodes); i++ {
			from, to := nodes[i-1], nodes[i]
			coord1 := [2]float64{geoNodes[from].X, geoNodes[from].Y}
			coord2 := [2]float64{geoNodes[to].X, geoNodes[to].Y}

			segDistance := int32(util.HaversineDistance(coord1, coord2))
			geometry.Distances = append(geometry.Distances, segDistance)
			distance += segDistance
		}

		from, to := nodes[0], nodes[len(nodes)-1]
		if reverse {
			from, to = to, from
			forward, backward = backward, forward
		}
		_fwdEdge := graph.InternalEdge{
			From: from, To: to,
			Distance: distance,
			Forward:  forward, Backward: backward,
			AnnotationId: int32(len(annotations)),
			GeometryId:   uint32(len(geometries)),
		}

		_fwdEdges = append(_fwdEdges, _fwdEdge)
		geometries = append(geometries, *geometry)
		annotations = append(annotations, *annotation)
	}

	sort.Slice(_fwdEdges, func(l, r int) bool {
		left, right := &_fwdEdges[l], &_fwdEdges[r]
		if left.From == right.From {
			return left.To < right.To
		}
		return left.From < right.From
	})

	// find minimal edge in both directions
	for i := 0; i < len(_fwdEdges); {
		startIdx := i
		from, to := _fwdEdges[i].From, _fwdEdges[i].To
		minForward, minBackward := math.MaxInt32, math.MaxInt32
		minFwdIdx, minBwdIdx := -1, -1

		for i < len(_fwdEdges) &&
			_fwdEdges[i].From == from &&
			_fwdEdges[i].To == to {

			edge := &_fwdEdges[i]

			if edge.Forward && edge.Distance < int32(minForward) {
				minForward = int(edge.Distance)
				minFwdIdx = i
			}
			if edge.Backward && edge.Distance < int32(minBackward) {
				minBackward = int(edge.Distance)
				minBwdIdx = i
			}
			i++
		}

		if minFwdIdx == minBwdIdx {
			// maybe, bidirection way
			_fwdEdges[minFwdIdx].Forward = true
			_fwdEdges[minFwdIdx].Backward = true
			_fwdEdges[minFwdIdx].Split = false
		} else {
			// maybe, oneway.
			// but,
			// https://www.openstreetmap.org/way/608946830#map=19/37.57055/126.99844
			// https://www.openstreetmap.org/way/218835937#map=18/37.57049/126.99852
			has_forward := minFwdIdx != -1
			has_backward := minBwdIdx != -1

			if has_forward {
				// has forward direction
				_fwdEdges[minFwdIdx].Forward = true
				_fwdEdges[minFwdIdx].Backward = false
				_fwdEdges[minFwdIdx].Split = has_backward
			}
			if has_backward {
				// has backward direction
				_fwdEdges[minBwdIdx].From, _fwdEdges[minBwdIdx].To =
					_fwdEdges[minBwdIdx].To, _fwdEdges[minBwdIdx].From
				_fwdEdges[minBwdIdx].Forward = true
				_fwdEdges[minBwdIdx].Backward = false
				_fwdEdges[minBwdIdx].Split = has_forward
			}
		}

		for j := startIdx; j < i; j++ {
			if j == minFwdIdx || j == minBwdIdx {
				continue
			}
			_fwdEdges[j].From, _fwdEdges[j].To = math.MaxInt32, math.MaxInt32
		}
	}

	// re-sort
	sort.Slice(_fwdEdges, func(l, r int) bool {
		left, right := &_fwdEdges[l], &_fwdEdges[r]
		if left.From == right.From {
			return left.To < right.To
		}
		return left.From < right.From
	})

	// no use invalid edge.
	invalidIdx := sort.Search(len(_fwdEdges), func(i int) bool {
		return _fwdEdges[i].From >= math.MaxInt32
	})
	_fwdEdges = _fwdEdges[:invalidIdx]

	extractor.InternalEdges = _fwdEdges
	extractor.Geometries = geometries
	extractor.EdgeAnnotations = annotations
}

func (extractor *Extractor) prepareRestrictions() {
	allRestrictions := &extractor.AllRestrictions

	// soretd by from osm node id
	sort.Slice(*allRestrictions, func(l, r int) bool {
		return (*allRestrictions)[l].From < (*allRestrictions)[r].From
	})
	// after tail data is invalid.
	tail := sort.Search(len(*allRestrictions), func(i int) bool {
		return (*allRestrictions)[i].From >= math.MaxInt64
	})

	restrictions := make([]graph.InternalRestriction, 0, tail)
	// until tail...
	for _, externalRestriction := range (*allRestrictions)[:tail] {
		FromOsm, ViaOsm, ToOsm := externalRestriction.From,
			externalRestriction.Via,
			externalRestriction.To
		from, via, to := extractor.getInternalNodeId(FromOsm),
			extractor.getInternalNodeId(ViaOsm),
			extractor.getInternalNodeId(ToOsm)

		restriction := graph.InternalRestriction{From: from, Via: via, To: to, Only: externalRestriction.Only}
		restrictions = append(restrictions, restriction)
	}

	// sorted by from internal node id.
	sort.Slice(restrictions, func(l, r int) bool {
		r1, r2 := &restrictions[l], &restrictions[r]
		return r1.From < r2.From
	})

	extractor.Restrictions = restrictions
	*allRestrictions = nil
}

func (extractor *Extractor) writeNodes() {
	log.Println("Location/Unique node count", len(extractor.AllNodes))

	newFile := files.ToDataPath(extractor.rawFilePath, files.GEONODE)
	err := files.StoreGeoNodes(newFile, extractor.AllNodes)
	if err == nil {
		log.Println("Saved to", newFile)
	} else {
		log.Println("ERROR", err)
	}
}

func (extractor *Extractor) writeEdges() {
	log.Println("Edge count", len(extractor.InternalEdges))

	newFile := files.ToDataPath(extractor.rawFilePath, files.NBGEDGE)
	err := files.StoreEdges(newFile, extractor.InternalEdges)
	if err == nil {
		log.Println("Saved to", newFile)
	} else {
		log.Println("ERROR", err)
	}

	newFile = files.ToDataPath(extractor.rawFilePath, files.ANNOTATION)
	err = files.StoreEdgeAnnotations(newFile, extractor.EdgeAnnotations)
	if err == nil {
		log.Println("Saved to", newFile)
	} else {
		log.Println("ERROR", err)
	}

	newFile = files.ToDataPath(extractor.rawFilePath, files.GEOMETRY)
	err = files.StoreEdgeGeometries(newFile, extractor.Geometries)
	if err == nil {
		log.Println("Saved to", newFile)
	} else {
		log.Println("ERROR", err)
	}
}

func (extractor *Extractor) writeRestrictions() {
	log.Println("Restriction count", len(extractor.Restrictions))

	newFile := files.ToDataPath(extractor.rawFilePath, files.RESTRICTION)
	err := files.StoreTurnRestrictions(newFile, extractor.Restrictions)
	if err == nil {
		log.Println("Saved to", newFile)
	} else {
		log.Println("ERROR", err)
	}
}
