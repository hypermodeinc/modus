package utils

type MaxHeapElement struct {
	value float64
	index string
	// An element that is "filteredOut" is one that should be removed
	// from final consideration due to it not matching the passed in
	// filter criteria.
}

func InitHeapElement(
	val float64, i string, filteredOut bool) MaxHeapElement {
	return MaxHeapElement{
		value: val,
		index: i,
	}
}

func (e MaxHeapElement) GetValue() float64 {
	return e.value
}

func (e MaxHeapElement) GetIndex() string {
	return e.index
}

type MaxTupleHeap []MaxHeapElement

func (h MaxTupleHeap) Len() int {
	return len(h)
}

func (h MaxTupleHeap) Less(i, j int) bool {
	return h[j].value < h[i].value
}

func (h MaxTupleHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MaxTupleHeap) Push(x interface{}) {
	*h = append(*h, x.(MaxHeapElement))
}

func (h *MaxTupleHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}
