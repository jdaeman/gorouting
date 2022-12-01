package engine_test

import (
	"fmt"
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
	rt := rtreego.NewTree(2, 25, 50)

	dataReader := files.NewReader("../extract/data/map.osm")
	segments := dataReader.LoadEdgeBasedNodeSegements()
	geoNodes := dataReader.LoadGeoNodes()

	for _, seg := range segments {
		a, b := seg.Forward_id, seg.Backward_id
		u, v := seg.U, seg.V
		pos := seg.Pos

		point1, point2 := geoNodes[u], geoNodes[v]
		rect, err := rtreego.NewRectFromPoints(rtreego.Point{point1.X, point1.Y}, rtreego.Point{point2.X, point2.Y})
		if err != nil {
			panic(err)
		}

		rt.Insert(&Thing{where: rect, phantom: engine.PhantomNode{ForwardId: a, BackwardId: b, FwdPosition: pos}})
	}

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
}

func TestDFS(t *testing.T) {
	config := engine.NewEngineConfig("../extract/data/map.osm")
	eng := engine.NewRoutingEngine(*config)

	order := eng.Dfs(0)

	viewFactory := view.NewViewFactory("../extract/data/map.osm")
	viewFactory.DrawingEdges(order, -1)
}

func TestSimpleSearch(t *testing.T) {
	config := engine.NewEngineConfig("../extract/data/map.osm")
	eng := engine.NewRoutingEngine(*config)

	order := eng.SimplePathSearch(0, 4557)
	viewFactory := view.NewViewFactory("../extract/data/map.osm")
	viewFactory.DrawingEdges(order, 8)
}

func TestBidirSearch(t *testing.T) {
	config := engine.NewEngineConfig("../extract/data/map.osm")
	eng := engine.NewRoutingEngine(*config)

	order := eng.SimplePathSearch(0, 0)
	fmt.Println(order)

	eng.ShortestPathSearch([]int32{0, 1}, []int32{4557, 4456})
}
