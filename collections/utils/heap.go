package utils

type MinHeapElement struct {
	value float64
	index string
	// An element that is "filteredOut" is one that should be removed
	// from final consideration due to it not matching the passed in
	// filter criteria.
}

func InitHeapElement(
	val float64, i string, filteredOut bool) MinHeapElement {
	return MinHeapElement{
		value: val,
		index: i,
	}
}

func (e MinHeapElement) GetValue() float64 {
	return e.value
}

func (e MinHeapElement) GetIndex() string {
	return e.index
}

type MinTupleHeap []MinHeapElement

func (h MinTupleHeap) Len() int {
	return len(h)
}

func (h MinTupleHeap) Less(i, j int) bool {
	return h[j].value < h[i].value
}

func (h MinTupleHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinTupleHeap) Push(x interface{}) {
	*h = append(*h, x.(MinHeapElement))
}

func (h *MinTupleHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}
