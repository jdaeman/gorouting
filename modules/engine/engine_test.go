package engine_test

import (
	"fmt"
	"log"
	"modules/engine"
	"modules/files"
	"modules/view"
	"testing"

	"github.com/dhconnelly/rtreego"
)

type Thing struct {
	where   *rtreego.Rect
	phantom engine.PhantomNode
}

func (t *Thing) Bounds() *rtreego.Rect {
	return t.where
}

func TestRTree(t *testing.T) {
	// 2-dimension
	// branch factor: 25 ~ 50
	rt := rtreego.NewTree(2, 15, 30)

	dataReader := files.NewReader("../../extractor/data/south-korea-latest.osm.pbf")
	segments := dataReader.LoadEdgeBasedNodeSegements()
	geoNodes := dataReader.LoadGeoNodes()

	percent := make([]bool, 10)

	log.Println("RTREE START")
	for i, seg := range segments {
		a, b := seg.Forward_id, seg.Backward_id
		u, v := seg.U, seg.V
		pos := seg.Pos

		point1, point2 := geoNodes[u], geoNodes[v]
		rect, err := rtreego.NewRectFromPoints(rtreego.Point{point1.X, point1.Y}, rtreego.Point{point2.X, point2.Y})
		if err != nil {
			panic(err)
		}

		rt.Insert(&Thing{where: rect, phantom: engine.PhantomNode{ForwardId: a, BackwardId: b, FwdPosition: pos}})
		progress := int((float64(i+1) / float64(len(segments))) * 100.)
		progress %= 10
		if !percent[progress] {
			fmt.Printf("%d%%...", progress*10)
			percent[progress] = true
		}
	}
	log.Println("RTREE END")

	fmt.Println(rt.Size())
	matched := rt.NearestNeighbor(rtreego.Point{127.0016072, 37.5860800})
	//matched := rt.NearestNeighbor(rtreego.Point{127.7653880, 37.9038354})
	a := matched.(*Thing)
	fmt.Println(a.where)
	fmt.Println(a.where.Size())
	fmt.Println(a.where.PointCoord(0), a.where.PointCoord(1))
	fmt.Println(a.phantom)
	fmt.Println(a.phantom.IsValidForward())
	fmt.Println(a.phantom.IsValidBackward())

	matched = rt.NearestNeighbor(rtreego.Point{127.7653880, 37.9038354})
	a = matched.(*Thing)
	fmt.Println(a.phantom)
	fmt.Println(a.where.PointCoord(0), a.where.PointCoord(1))
}

func TestDFS(t *testing.T) {
	config := engine.NewEngineConfig("../../extractor/data/south-korea-latest.osm.pbf")
	eng := engine.NewRoutingEngine(*config)

	order := eng.Dfs(3757)

	viewFactory := view.NewViewFactory("../../extractor/data/south-korea-latest.osm.pbf")
	viewFactory.DrawingEdges(order, -1)
}

func TestSimpleSearch(t *testing.T) {
	config := engine.NewEngineConfig("../../extractor/data/south-korea-latest.osm.pbf")
	eng := engine.NewRoutingEngine(*config)

	order := eng.SimplePathSearch(0, 2000)
	fmt.Println(order)
	viewFactory := view.NewViewFactory("../../extractor/data/south-korea-latest.osm.pbf")
	viewFactory.DrawingEdges(order, -1)
}

func TestBidirSearch(t *testing.T) {
	config := engine.NewEngineConfig("../../extractor/data/south-korea-latest.osm.pbf")
	eng := engine.NewRoutingEngine(*config)

	// order := eng.SimplePathSearch(0, 2000)
	// fmt.Println(order)

	//eng.ShortestPathSearch([]int32{0, 1}, []int32{4557, 4456})
	order := eng.ShortestPathSearch([]int32{0}, []int32{2000})
	fmt.Println(order)
}
