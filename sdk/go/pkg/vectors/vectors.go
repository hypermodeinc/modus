/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package vectors

import (
	"math"

	"golang.org/x/exp/constraints"
)

// Add adds two vectors together and returns the result.
func Add[T constraints.Integer | constraints.Float](a, b []T) []T {
	assertEqualLength(a, b)
	result := make([]T, len(a))
	for i := range a {
		result[i] = a[i] + b[i]
	}
	return result
}

// AddInPlace adds two vectors together and stores the result in the first vector.
func AddInPlace[T constraints.Integer | constraints.Float](a, b []T) {
	assertEqualLength(a, b)
	for i := range a {
		a[i] += b[i]
	}
}

// Subtract subtracts one vector from another and returns the result.
func Subtract[T constraints.Integer | constraints.Float](a, b []T) []T {
	assertEqualLength(a, b)
	result := make([]T, len(a))
	for i := range a {
		result[i] = a[i] - b[i]
	}
	return result
}

// SubtractInPlace subtracts one vector from another and stores the result in the first vector.
func SubtractInPlace[T constraints.Integer | constraints.Float](a, b []T) {
	assertEqualLength(a, b)
	for i := range a {
		a[i] -= b[i]
	}
}

// AddNumber adds a number to each element of a vector and returns the result.
func AddNumber[T constraints.Integer | constraints.Float](a []T, b T) []T {
	result := make([]T, len(a))
	for i := range a {
		result[i] = a[i] + b
	}
	return result
}

// AddNumberInPlace adds a number to each element of a vector and stores the result in the vector.
func AddNumberInPlace[T constraints.Integer | constraints.Float](a []T, b T) {
	for i := range a {
		a[i] += b
	}
}

// SubtractNumber subtracts a number from each element of a vector and returns the result.
func SubtractNumber[T constraints.Integer | constraints.Float](a []T, b T) []T {
	result := make([]T, len(a))
	for i := range a {
		result[i] = a[i] - b
	}
	return result
}

// SubtractNumberInPlace subtracts a number from each element of a vector and stores the result in the vector.
func SubtractNumberInPlace[T constraints.Integer | constraints.Float](a []T, b T) {
	for i := range a {
		a[i] -= b
	}
}

// MultiplyNumber multiplies each element of a vector by a number and returns the result.
func MultiplyNumber[T constraints.Integer | constraints.Float](a []T, b T) []T {
	result := make([]T, len(a))
	for i := range a {
		result[i] = a[i] * b
	}
	return result
}

// MultiplyNumberInPlace multiplies each element of a vector by a number and stores the result in the vector.
func MultiplyNumberInPlace[T constraints.Integer | constraints.Float](a []T, b T) {
	for i := range a {
		a[i] *= b
	}
}

// DivideNumber divides each element of a vector by a number and returns the result.
func DivideNumber[T constraints.Integer | constraints.Float](a []T, b T) []T {
	assertNotZero(b)
	result := make([]T, len(a))
	for i := range a {
		result[i] = a[i] / b
	}
	return result
}

// DivideNumberInPlace divides each element of a vector by a number and stores the result in the vector.
func DivideNumberInPlace[T constraints.Integer | constraints.Float](a []T, b T) {
	assertNotZero(b)
	for i := range a {
		a[i] /= b
	}
}

// Dot computes the dot product of two vectors.
func Dot[T constraints.Integer | constraints.Float](a, b []T) T {
	assertEqualLength(a, b)
	var result T = 0
	for i := 0; i < len(a); i++ {
		result += a[i] * b[i]
	}
	return result
}

// Magnitude computes the magnitude of a vector.
func Magnitude[T constraints.Integer | constraints.Float](a []T) T {
	return T(math.Sqrt(float64(Dot(a, a))))
}

// Normalize normalizes a vector to have a magnitude of 1.
func Normalize[T constraints.Integer | constraints.Float](a []T) []T {
	mag := Magnitude(a)
	return DivideNumber(a, mag)
}

// Sum computes the sum of all elements in a vector.
func Sum[T constraints.Integer | constraints.Float](a []T) T {
	var result T = 0
	for i := range a {
		result += a[i]
	}
	return result
}

// Product computes the product of all elements in a vector.
func Product[T constraints.Integer | constraints.Float](a []T) T {
	var result T = 1
	for i := 0; i < len(a); i++ {
		result *= a[i]
	}
	return result
}

// func Mean computes the mean of a vector.
func Mean[T constraints.Integer | constraints.Float](a []T) T {
	assertNonEmpty(a)
	return T(float64(Sum(a)) / float64(len(a)))
}

// Min computes the minimum element in a vector.
func Min[T constraints.Integer | constraints.Float](a []T) T {
	assertNonEmpty(a)
	result := a[0]
	for i := 0; i < len(a); i++ {
		if a[i] < result {
			result = a[i]
		}
	}
	return result
}

// Max computes the maximum element in a vector.
func Max[T constraints.Integer | constraints.Float](a []T) T {
	assertNonEmpty(a)
	result := a[0]
	for i := 0; i < len(a); i++ {
		if a[i] > result {
			result = a[i]
		}
	}
	return result
}

// Abs computes the absolute value of each element in a vector.
func Abs[T constraints.Integer | constraints.Float](a []T) []T {
	result := make([]T, len(a))
	for i := range a {
		result[i] = a[i]
		if a[i] < 0 {
			result[i] = -a[i]
		}
	}
	return result
}

// AbsInPlace computes the absolute value of each element in a vector and stores the result in the vector.
func AbsInPlace[T constraints.Integer | constraints.Float](a []T) {
	for i := range a {
		if a[i] < 0 {
			a[i] = -a[i]
		}
	}
}

// EuclidianDistance computes the Euclidian distance between two vectors.
func EuclidianDistance[T constraints.Integer | constraints.Float](a, b []T) float64 {
	assertEqualLength(a, b)
	var result float64 = 0
	for i := 0; i < len(a); i++ {
		result += math.Pow(float64(a[i]-b[i]), 2)
	}
	return math.Sqrt(result)
}

func assertEqualLength[T constraints.Integer | constraints.Float](a, b []T) {
	if len(a) != len(b) {
		panic("vectors must be the same length")
	}
}

func assertNonEmpty[T constraints.Integer | constraints.Float](a []T) {
	if len(a) == 0 {
		panic("vector must be non-empty")
	}
}

func assertNotZero[T constraints.Integer | constraints.Float](a T) {
	if a == 0 {
		panic("value must be non-zero")
	}
}
