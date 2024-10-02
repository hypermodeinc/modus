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
	"fmt"
	"reflect"

	"github.com/chewxy/math32"
	"github.com/viterin/vek/vek32"
)

// DistanceFunc is a function that computes the distance between two vectors.
type DistanceFunc func(a, b []float32) (float32, error)

var (
	ErrDifferentVectorLengths = fmt.Errorf("vectors have different lengths")
)

// CosineDistance computes the cosine distance between two vectors.
func CosineDistance(a, b []float32) (float32, error) {
	if len(a) != len(b) {
		return 0, ErrDifferentVectorLengths
	}
	return 1 - vek32.CosineSimilarity(a, b), nil
}

// EuclideanDistance computes the Euclidean distance between two vectors.
func EuclideanDistance(a, b []float32) (float32, error) {
	if len(a) != len(b) {
		return 0, ErrDifferentVectorLengths
	}
	// TODO: can we speedup with vek?
	var sum float32 = 0
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return math32.Sqrt(sum), nil
}

var distanceFuncs = map[string]DistanceFunc{
	"euclidean": EuclideanDistance,
	"cosine":    CosineDistance,
}

func distanceFuncToName(fn DistanceFunc) (string, bool) {
	for name, f := range distanceFuncs {
		fnptr := reflect.ValueOf(fn).Pointer()
		fptr := reflect.ValueOf(f).Pointer()
		if fptr == fnptr {
			return name, true
		}
	}
	return "", false
}

// RegisterDistanceFunc registers a distance function with a name.
// A distance function must be registered here before a graph can be
// exported and imported.
func RegisterDistanceFunc(name string, fn DistanceFunc) {
	distanceFuncs[name] = fn
}
