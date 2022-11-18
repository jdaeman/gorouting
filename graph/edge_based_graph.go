package graph

// GenerateEdgeExpandedNodes()

// GenerateEdgeExpandedEdges()

//

type EdgeBasedNode struct {
}

type EdgeBasedNodeSegment struct {
}

type EdgeBasedEdge struct {
	Source int32 // EdgeBasedNode id
	Target int32 // EdgeNasedNode id

	Distance int32

	Forward  bool
	Backward bool
}
