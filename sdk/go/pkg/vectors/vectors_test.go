/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package vectors

import (
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	a := []uint8{1, 2, 3}
	b := []uint8{4, 5, 6}
	expected := []uint8{5, 7, 9}

	result := Add(a, b)

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Add failed, expected %v, got %v", expected, result)
		}
	}
}

func TestAddInPlace(t *testing.T) {
	a := []uint8{1, 2, 3}
	b := []uint8{4, 5, 6}
	expected := []uint8{5, 7, 9}

	AddInPlace(a, b)

	for i := range a {
		if a[i] != expected[i] {
			t.Errorf("AddInPlace failed, expected %v, got %v", expected, a)
		}
	}
}

func TestSubtract(t *testing.T) {
	a := []uint8{4, 5, 6}
	b := []uint8{1, 2, 3}
	expected := []uint8{3, 3, 3}

	result := Subtract(a, b)

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Subtract failed, expected %v, got %v", expected, result)
		}
	}
}

func TestSubtractInPlace(t *testing.T) {
	a := []uint8{4, 5, 6}
	b := []uint8{1, 2, 3}
	expected := []uint8{3, 3, 3}

	SubtractInPlace(a, b)

	for i := range a {
		if a[i] != expected[i] {
			t.Errorf("SubtractInPlace failed, expected %v, got %v", expected, a)
		}
	}
}

func TestAddNumber(t *testing.T) {
	a := []uint8{1, 2, 3}
	b := uint8(4)
	expected := []uint8{5, 6, 7}

	result := AddNumber(a, b)

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("AddNumber failed, expected %v, got %v", expected, result)
		}
	}
}

func TestAddNumberInPlace(t *testing.T) {
	a := []uint8{1, 2, 3}
	b := uint8(4)
	expected := []uint8{5, 6, 7}

	AddNumberInPlace(a, b)

	for i := range a {
		if a[i] != expected[i] {
			t.Errorf("AddNumberInPlace failed, expected %v, got %v", expected, a)
		}
	}
}

func TestSubtractNumber(t *testing.T) {
	a := []uint8{4, 5, 6}
	b := uint8(1)
	expected := []uint8{3, 4, 5}

	result := SubtractNumber(a, b)

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("SubtractNumber failed, expected %v, got %v", expected, result)
		}
	}
}

func TestSubtractNumberInPlace(t *testing.T) {
	a := []uint8{4, 5, 6}
	b := uint8(1)
	expected := []uint8{3, 4, 5}

	SubtractNumberInPlace(a, b)

	for i := range a {
		if a[i] != expected[i] {
			t.Errorf("SubtractNumberInPlace failed, expected %v, got %v", expected, a)
		}
	}
}

func TestMultiplyNumber(t *testing.T) {
	a := []uint8{1, 2, 3}
	b := uint8(4)
	expected := []uint8{4, 8, 12}

	result := MultiplyNumber(a, b)

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("MultiplyNumber failed, expected %v, got %v", expected, result)
		}
	}
}

func TestMultiplyNumberInPlace(t *testing.T) {
	a := []uint8{1, 2, 3}
	b := uint8(4)
	expected := []uint8{4, 8, 12}

	MultiplyNumberInPlace(a, b)

	for i := range a {
		if a[i] != expected[i] {
			t.Errorf("MultiplyNumberInPlace failed, expected %v, got %v", expected, a)
		}
	}
}

func TestDivideNumber(t *testing.T) {
	a := []uint8{4, 8, 12}
	b := uint8(4)
	expected := []uint8{1, 2, 3}

	result := DivideNumber(a, b)

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("DivideNumber failed, expected %v, got %v", expected, result)
		}
	}
}

func TestDivideNumberInPlace(t *testing.T) {
	a := []uint8{4, 8, 12}
	b := uint8(4)
	expected := []uint8{1, 2, 3}

	DivideNumberInPlace(a, b)

	for i := range a {
		if a[i] != expected[i] {
			t.Errorf("DivideNumberInPlace failed, expected %v, got %v", expected, a)
		}
	}
}

func TestDot(t *testing.T) {
	a := []float64{1.0, 2.0, 3.0}
	b := []float64{4.0, 5.0, 6.0}
	expected := 32.0 // (1*4 + 2*5 + 3*6)

	result := Dot(a, b)

	if result != expected {
		t.Errorf("Dot failed, expected %v, got %v", expected, result)
	}
}

func TestMagnitude(t *testing.T) {
	a := []float64{3.0, 4.0}
	expected := 5.0 // sqrt(3^2 + 4^2)

	result := Magnitude(a)

	if result != expected {
		t.Errorf("Magnitude failed, expected %v, got %v", expected, result)
	}
}

func TestNormalize(t *testing.T) {
	a := []float64{3.0, 4.0}
	mag := 5.0 // sqrt(3^2 + 4^2)
	expected := []float64{3.0 / mag, 4.0 / mag}

	result := Normalize(a)

	for i := range result {
		if math.Abs(result[i]-expected[i]) > 1e-9 { // Using epsilon for floating-point comparison
			t.Errorf("Normalize failed, expected %v, got %v", expected, result)
		}
	}
}

func TestSum(t *testing.T) {
	a := []uint8{1, 2, 3}
	expected := uint8(6) // 1 + 2 + 3

	result := Sum(a)

	if result != expected {
		t.Errorf("Sum failed, expected %v, got %v", expected, result)
	}
}

func TestProduct(t *testing.T) {
	a := []uint8{1, 2, 3}
	expected := uint8(6) // 1 * 2 * 3

	result := Product(a)

	if result != expected {
		t.Errorf("Product failed, expected %v, got %v", expected, result)
	}
}

func TestMean(t *testing.T) {
	a := []uint8{1, 2, 3}
	expected := uint8(2) // (1 + 2 + 3) / 3

	result := Mean(a)

	if result != expected {
		t.Errorf("Mean failed, expected %v, got %v", expected, result)
	}
}

func TestMin(t *testing.T) {
	a := []uint8{1, 2, 3}
	expected := uint8(1)

	result := Min(a)

	if result != expected {
		t.Errorf("Min failed, expected %v, got %v", expected, result)
	}
}

func TestMax(t *testing.T) {
	a := []uint8{1, 2, 3}
	expected := uint8(3)

	result := Max(a)

	if result != expected {
		t.Errorf("Max failed, expected %v, got %v", expected, result)
	}
}

func TestAbs(t *testing.T) {
	a := []float64{-1.0, 2.0, -3.0}
	expected := []float64{1.0, 2.0, 3.0}

	result := Abs(a)

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("Abs failed, expected %v, got %v", expected, result)
		}
	}
}

func TestAbsInPlace(t *testing.T) {
	a := []float64{-1.0, 2.0, -3.0}
	expected := []float64{1.0, 2.0, 3.0}

	AbsInPlace(a)

	for i := range a {
		if a[i] != expected[i] {
			t.Errorf("AbsInPlace failed, expected %v, got %v", expected, a)
		}
	}
}

func TestEuclidianDistance(t *testing.T) {
	a := []float64{1.0, 2.0, 3.0}
	b := []float64{4.0, 5.0, 6.0}
	expected := math.Sqrt(27) // sqrt((1-4)^2 + (2-5)^2 + (3-6)^2)

	result := EuclidianDistance(a, b)

	if math.Abs(result-expected) > 1e-9 {
		t.Errorf("EuclidianDistance failed, expected %v, got %v", expected, result)
	}
}
