package engine

import (
	"modules/graph"

	"github.com/dhconnelly/rtreego"
)

type rt_seg struct {
	fwd, rev int32
	fwdPos   int32
}

type rt_node struct {
	where *rtreego.Rect
	seg   rt_seg
}

func (t *rt_node) Bounds() *rtreego.Rect {
	return t.where
}

type RTree struct {
	tree *rtreego.Rtree
}

// Constructor
func NewRTree(nodes []graph.ResultNode, segments []graph.EdgeBasedNodeSegment) *RTree {
	if nodes == nil || segments == nil {
		return nil
	}

	rt := rtreego.NewTree(2, 25, 50)

	for _, seg := range segments {
		fwd, rev := seg.Forward_id, seg.Backward_id
		coord_u, coord_v := seg.U, seg.V
		fwdPos := seg.Pos

		point1, point2 := nodes[coord_u], nodes[coord_v]
		rect, err := rtreego.NewRectFromPoints(
			rtreego.Point{point1.X, point1.Y},
			rtreego.Point{point2.X, point2.Y},
		)
		if err != nil {
			panic(err)
		}

		rt.Insert(&rt_node{
			where: rect,
			seg: rt_seg{
				fwd: fwd, rev: rev,
				fwdPos: fwdPos}},
		)
	}

	ret := &RTree{tree: rt}
	return ret
}

// requested coordinates to map-matched data structure
func (rtree *RTree) GetPhantomNodes(coords [][2]float64) []PhantomNode {
	rt := rtree.tree

	ret := make([]PhantomNode, 0, len(coords))
	for _, coord := range coords {
		x, y := coord[0], coord[1]

		matched := rt.NearestNeighbor(rtreego.Point{x, y})
		if matched == nil {
			break
		}

		node := matched.(*rt_node)
		rect, seg := node.where, node.seg
		mx, my := rect.PointCoord(0), rect.PointCoord(1)

		// dist := util.HaversineDistance([2]float64{x, y}, [2]float64{mx, my})
		// if dist > 1000 {
		// 	// too far distance...
		// 	break
		// }

		phantom := PhantomNode{
			ForwardId:   seg.fwd,
			BackwardId:  seg.rev,
			FwdPosition: seg.fwdPos,
			X:           mx,
			Y:           my,
		}

		ret = append(ret, phantom)
	}

	return ret
}
