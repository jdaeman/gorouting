package engine_test

import (
	"engine"
	"files"
	"fmt"
	"testing"
	"view"

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
		fmt.Println(point1, point2)
		rect, err := rtreego.NewRectFromPoints(rtreego.Point{point1.X, point1.Y}, rtreego.Point{point2.X, point2.Y})
		if err != nil {
			panic(err)
		}

		rt.Insert(&Thing{where: rect, phantom: engine.PhantomNode{ForwardId: a, BackwardId: b, Position: pos}})
	}

	fmt.Println(rt.Size())
	matched := rt.NearestNeighbor(rtreego.Point{127.0016072, 37.5860800})
	a := matched.(*Thing)
	fmt.Println(a.where)
	fmt.Println(a.phantom)
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

	// order := eng.SimplePathSearch(0, 4557)
	// fmt.Println(order)

	eng.ShortestPathSearch(0, 4557)
}
