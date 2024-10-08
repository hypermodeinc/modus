/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"github.com/hypermodeinc/modus/sdk/go/pkg/vectors"
)

func Add(a, b []float64) []float64 {
	return vectors.Add(a, b)
}

func AddInPlace(a, b []float64) []float64 {
	vectors.AddInPlace(a, b)
	return a
}

func Subtract(a, b []float64) []float64 {
	return vectors.Subtract(a, b)
}

func SubtractInPlace(a, b []float64) []float64 {
	vectors.SubtractInPlace(a, b)
	return a
}

func AddNumber(a []float64, b float64) []float64 {
	return vectors.AddNumber(a, b)
}

func AddNumberInPlace(a []float64, b float64) []float64 {
	vectors.AddNumberInPlace(a, b)
	return a
}

func SubtractNumber(a []float64, b float64) []float64 {
	return vectors.SubtractNumber(a, b)
}

func SubtractNumberInPlace(a []float64, b float64) []float64 {
	vectors.SubtractNumberInPlace(a, b)
	return a
}

func MultiplyNumber(a []float64, b float64) []float64 {
	return vectors.MultiplyNumber(a, b)
}

func MultiplyNumberInPlace(a []float64, b float64) []float64 {
	vectors.MultiplyNumberInPlace(a, b)
	return a
}

func DivideNumber(a []float64, b float64) []float64 {
	return vectors.DivideNumber(a, b)
}

func DivideNumberInPlace(a []float64, b float64) []float64 {
	vectors.DivideNumberInPlace(a, b)
	return a
}

func Dot(a, b []float64) float64 {
	return vectors.Dot(a, b)
}

func Magnitude(a []float64) float64 {
	return vectors.Magnitude(a)
}

func Normalize(a []float64) []float64 {
	return vectors.Normalize(a)
}

func Sum(a []float64) float64 {
	return vectors.Sum(a)
}

func Product(a []float64) float64 {
	return vectors.Product(a)
}

func Mean(a []float64) float64 {
	return vectors.Mean(a)
}

func Min(a []float64) float64 {
	return vectors.Min(a)
}

func Max(a []float64) float64 {
	return vectors.Max(a)
}

func Abs(a []float64) []float64 {
	return vectors.Abs(a)
}

func AbsInPlace(a []float64) []float64 {
	vectors.AbsInPlace(a)
	return a
}

func EuclidianDistance(a, b []float64) float64 {
	return vectors.EuclidianDistance(a, b)
}
