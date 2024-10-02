/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package utils

import (
	"errors"
	"fmt"
	"math"

	"github.com/chewxy/math32"
	"github.com/viterin/vek/vek32"
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

func IsBetterScoreForDistance(a, b float64) bool {
	return a < b
}

func Normalize(v []float32) ([]float32, error) {
	norm := norm(v)
	if norm == 0 {
		return nil, errors.New("can not normalize vector with zero norm")
	}
	v = vek32.DivNumber(v, norm)
	return v, nil
}

func norm(v []float32) float32 {
	vectorNorm, _ := DotProduct(v, v)
	return math32.Sqrt(vectorNorm)
}

func DotProduct(a, b []float32) (float32, error) {
	if len(a) != len(b) {
		return 0, errors.New("can not compute dot product on vectors of different lengths")
	}
	return vek32.Dot(a, b), nil
}

// assume normalization for vectors
func CosineDistance(a, b []float32) (float64, error) {
	dotProd, err := DotProduct(a, b)
	if err != nil {
		return 0, err
	}

	return 1 - float64(dotProd), nil
}

func ConcatStrings(strs ...string) string {
	total := ""
	for _, s := range strs {
		total += s
	}
	return total
}

func ConvertToFloat32_2DArray(vecs any) ([][]float32, error) {

	// return if already a slice of float32 slices
	if vec32s, ok := vecs.([][]float32); ok {
		return vec32s, nil
	}

	// convert [][]float64 to [][]float32
	if vec64s, ok := vecs.([][]float64); ok {
		vec32s := make([][]float32, len(vec64s))
		for i, vec := range vec64s {
			vec32s[i] = make([]float32, len(vec))
			for j, val := range vec {
				vec32s[i][j] = float32(val)
			}
		}
		return vec32s, nil
	}

	return nil, fmt.Errorf("expected a slice of float32 or float64 slices, got: %T", vecs)
}

func EqualFloat32Slices(a, b []float32) bool {
	const epsilon = 1e-9
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if math.Abs(float64(a[i]-b[i])) > epsilon {
			return false
		}
	}
	return true
}
