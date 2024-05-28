package utils

import (
	"errors"
	"math"

	c "hmruntime/vector/constraints"

	"github.com/chewxy/math32"
)

const (
	Euclidian            = "euclidian"
	Cosine               = "cosine"
	DotProd              = "dotproduct"
	plError              = "\nerror fetching posting list for data key: "
	dataError            = "\nerror fetching data for data key: "
	VecKeyword           = "__vector_"
	visitedVectorsLevel  = "visited_vectors_level_"
	distanceComputations = "vector_distance_computations"
	searchTime           = "vector_search_time"
	VecEntry             = "__vector_entry"
	VecDead              = "__vector_dead"
	VectorIndexMaxLevels = 5
	EfConstruction       = 16
	EfSearch             = 12
	numEdgesConst        = 2
	// ByteData indicates the key stores data.
	ByteData = byte(0x00)
	// DefaultPrefix is the prefix used for data, index and reverse keys so that relative
	DefaultPrefix = byte(0x00)
	// NsSeparator is the separator between the namespace and attribute.
	NsSeparator = "-"
)

func IsBetterScoreForDistance[T c.Float](a, b T) bool {
	return a < b
}

func IsBetterScoreForSimilarity[T c.Float](a, b T) bool {
	return a > b
}

func norm[T c.Float](v []T, floatBits int) T {
	vectorNorm, _ := DotProduct(v, v, floatBits)
	if floatBits == 32 {
		return T(math32.Sqrt(float32(vectorNorm)))
	}
	if floatBits == 64 {
		return T(math.Sqrt(float64(vectorNorm)))
	}
	panic("Invalid floatBits")
}

// This needs to implement signature of SimilarityType[T].distanceScore
// function, hence it takes in a floatBits parameter,
// but doesn't actually use it.
func DotProduct[T c.Float](a, b []T, floatBits int) (T, error) {
	var dotProduct T
	if len(a) != len(b) {
		err := errors.New("can not compute dot product on vectors of different lengths")
		return dotProduct, err
	}
	for i := range a {
		dotProduct += a[i] * b[i]
	}
	return dotProduct, nil
}

// This needs to implement signature of SimilarityType[T].distanceScore
// function, hence it takes in a floatBits parameter.
func CosineSimilarity[T c.Float](a, b []T, floatBits int) (T, error) {
	dotProd, err := DotProduct(a, b, floatBits)
	if err != nil {
		return 0, err
	}
	normA := norm[T](a, floatBits)
	normB := norm[T](b, floatBits)
	if normA == 0 || normB == 0 {
		err := errors.New("can not compute cosine similarity on zero vector")
		var empty T
		return empty, err
	}
	return dotProd / (normA * normB), nil
}

// This needs to implement signature of SimilarityType[T].distanceScore
// function, hence it takes in a floatBits parameter,
// but doesn't actually use it.
func EuclidianDistanceSq[T c.Float](a, b []T, floatBits int) (T, error) {
	if len(a) != len(b) {
		return 0, errors.New("can not subtract vectors of different lengths")
	}
	var distSq T
	for i := range a {
		val := a[i] - b[i]
		distSq += val * val
	}
	return distSq, nil
}

func ConcatStrings(strs ...string) string {
	total := ""
	for _, s := range strs {
		total += s
	}
	return total
}
