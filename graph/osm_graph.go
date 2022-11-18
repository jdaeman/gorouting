package graph

import (
	"errors"
	"math"
	"sort"
)

// type ExternalNode struct {
// 	Id int64   // osm node id
// 	X  float64 // longitude
// 	Y  float64 // latitude
// }

type ResultNode struct {
	Id int64   // osm node id
	X  float64 // longitude
	Y  float64 // latitude
}

type ResultWay struct {
	Id     int64
	Nodes  []int64
	Oneway bool
}

type ResultRestriction struct {
	From int64
	Via  int64
	To   int64

	Only bool
}

type InternalRestriction struct {
	From int32
	Via  int32
	To   int32

	Only bool
}

type NodeBasedEdgeWithOSM struct {
	OsmId   int64 // osm way id
	OsmFrom int64 // osm node id
	OsmTo   int64 // osm node id

	Oneway bool

	Edge NodeBasedEdge
}

type ExternalEdge struct {
	Id   int64 // osm way id
	From int64 // osm node id
	To   int64 // oms node id

	AnnoId   int32 // annotation id
	SegIndex int16

	Oneway bool
	//... other attributes
}

type EdgeAnnotation struct {
	Id int64 // osm way id
	// etc way attributes
}

func NewEdgeAnnotation(id int64) *EdgeAnnotation {
	ret := &EdgeAnnotation{Id: id}
	return ret
}

type InternalNode struct {
	Id int32 // osm node id -> indexed id
}

type InternalEdge struct {
	From     int32 // indexed node id
	To       int32 // indexed node id
	Distance int32

	AnnotationId int32
	GeometryId   uint32

	Forward  bool // From -> To. always true
	Backward bool
	Split    bool
}

type OsmGraph struct {
	maxNode int32
	graph   [][]InternalEdge
	// trun restriction map
}

func (Graph *OsmGraph) GetNodeCount() int32 {
	return Graph.maxNode
}

func NewOsmGraph(nodeCnt int32, edges []InternalEdge) *OsmGraph {
	graph := make([][]InternalEdge, nodeCnt)
	for _, edge := range edges {
		from := edge.From
		graph[from] = append(graph[from], edge)
	}

	ret := &OsmGraph{maxNode: nodeCnt, graph: graph}
	return ret
}

func (Graph *OsmGraph) FindEdge(from, to int32) *InternalEdge {
	for idx, edge := range Graph.graph[from] {
		if to == edge.To && edge.Distance < math.MaxInt32 {
			return &Graph.graph[from][idx]
		}
	}
	return nil
}

func (Graph *OsmGraph) FindConstEdge(from, to int32) (InternalEdge, error) {
	ret := Graph.FindEdge(from, to)
	if ret != nil {
		return *ret, nil
	}
	return InternalEdge{}, errors.New("Not found edge")
}

func (Graph *OsmGraph) GetOutDegree(from int32) int {
	return len(Graph.graph[from])
}

func (Graph *OsmGraph) DelEdge(from, to int32) {
	edge := Graph.FindEdge(from, to)
	if edge == nil {
		return
	}

	lastDegree := len(Graph.graph[from])
	if lastDegree == 0 {
		panic("do not delete edge")
	}

	edge.Distance = math.MaxInt32
	sort.Slice(Graph.graph[from], func(l, r int) bool {
		return Graph.graph[from][l].Distance < Graph.graph[from][r].Distance
	})
	Graph.graph[from] = Graph.graph[from][:lastDegree-1]
}

func (Graph *OsmGraph) GetTarget(from int32, offset int) int32 {
	if offset >= len(Graph.graph[from]) {
		return -1
	}

	return Graph.graph[from][offset].To
}

func (Graph *OsmGraph) SetNewTarget(from int32, to int32, new int32) {
	graph := Graph.graph
	for idx, edge := range graph[from] {
		if edge.To == to {
			graph[from][idx].To = new
		}
	}
}
