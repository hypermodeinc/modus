/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package postgresql

import "fmt"

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (p *Point) String() string {
	return fmt.Sprintf("(%v,%v)", p.X, p.Y)
}

func (p *Point) MarshalJSON() ([]byte, error) {
	s := p.String()
	b := make([]byte, len(s)+2)
	b[0] = '"'
	copy(b[1:], s)
	b[len(b)-1] = '"'
	return b, nil
}

func (p *Point) UnmarshalJSON(data []byte) error {
	if len(data) < 7 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("invalid point: %s", string(data))
	}

	loc, err := ParsePoint(string(data[1 : len(data)-1]))
	if err != nil {
		return err
	}

	p.X = loc.X
	p.Y = loc.Y
	return nil
}

func NewPoint(x, y float64) *Point {
	return &Point{x, y}
}

func ParsePoint(s string) (*Point, error) {
	var p Point
	_, err := fmt.Sscanf(s, "(%f,%f)", &p.X, &p.Y)
	return &p, err
}
