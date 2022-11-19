package graph

// GenerateEdgeExpandedNodes()

// GenerateEdgeExpandedEdges()

//

type EdgeBasedNode struct {
	GeometryId   uint32
	AnnotationId int32
}

type EdgeBasedNodeSegment struct {
	Forward_id  int32 // EdgeBasedNode id
	Backward_id int32 // EdgeBasedNode id

	U, V int32 // node id, for location
}

type EdgeBasedEdge struct {
	Source int32 // EdgeBasedNode id
	Target int32 // EdgeNasedNode id

	Distance int32

	Forward  bool
	Backward bool
}
