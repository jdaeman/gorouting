package engine

type HeapNode struct {
	u      int32
	weight int32
}

type MinHeap []HeapNode

func (h MinHeap) Len() int {
	return len(h)
}

func (h MinHeap) Swap(parent, child int) {
	h[parent], h[child] = h[child], h[parent]
}

func (h MinHeap) Less(parent, child int) bool {
	return h[parent].weight < h[child].weight
}

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(HeapNode))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	elem := old[n-1]
	*h = old[:n-1]
	return elem
}
