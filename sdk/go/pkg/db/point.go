/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package db

import (
	"fmt"
	"strings"
)

// Represents a point in 2D space, having X and Y coordinates.
// Correctly serializes to and from a SQL point type, in (X, Y) order.
//
// Note that this struct is identical to the Location struct, but uses different field names.
type Point struct {

	// The X coordinate of the point.
	X float64 `json:"x"`

	// The Y coordinate of the point.
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

// Creates a new Point with the specified X and Y coordinates.
func NewPoint(x, y float64) *Point {
	return &Point{x, y}
}

// Parses a Point from a string in the format "(X,Y)", or "POINT (X Y)".
func ParsePoint(s string) (*Point, error) {
	var p Point
	var err error
	if strings.HasPrefix(s, "POINT (") {
		_, err = fmt.Sscanf(s, "POINT (%f %f)", &p.X, &p.Y)
	} else {
		_, err = fmt.Sscanf(s, "(%f,%f)", &p.X, &p.Y)
	}
	return &p, err
}
