package files

import (
	"modules/graph"
)

type DataRepository struct {
	geoNodes    []graph.ResultNode
	nodes       []graph.EdgeBasedNode
	edges       []graph.EdgeBasedEdge
	segments    []graph.EdgeBasedNodeSegment
	annotations []graph.EdgeAnnotation
	geometries  []graph.Geometry
}

func NewDataRepository(datapath string) *DataRepository {
	reader := NewReader(datapath)

	geoNodes := reader.LoadGeoNodes()
	nodes := reader.LoadEdgeBasedNodes()
	edges := reader.LoadEdgeBasedEdges()
	segments := reader.LoadEdgeBasedNodeSegements()
	annotations := reader.LoadEdgeAnnotations()
	geometries := reader.LoadEdgeGeometries()

	ret := &DataRepository{
		geoNodes:    geoNodes,
		nodes:       nodes,
		edges:       edges,
		segments:    segments,
		annotations: annotations,
		geometries:  geometries,
	}

	return ret
}

func (repo *DataRepository) GetGeoNodes() []graph.ResultNode {
	return repo.geoNodes
}

func (repo *DataRepository) GetNodes() []graph.EdgeBasedNode {
	return repo.nodes
}

func (repo *DataRepository) GetEdges() []graph.EdgeBasedEdge {
	return repo.edges
}

func (repo *DataRepository) GetSegments() []graph.EdgeBasedNodeSegment {
	return repo.segments
}

// func (repo *DataRepository) GetAnnotations() []graph.EdgeAnnotation {
// 	return repo.annotations
// }

// func (repo *DataRepository) GetGeometries() []graph.Geometry {
// 	return repo.geometries
// }

func (repo *DataRepository) GetGeometryId(u int32) uint32 {
	return repo.nodes[u].GeometryId
}

func (repo *DataRepository) GetAnnotationId(u int32) int32 {
	return repo.nodes[u].AnnotationId
}

func (repo *DataRepository) GetAnnotation(anno int32) graph.EdgeAnnotation {
	return repo.annotations[anno]
}

func (repo *DataRepository) GetGeometry(geoId uint32) ([]int32, []int32) {
	forward := true
	if geoId&0x80000000 > 0 {
		forward = false
	}

	geoId &= 0x7fffffff
	var geos, distances []int32
	// copy
	{
		orgGeos := repo.geometries[geoId].Nodes
		orgDists := repo.geometries[geoId].Distances

		if len(orgGeos) != len(orgDists) {
			panic("Geometry data invalid")
		}

		geos = make([]int32, len(orgGeos))
		distances = make([]int32, len(orgDists))

		copyLen := copy(geos, orgGeos)
		copyLen += copy(distances, orgDists)
		if copyLen < len(orgGeos)*2 {
			panic("Copy fail")
		}
	}

	if !forward {
		f, t := 0, len(geos)-1
		for f < t {
			geos[f], geos[t] = geos[t], geos[f]
			distances[f], distances[t] = distances[t], distances[f]
			f, t = f+1, t-1
		}
	}

	return geos, distances
}

// TODO
//
func (repo *DataRepository) GetGeoNodeIds(geoId uint32) []int32 {
	forward := true
	if geoId&0x80000000 > 0 {
		forward = false
	}

	geoId &= 0x7fffffff
	var geos []int32
	{
		org := repo.geometries[geoId].Nodes
		geos = make([]int32, len(org))
		if copy(geos, org) != len(org) {
			panic("Copy fail")
		}
	}

	if !forward {
		f, t := 0, len(geos)-1
		for f < t {
			geos[f], geos[t] = geos[t], geos[f]
			f, t = f+1, t-1
		}
	}
	return geos
}

// TODO
//
func (repo *DataRepository) GetGeometryDistance(geoId uint32) []int32 {
	forward := true
	if geoId&0x80000000 > 0 {
		forward = false
	}

	geoId &= 0x7fffffff
	var distances []int32
	{
		dist := repo.geometries[geoId].Distances
		distances = make([]int32, len(dist))
		if copy(distances, dist) != len(dist) {
			panic("Copy fail")
		}
	}

	if !forward {
		f, t := 0, len(distances)-1
		for f < t {
			distances[f], distances[t] = distances[t], distances[f]
			f, t = f+1, t-1
		}
	}

	return distances
}

func (repo *DataRepository) GetLocations(ids []int32) [][]float64 {
	geoNodes := repo.geoNodes

	ret := make([][]float64, 0, len(ids))

	for _, id := range ids {
		x, y := geoNodes[id].X, geoNodes[id].Y
		ret = append(ret, []float64{x, y})
	}

	return ret
}

func (repo *DataRepository) GetLocationsYX(ids []int32) [][]float64 {
	geoNodes := repo.geoNodes

	ret := make([][]float64, 0, len(ids))

	for _, id := range ids {
		x, y := geoNodes[id].X, geoNodes[id].Y
		ret = append(ret, []float64{y, x})
	}

	return ret
}
