package utils

import (
	"container/heap"
	"testing"
)

func TestHeap(t *testing.T) {
	h := &MinTupleHeap{}
	heap.Init(h)

	// Test Len before pushing any elements
	if h.Len() != 0 {
		t.Errorf("Expected length of 0, got %d", h.Len())
	}

	// Test Push
	heap.Push(h, MinHeapElement{value: 3.0, index: "three"})
	heap.Push(h, MinHeapElement{value: 1.0, index: "one"})
	heap.Push(h, MinHeapElement{value: 2.0, index: "two"})

	// Test Len
	if h.Len() != 3 {
		t.Errorf("Expected length of 3, got %d", h.Len())
	}

	// Test Less
	if !h.Less(0, 1) {
		t.Errorf("Expected h[0] < h[1], got h[0] = %v, h[1] = %v", (*h)[0], (*h)[1])
	}

	// Test Pop
	expectedValues := []float64{1.0, 2.0, 3.0}
	expectedIndices := []string{"one", "two", "three"}
	initialLen := h.Len() // Store initial length of heap

	for i := 0; i < initialLen; i++ {
		popped := heap.Pop(h).(MinHeapElement)
		if popped.value != expectedValues[i] || popped.index != expectedIndices[i] {
			t.Errorf("Expected pop value of %v and index '%s', got %v and '%s'", expectedValues[i], expectedIndices[i], popped.value, popped.index)
		}
	}

	// Test Len after popping all elements
	if h.Len() != 0 {
		t.Errorf("Expected length of 0, got %d", h.Len())
	}
}

func TestHeapSwap(t *testing.T) {
	h := &MinTupleHeap{
		MinHeapElement{value: 1.0, index: "one"},
		MinHeapElement{value: 2.0, index: "two"},
	}
	h.Swap(0, 1)

	if (*h)[0].value != 2.0 || (*h)[0].index != "two" || (*h)[1].value != 1.0 || (*h)[1].index != "one" {
		t.Errorf("Expected heap to be swapped, got %v", h)
	}
}
