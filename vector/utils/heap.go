package utils

import (
	c "hmruntime/vector/constraints"
)

type MinHeapElement[T c.Float] struct {
	value T
	index uint64
	// An element that is "filteredOut" is one that should be removed
	// from final consideration due to it not matching the passed in
	// filter criteria.
}

func InitHeapElement[T c.Float](
	val T, i uint64, filteredOut bool) *MinHeapElement[T] {
	return &MinHeapElement[T]{
		value: val,
		index: i,
	}
}

func (e MinHeapElement[T]) GetValue() T {
	return e.value
}

func (e MinHeapElement[T]) GetIndex() uint64 {
	return e.index
}

type MinTupleHeap[T c.Float] []MinHeapElement[T]

func (h MinTupleHeap[T]) Len() int {
	return len(h)
}

func (h MinTupleHeap[T]) Less(i, j int) bool {
	return h[i].value < h[j].value
}

func (h MinTupleHeap[T]) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinTupleHeap[T]) Push(x interface{}) {
	*h = append(*h, x.(MinHeapElement[T]))
}

func (h *MinTupleHeap[T]) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}
