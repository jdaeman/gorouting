package engine

import "math"

type PhantomNode struct {
	ForwardId  int32
	BackwardId int32

	FwdPosition int32

	X, Y float64
}

func (phantom PhantomNode) IsValidForward() bool {
	return phantom.ForwardId != math.MaxInt32
}

func (phantom PhantomNode) IsValidBackward() bool {
	return phantom.BackwardId != math.MaxInt32
}
