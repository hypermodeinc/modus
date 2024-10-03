/*
 * The code in this file originates from https://github.com/coder/hnsw
 * and is licensed under the terms of the Creative Commons Zero v1.0 Universal license
 * See the LICENSE file in the "hnsw" directory that accompanied this code for further details.
 * See also: https://github.com/coder/hnsw/blob/main/LICENSE
 *
 * SPDX-FileCopyrightText: Ammar Bandukwala <ammar@ammar.io>
 * SPDX-License-Identifier: CC0-1.0
 */

package heap

import (
	"math/rand"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

type Int int

func (i Int) Less(j Int) bool {
	return i < j
}

func TestHeap(t *testing.T) {
	h := Heap[Int]{}

	for i := 0; i < 20; i++ {
		h.Push(Int(rand.Int() % 100))
	}

	require.Equal(t, 20, h.Len())

	var inOrder []Int
	for h.Len() > 0 {
		inOrder = append(inOrder, h.Pop())
	}

	if !slices.IsSorted(inOrder) {
		t.Errorf("Heap did not return sorted elements: %+v", inOrder)
	}
}
