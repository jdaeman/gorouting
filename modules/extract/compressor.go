package extract

import (
	"modules/graph"
)

var lastid int32

type Compressor struct {
	uncompGraph *graph.NodeBasedGraph
}

func NewCompressor(rawFilePath string) *Compressor {

	//log.Println("Edge count", len(edges))

	{
		//anno := files.LoadEdgeAnnotations(files.ToDataPath(rawFilePath, ".anno"))
		// fmt.Println("way id ", anno[0].Id)
		// fmt.Println("1017 ", geos[1017].Id)
		// fmt.Println("1016 ", geos[1016].Id)
		// fmt.Println("3 ", geos[3].Id)
		// fmt.Println("199 ", geos[199].Id)
		// fmt.Println("69 ", geos[69].Id)
	}

	return nil
}

// func (compressor *Compressor) CompressGeometry(u, v, w int32, order int16) {

// 	uncompGraph := compressor.uncompGraph
// 	geometries := &compressor.geometries
// 	mapping := &compressor.mappingGeoId

// 	edge := uncompGraph.FindEdge(u, w)
// 	if edge == nil {
// 		panic("edge is nil")
// 	}

// 	var geoId int32
// 	geoId, exist := (*mapping)[edge.AnnoId]

// 	if !exist {
// 		geoId = int32(len(*geometries))
// 		(*mapping)[edge.AnnoId] = geoId

// 		*geometries = append(*geometries, graph.Geometry{})
// 	}

// 	nodes := &(*geometries)[geoId].Node
// 	orders := &(*geometries)[geoId].Order

// 	if nodes == nil {
// 		*nodes = make([]int32, 0)
// 		*orders = make([]int16, 0)
// 	}
// 	*nodes = append(*nodes, v)
// 	*orders = append(*orders, order)

// 	fmt.Print("ORDER ===> ")
// 	for _, order := range *orders {
// 		fmt.Print(int(order), " ")
// 	}
// 	fmt.Println()

// 	if lastid == -1 || lastid == edge.AnnoId {
// 		fmt.Println((*geometries)[geoId].Node, edge.AnnoId)
// 		lastid = edge.AnnoId
// 	}
// }

// func (compressor *Compressor) Compress() int {
// 	graph := compressor.uncompGraph

// 	number_of_nodes := graph.GetNumberOfNodes()

// 	for v := int32(0); v < number_of_nodes; v++ {
// 		if graph.GetOutDegree(v) != 2 {
// 			continue
// 		}
// 		//    reverse_e2   forward_e2
// 		// u <---------- v -----------> w
// 		//    ----------> <-----------
// 		//    forward_e1   reverse_e1
// 		//
// 		// Will be compressed to:
// 		//
// 		//    reverse_e1
// 		// u <---------- w
// 		//    ---------->
// 		//    forward_e1
// 		//
// 		// forward_e2, reverse_e2 will be removed
// 		//
// 		reverse_edge_order := int32(0)
// 		if graph.GetEdgeData(graph.BeginEdges(v)).Reverse {
// 			reverse_edge_order = 1
// 		}
// 		// reverse edge is oneway edge
// 		// in this case, begin edge is not forward
// 		//
// 		// always, u < v in bi-direction edge.
// 		// edge is sorted by source id.
// 		// so, v < w and w < u
// 		// v < w < u
// 		//
// 		forward_e2 := graph.BeginEdges(v) + reverse_edge_order
// 		reverse_e2 := graph.BeginEdges(v) + 1 - reverse_edge_order

// 		w, u := graph.GetTarget(forward_e2), graph.GetTarget(reverse_e2)
// 		// if graph.FindEdgeEitherDirection(u, w) != -1 {
// 		// 	continue
// 		// }

// 		forward_e1 := graph.FindEdge(u, v)
// 		reverse_e1 := graph.FindEdge(w, v)

// 		forward_distance2 := graph.GetEdgeData(forward_e2).Distance
// 		reverse_distance2 := graph.GetEdgeData(reverse_e2).Distance

// 		graph.GetEdgePtr(forward_e1).Distance += forward_distance2
// 		graph.GetEdgePtr(reverse_e1).Distance += reverse_distance2

// 		graph.SetTarget(forward_e1, w)
// 		graph.SetTarget(reverse_e1, u)

// 		graph.DeleteEdge(v, forward_e2)
// 		graph.DeleteEdge(v, reverse_e2)
// 	}

// 	return 0
// }

// func (compressor *Compressor) Compress() int {

// 	// restriction map
// 	var v int32

// 	graph := compressor.uncompGraph
// 	maxNode := graph.GetNodeCount()
// 	compressedNode := 0

// 	for v = 0; v < maxNode; v++ {
// 		if graph.GetOutDegree(v) != 2 {
// 			continue
// 		}

// 		// if _, exist := restricion_vias[v]; exist {
// 		// 	continue
// 		// }

// 		// forward way direction
// 		// from --------------------> to
// 		//
// 		// always (u->v).segIndex is less than (v-w).segIndex
// 		//
// 		// u <---------- v -----------> w
// 		//    ----------> <-----------
// 		//
// 		// Will be compressed to:
// 		//
// 		// u <---------- w
// 		//    ---------->
// 		//

// 		u := graph.GetTarget(v, 0)
// 		w := graph.GetTarget(v, 1)
// 		if u == -1 || w == -1 {
// 			panic("graph is invalid")
// 		}

// 		if graph.FindEdge(u, w) != nil {
// 			// u ~ w is directly connected
// 			continue
// 		}

// 		if graph.FindEdge(u, v).AnnoId != graph.FindEdge(v, w).AnnoId {
// 			// different osm way
// 			// not compress
// 			continue
// 		}

// 		if graph.FindEdge(u, v).AnnoId != 0 {
// 			continue
// 		}

// 		if graph.FindEdge(u, v).SegIndex > graph.FindEdge(v, w).SegIndex {
// 			u, w = w, u
// 		}

// 		// fmt.Println("Seg index u->v", graph.FindEdge(u, v).SegIndex)
// 		// fmt.Println("Seg index v->w", graph.FindEdge(v, w).SegIndex)
// 		// fmt.Println("fwd, rev", graph.FindEdge(v, w).Forward, graph.FindEdge(v, w).Reverse)
// 		// fmt.Println("u,v,w", u, v, w)

// 		fwd1, _ := graph.FindConstEdge(u, v)
// 		fwd2, _ := graph.FindConstEdge(v, w)
// 		rev1, _ := graph.FindConstEdge(w, v)
// 		rev2, _ := graph.FindConstEdge(v, u)

// 		// u -----> v <----- w
// 		graph.DelEdge(v, u)
// 		graph.DelEdge(v, w)

// 		// u ----------> w
// 		//   <----------
// 		graph.SetNewTarget(u, v, w)
// 		graph.SetNewTarget(w, v, u)

// 		graph.FindEdge(u, w).Distance = fwd1.Distance + fwd2.Distance
// 		graph.FindEdge(w, u).Distance = rev1.Distance + rev2.Distance
// 		graph.FindEdge(u, w).SegIndex = fwd2.SegIndex
// 		graph.FindEdge(w, u).SegIndex = fwd1.SegIndex

// 		compressor.CompressGeometry(u, v, w, fwd1.SegIndex)

// 		compressedNode++
// 	}

// 	// uncomp geometry
// 	return compressedNode
// }
