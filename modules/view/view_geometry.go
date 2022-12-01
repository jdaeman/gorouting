package view

import (
	"fmt"
	"modules/files"
	"modules/graph"
	"strconv"

	"github.com/jdaeman/go-shp"
)

type ViewFactory struct {
	edge_based_nodes []graph.EdgeBasedNode
	geometires       []graph.Geometry
	geoNodes         []graph.ResultNode
	annotations      []graph.EdgeAnnotation
}

func NewViewFactory(filepath string) *ViewFactory {
	annotations := files.LoadEdgeAnnotations(files.ToDataPath(filepath, files.ANNOTATION))
	edge_based_nodes := files.LoadEdgeBasedNodes(files.ToDataPath(filepath, files.EBNODE))
	geometries := files.LoadEdgeGeometries(files.ToDataPath(filepath, files.GEOMETRY))
	geoNodes := files.LoadGeoNodes(files.ToDataPath(filepath, files.GEONODE))

	ret := &ViewFactory{
		edge_based_nodes: edge_based_nodes,
		geometires:       geometries,
		geoNodes:         geoNodes,
		annotations:      annotations,
	}

	return ret
}

func (v ViewFactory) get(geoId uint32) []int32 {
	forward := true
	if geoId&0x80000000 > 0 {
		forward = false
	}

	geoId &= 0x7fffffff
	var geos []int32
	{
		org := v.geometires[geoId].Nodes
		geos = make([]int32, len(org))
		if copy(geos, org) != len(org) {
			panic("Copy fail")
		}
	}

	if !forward {
		f, t := 0, len(geos)-1
		fmt.Println("REVERSE", f, t)
		for f < t {
			geos[f], geos[t] = geos[t], geos[f]
			f, t = f+1, t-1
		}
	}

	return geos
}

func (v ViewFactory) DrawingEdges(path []int32, lastSeg int32) {
	w, err := shp.Create("1_path", shp.POLYLINE)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	polylines := make([]shp.PolyLine, 0, len(path))
	for i, u := range path {
		geoId := v.edge_based_nodes[u].GeometryId
		annoId := v.edge_based_nodes[u].AnnotationId
		wayid := v.annotations[annoId].Id
		geos := v.get(geoId)
		fmt.Println(i+1, "Geoid", geoId, geoId&0x7fffffff, wayid, geos)

		points := make([]shp.Point, 0, len(geos))

		for l, id := range geos {
			if i == len(path)-1 && l == int(lastSeg) {
				break
			}
			x, y := v.geoNodes[id].X, v.geoNodes[id].Y
			points = append(points, shp.Point{X: x, Y: y})
		}
		polylines = append(polylines, *shp.NewPolyLine([][]shp.Point{points}))
	}

	for i := range polylines {
		w.Write(&polylines[i])
	}
}

func (v ViewFactory) DrawingEdge(u int32) {
	geometryId := v.edge_based_nodes[u].GeometryId
	annoId := v.edge_based_nodes[u].AnnotationId

	fmt.Println(geometryId, v.annotations[annoId].Id)

	geometryId &= 0x7fffffff
	nodes := v.geometires[geometryId].Nodes

	points := make([]shp.Point, 0, len(nodes))
	for _, node := range nodes {
		points = append(points, shp.Point{X: v.geoNodes[node].X, Y: v.geoNodes[node].Y})
	}

	w, err := shp.Create("1_"+strconv.Itoa(int(u)), shp.POLYLINE)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	polyline := shp.NewPolyLine([][]shp.Point{
		points,
	})

	w.Write(polyline)
}
