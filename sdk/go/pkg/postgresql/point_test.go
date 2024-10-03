/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package postgresql_test

import (
	"encoding/json"
	"testing"

	"github.com/hypermodeAI/functions-go/pkg/postgresql"
)

func TestPointString(t *testing.T) {
	point := postgresql.NewPoint(12.345678901234567, -56.7890123456789)
	expected := "(12.345678901234567,-56.7890123456789)"
	result := point.String()

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestPointMarshalJSON(t *testing.T) {
	point := postgresql.NewPoint(12.345678901234567, -56.7890123456789)
	expected := `"(12.345678901234567,-56.7890123456789)"`
	result, err := json.Marshal(point)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if string(result) != expected {
		t.Errorf("Expected %s, but got %s", expected, string(result))
	}
}

func TestPointUnmarshalJSON(t *testing.T) {
	data := []byte(`"(12.345678901234567,-56.7890123456789)"`)
	expected := postgresql.NewPoint(12.345678901234567, -56.7890123456789)
	point := &postgresql.Point{}

	err := json.Unmarshal(data, point)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if point.X != expected.X || point.Y != expected.Y {
		t.Errorf("Expected %+v, but got %+v", expected, point)
	}
}

func TestNewPoint(t *testing.T) {
	x := 12.345678901234567
	y := -56.7890123456789
	point := postgresql.NewPoint(x, y)

	if point.X != x || point.Y != y {
		t.Errorf("Expected X: %f, Y: %f, but got %+v", x, y, point)
	}
}

func TestParsePoint(t *testing.T) {
	s := "(12.345678901234567,-56.7890123456789)"
	expected := postgresql.NewPoint(12.345678901234567, -56.7890123456789)
	point, err := postgresql.ParsePoint(s)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if point.X != expected.X || point.Y != expected.Y {
		t.Errorf("Expected %+v, but got %+v", expected, point)
	}
}

func TestParsePointError(t *testing.T) {
	s := "invalid-point"
	_, err := postgresql.ParsePoint(s)

	if err == nil {
		t.Error("Expected an error, but received none")
	}
}
