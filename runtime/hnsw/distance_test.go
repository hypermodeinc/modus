package hnsw

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEuclideanDistance(t *testing.T) {
	a := []float32{1, 2, 3}
	b := []float32{4, 5, 6}
	expected := float32(5.196152)
	actual, _ := EuclideanDistance(a, b)
	require.Equal(t, expected, actual)
}

func TestCosineSimilarity(t *testing.T) {
	var a, b []float32
	// Same magnitude, same direction.
	a = []float32{1, 1, 1}
	b = []float32{0.8, 0.8, 0.8}
	distance, _ := CosineDistance(a, b)
	require.InDelta(t, 0, distance, 0.000001)

	// Perpendicular vectors.
	a = []float32{1, 0}
	b = []float32{0, 1}
	distance, _ = CosineDistance(a, b)
	require.InDelta(t, 1, distance, 0.000001)

	// Equivalent vectors.
	a = []float32{1, 0}
	b = []float32{1, 0}
	distance, _ = CosineDistance(a, b)
	require.InDelta(t, 0, distance, 0.000001)
}

func BenchmarkCosineSimilarity(b *testing.B) {
	v1 := randFloats(1536)
	v2 := randFloats(1536)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := CosineDistance(v1, v2)
		if err != nil {
			b.Fatal(err)
		}
	}
}
