/*
 * The code in this file originates from https://github.com/coder/hnsw
 * and is licensed under the terms of the Creative Commons Zero v1.0 Universal license
 * See the LICENSE file in the "hnsw" directory that accompanied this code for further details.
 * See also: https://github.com/coder/hnsw/blob/main/LICENSE
 *
 * SPDX-FileCopyrightText: Ammar Bandukwala <ammar@ammar.io>
 * SPDX-License-Identifier: CC0-1.0
 */

package hnsw

import (
	"cmp"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_maxLevel(t *testing.T) {
	var m int

	m, _ = maxLevel(0.5, 10)
	require.Equal(t, 4, m)

	m, _ = maxLevel(0.5, 1000)
	require.Equal(t, 11, m)
}

func Test_layerNode_search(t *testing.T) {
	entry := &layerNode[int]{
		Node: Node[int]{
			Value: Vector{0},
			Key:   0,
		},
		neighbors: map[int]*layerNeighborNode[int]{
			1: {
				node: &layerNode[int]{
					Node: Node[int]{
						Value: Vector{1},
						Key:   1,
					},
				},
				distance: 1.0,
			},
			2: {
				node: &layerNode[int]{
					Node: Node[int]{
						Value: Vector{2},
						Key:   2,
					},
				},
				distance: 2.0,
			},
			3: {
				node: &layerNode[int]{
					Node: Node[int]{
						Value: Vector{3},
						Key:   3,
					},
					neighbors: map[int]*layerNeighborNode[int]{
						4: {
							node: &layerNode[int]{
								Node: Node[int]{
									Value: Vector{4},
									Key:   4,
								},
							},
							distance: 4.0,
						},
						5: {
							node: &layerNode[int]{
								Node: Node[int]{
									Value: Vector{5},
									Key:   5,
								},
							},
							distance: 5.0,
						},
					},
				},
				distance: 3.0,
			},
		},
	}

	best, _ := entry.search(2, 4, []float32{4}, EuclideanDistance)

	require.Equal(t, 4, best[0].node.Key)
	require.Equal(t, 3, best[1].node.Key)
	require.Len(t, best, 2)
}

func newTestGraph[K cmp.Ordered]() *Graph[K] {
	return &Graph[K]{
		M:        6,
		Distance: EuclideanDistance,
		Ml:       0.5,
		EfSearch: 20,
		Rng:      rand.New(rand.NewSource(0)),
	}
}

func TestGraph_AddSearch(t *testing.T) {
	g := newTestGraph[int]()

	for i := range 128 {
		err := g.Add(
			Node[int]{
				Key:   i,
				Value: Vector{float32(i)},
			},
		)
		require.NoError(t, err)
	}

	al := Analyzer[int]{Graph: g}

	// Layers should be approximately log2(128) = 7
	// Look for an approximate doubling of the number of nodes in each layer.
	require.Equal(t, []int{
		128,
		67,
		28,
		12,
		6,
		2,
		1,
		1,
	}, al.Topography())

	nearest, _ := g.Search(
		[]float32{64.5},
		4,
	)

	require.Len(t, nearest, 4)
	require.EqualValues(
		t,
		[]SearchResultNode[int]{
			{Node: Node[int]{Key: 64, Value: Vector{64}}, Distance: 0.5},
			{Node: Node[int]{Key: 65, Value: Vector{65}}, Distance: 0.5},
			{Node: Node[int]{Key: 62, Value: Vector{62}}, Distance: 2.5},
			{Node: Node[int]{Key: 63, Value: Vector{63}}, Distance: 1.5},
		},
		nearest,
	)
}

func TestGraph_AddDelete(t *testing.T) {
	g := newTestGraph[int]()
	for i := range 128 {
		err := g.Add(Node[int]{
			Key:   i,
			Value: Vector{float32(i)},
		})
		require.NoError(t, err)
	}

	require.Equal(t, 128, g.Len())
	an := Analyzer[int]{Graph: g}

	preDeleteConnectivity := an.Connectivity()

	// Delete every even node.
	for i := 0; i < 128; i += 2 {
		ok := g.Delete(i)
		require.True(t, ok)
	}

	require.Equal(t, 64, g.Len())

	postDeleteConnectivity := an.Connectivity()

	// Connectivity should be the same for the lowest layer.
	require.Equal(
		t, preDeleteConnectivity[0],
		postDeleteConnectivity[0],
	)

	t.Run("DeleteNotFound", func(t *testing.T) {
		ok := g.Delete(-1)
		require.False(t, ok)
	})
}

func Benchmark_HNSW(b *testing.B) {
	b.ReportAllocs()

	sizes := []int{100, 1000, 10000}

	// Use this to ensure that complexity is O(log n) where n = h.Len().
	for _, size := range sizes {
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			g := Graph[int]{}
			g.Ml = 0.5
			g.Distance = EuclideanDistance
			for i := range size {
				err := g.Add(Node[int]{
					Key:   i,
					Value: Vector{float32(i)},
				})
				require.NoError(b, err)
			}
			b.ResetTimer()

			b.Run("Search", func(b *testing.B) {
				for i := 0; b.Loop(); i++ {
					_, err := g.Search(
						[]float32{float32(i % size)},
						4,
					)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

func randFloats(n int) []float32 {
	x := make([]float32, n)
	for i := range x {
		x[i] = rand.Float32()
	}
	return x
}

func Benchmark_HNSW_1536(b *testing.B) {
	b.ReportAllocs()

	g := newTestGraph[int]()
	const size = 1000
	points := make([]Node[int], size)
	for i := range size {
		points[i] = Node[int]{
			Key:   i,
			Value: Vector(randFloats(1536)),
		}
		err := g.Add(points[i])
		require.NoError(b, err)
	}
	b.ResetTimer()

	b.Run("Search", func(b *testing.B) {
		for i := 0; b.Loop(); i++ {
			_, err := g.Search(
				points[i%size].Value,
				4,
			)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func TestGraph_DefaultCosine(t *testing.T) {
	g := NewGraph[int]()
	err := g.Add(
		Node[int]{Key: 1, Value: Vector{1, 1}},
		Node[int]{Key: 2, Value: Vector{0, 1}},
		Node[int]{Key: 3, Value: Vector{1, -1}},
	)
	require.NoError(t, err)

	neighbors, _ := g.Search(
		[]float32{0.5, 0.5},
		1,
	)

	require.Equal(
		t,
		[]SearchResultNode[int]{
			{Node: Node[int]{Key: 1, Value: Vector{1, 1}}, Distance: 0},
		},
		neighbors,
	)
}
