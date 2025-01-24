/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package utils_test

import (
	"encoding/json"
	"testing"

	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/stretchr/testify/assert"
)

func testCast[T any](t *testing.T, expected T, inputs ...any) {
	for _, input := range inputs {
		actual, e := utils.Cast[T](input)
		assert.Nil(t, e)
		assert.Equal(t, expected, actual, "input: %v (%[1]T)", input)
	}
}

func TestCast_Int(t *testing.T) {
	testCast(t, int(42),
		42,
		42.0,
		"42",
		json.Number("42"),
		func() *int { x := 42; return &x }(),
	)
}

func TestCast_Int8(t *testing.T) {
	testCast(t, int8(42),
		42,
		42.0,
		"42",
		json.Number("42"),
		func() *int { x := 42; return &x }(),
	)
}

func TestCast_Int16(t *testing.T) {
	testCast(t, int16(42),
		42,
		42.0,
		"42",
		json.Number("42"),
		func() *int { x := 42; return &x }(),
	)
}

func TestCast_Int32(t *testing.T) {
	testCast(t, int32(42),
		42,
		42.0,
		"42",
		json.Number("42"),
		func() *int { x := 42; return &x }(),
	)
}

func TestCast_Int64(t *testing.T) {
	testCast(t, int64(42),
		42,
		42.0,
		"42",
		json.Number("42"),
		func() *int { x := 42; return &x }(),
	)
}

func TestCast_Uint(t *testing.T) {
	testCast(t, uint(42),
		42,
		42.0,
		"42",
		json.Number("42"),
		func() *int { x := 42; return &x }(),
	)
}

func TestCast_Uint8(t *testing.T) {
	testCast(t, uint8(42),
		42,
		42.0,
		"42",
		json.Number("42"),
		func() *int { x := 42; return &x }(),
	)
}

func TestCast_Uint16(t *testing.T) {
	testCast(t, uint16(42),
		42,
		42.0,
		"42",
		json.Number("42"),
		func() *int { x := 42; return &x }(),
	)
}

func TestCast_Uint32(t *testing.T) {
	testCast(t, uint32(42),
		42,
		42.0,
		"42",
		json.Number("42"),
		func() *int { x := 42; return &x }(),
	)
}

func TestCast_Uint64(t *testing.T) {
	testCast(t, uint64(42),
		42,
		42.0,
		"42",
		json.Number("42"),
		func() *int { x := 42; return &x }(),
	)
}

func TestCast_Float32(t *testing.T) {
	testCast(t, float32(42),
		42,
		42.0,
		"42",
		json.Number("42"),
		func() *int { x := 42; return &x }(),
	)
}

func TestCast_Float64(t *testing.T) {
	testCast(t, float64(42),
		42,
		42.0,
		"42",
		json.Number("42"),
		func() *int { x := 42; return &x }(),
	)
}

func TestCast_Uintptr(t *testing.T) {
	testCast(t, uintptr(42),
		42,
		42.0,
		"42",
		json.Number("42"),
		uintptr(42),
		uint32(42),
		func() *uint32 { x := uint32(42); return &x }(),
		// TODO: *uintptr fails
	)
}

func TestCast_Bool(t *testing.T) {
	testCast(t, true,
		1,
		1.0,
		"true",
		json.Number("1"),
		func() *bool { x := true; return &x }(),
	)
}
