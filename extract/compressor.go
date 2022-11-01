package extract

import (
	"files"
	"fmt"
	"graph"
)

var lastid int32

type Compressor struct {
	uncompGraph *graph.OsmGraph

	// key: annoId
	// value: geoId
	mappingGeoId map[int32]int32
	geometries   []graph.Geometry

	// geometry
	// node_based_edge 로 만드는 작업이 있어야 함.
}

func NewCompressor(rawFilePath string) *Compressor {
	geos := files.LoadGeoNodes(files.ToDataPath(rawFilePath, ".geos"))
	uncomps := files.LoadUncompEdges(files.ToDataPath(rawFilePath, ".uncomp"))
	//anno := files.LoadEdgeAnnotations(files.ToDataPath(rawFilePath, ".anno"))

	// fmt.Println("way id ", anno[0].Id)
	// fmt.Println("1017 ", geos[1017].Id)
	// fmt.Println("1016 ", geos[1016].Id)
	// fmt.Println("3 ", geos[3].Id)
	// fmt.Println("199 ", geos[199].Id)
	// fmt.Println("69 ", geos[69].Id)

	mapping := make(map[int32]int32)
	geometries := make([]graph.Geometry, 0)

	nodeCount := len(geos)
	graph := graph.NewOsmGraph(int32(nodeCount), uncomps)

	compressor := &Compressor{uncompGraph: graph, mappingGeoId: mapping, geometries: geometries}
	return compressor
}

func (compressor *Compressor) CompressGeometry(u, v, w int32, order int16) {

	uncompGraph := compressor.uncompGraph
	geometries := &compressor.geometries
	mapping := &compressor.mappingGeoId

	edge := uncompGraph.FindEdge(u, w)
	if edge == nil {
		panic("edge is nil")
	}

	var geoId int32
	geoId, exist := (*mapping)[edge.AnnoId]

	if !exist {
		geoId = int32(len(*geometries))
		(*mapping)[edge.AnnoId] = geoId

		*geometries = append(*geometries, graph.Geometry{})
	}

	nodes := &(*geometries)[geoId].Node
	orders := &(*geometries)[geoId].Order

	if nodes == nil {
		*nodes = make([]int32, 0)
		*orders = make([]int16, 0)
	}
	*nodes = append(*nodes, v)
	*orders = append(*orders, order)

	fmt.Print("ORDER ===> ")
	for _, order := range *orders {
		fmt.Print(int(order), " ")
	}
	fmt.Println()

	if lastid == -1 || lastid == edge.AnnoId {
		fmt.Println((*geometries)[geoId].Node, edge.AnnoId)
		lastid = edge.AnnoId
	}
}

func (compressor *Compressor) Compress() int {

	// restriction map
	var v int32

	graph := compressor.uncompGraph
	maxNode := graph.GetNodeCount()
	compressedNode := 0

	for v = 0; v < maxNode; v++ {
		if graph.GetOutDegree(v) != 2 {
			continue
		}

		// if _, exist := restricion_vias[v]; exist {
		// 	continue
		// }

		// forward way direction
		// from --------------------> to
		//
		// always (u->v).segIndex is less than (v-w).segIndex
		//
		// u <---------- v -----------> w
		//    ----------> <-----------
		//
		// Will be compressed to:
		//
		// u <---------- w
		//    ---------->
		//

		u := graph.GetTarget(v, 0)
		w := graph.GetTarget(v, 1)
		if u == -1 || w == -1 {
			panic("graph is invalid")
		}

		if graph.FindEdge(u, w) != nil {
			// u ~ w is directly connected
			continue
		}

		if graph.FindEdge(u, v).AnnoId != graph.FindEdge(v, w).AnnoId {
			// different osm way
			// not compress
			continue
		}

		if graph.FindEdge(u, v).AnnoId != 0 {
			continue
		}

		if graph.FindEdge(u, v).SegIndex > graph.FindEdge(v, w).SegIndex {
			u, w = w, u
		}

		// fmt.Println("Seg index u->v", graph.FindEdge(u, v).SegIndex)
		// fmt.Println("Seg index v->w", graph.FindEdge(v, w).SegIndex)
		// fmt.Println("fwd, rev", graph.FindEdge(v, w).Forward, graph.FindEdge(v, w).Reverse)
		// fmt.Println("u,v,w", u, v, w)

		fwd1, _ := graph.FindConstEdge(u, v)
		fwd2, _ := graph.FindConstEdge(v, w)
		rev1, _ := graph.FindConstEdge(w, v)
		rev2, _ := graph.FindConstEdge(v, u)

		// u -----> v <----- w
		graph.DelEdge(v, u)
		graph.DelEdge(v, w)

		// u ----------> w
		//   <----------
		graph.SetNewTarget(u, v, w)
		graph.SetNewTarget(w, v, u)

		graph.FindEdge(u, w).Distance = fwd1.Distance + fwd2.Distance
		graph.FindEdge(w, u).Distance = rev1.Distance + rev2.Distance
		graph.FindEdge(u, w).SegIndex = fwd2.SegIndex
		graph.FindEdge(w, u).SegIndex = fwd1.SegIndex

		compressor.CompressGeometry(u, v, w, fwd1.SegIndex)

		compressedNode++
	}

	// uncomp geometry
	return compressedNode
}
