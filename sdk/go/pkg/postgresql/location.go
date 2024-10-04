/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package postgresql

import "fmt"

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func (l *Location) String() string {
	return fmt.Sprintf("(%v,%v)", l.Longitude, l.Latitude)
}

func (l *Location) MarshalJSON() ([]byte, error) {
	s := l.String()
	b := make([]byte, len(s)+2)
	b[0] = '"'
	copy(b[1:], s)
	b[len(b)-1] = '"'
	return b, nil
}

func (l *Location) UnmarshalJSON(data []byte) error {
	if len(data) < 7 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("invalid location: %s", string(data))
	}

	loc, err := ParseLocation(string(data[1 : len(data)-1]))
	if err != nil {
		return err
	}

	l.Longitude = loc.Longitude
	l.Latitude = loc.Latitude
	return nil
}

func NewLocation(longitude, latitude float64) *Location {
	return &Location{longitude, latitude}
}

func ParseLocation(s string) (*Location, error) {
	var l Location
	_, err := fmt.Sscanf(s, "(%f,%f)", &l.Longitude, &l.Latitude)
	return &l, err
}
