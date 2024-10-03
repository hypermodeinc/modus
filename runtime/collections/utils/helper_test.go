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
	"reflect"
	"testing"
)

func TestConvertToFloat32_2DArray(t *testing.T) {
	expected := [][]float32{
		{1.0, 2.0},
		{3.0, 4.0},
	}
	// Test with float32, not wrapped in interface
	input32 := [][]float32{
		{1.0, 2.0},
		{3.0, 4.0},
	}
	result, err := ConvertToFloat32_2DArray(input32)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Test with float64, not wrapped in interface
	input := [][]float64{
		{1.0, 2.0},
		{3.0, 4.0},
	}
	result, err = ConvertToFloat32_2DArray(input)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}

	// Test with invalid input
	inputInvalid := []any{1.0, "invalid"}
	_, err = ConvertToFloat32_2DArray(inputInvalid)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
