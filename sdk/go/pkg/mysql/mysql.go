/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package mysql

import "github.com/hypermodeinc/modus/sdk/go/pkg/db"

const dbType = "mysql"

// Executes a database query that does not return rows.
func Execute(connection, statement string, params ...any) (*db.Response, error) {
	return db.Execute(connection, dbType, statement, params...)
}

// Executes a database query that returns rows.
// The structure of the rows is determined by the type parameter.
func Query[T any](connection, statement string, params ...any) (*db.QueryResponse[T], error) {
	return db.Query[T](connection, dbType, statement, params...)
}

// Executes a database query that returns a single scalar value.
// The type parameter determines the type of the scalar value.
func QueryScalar[T any](connection, statement string, params ...any) (*db.ScalarResponse[T], error) {
	return db.QueryScalar[T](connection, dbType, statement, params...)
}

// Represents a point in 2D space, having X and Y coordinates.
// Correctly serializes to and from a SQL point type, in (X, Y) order.
//
// Note that this struct is identical to the Location struct, but uses different field names.
type Point = db.Point

// Represents a location on Earth, having longitude and latitude coordinates.
// Correctly serializes to and from a SQL point type, in (longitude, latitude) order.
//
// Note that this struct is identical to the Point struct, but uses different field names.
type Location = db.Location

// Creates a new Point with the specified X and Y coordinates.
func NewPoint(x, y float64) *Point {
	return db.NewPoint(x, y)
}

// Creates a new Location with the specified longitude and latitude coordinates.
func NewLocation(longitude, latitude float64) *Location {
	return db.NewLocation(longitude, latitude)
}
