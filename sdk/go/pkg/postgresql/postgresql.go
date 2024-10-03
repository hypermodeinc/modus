/*
 * Copyright 2024 Hypermode, Inc.
 */

package postgresql

import "github.com/hypermodeAI/functions-go/pkg/db"

const dbType = "postgresql"

func Query[T any](hostName, statement string, params ...any) ([]T, uint, error) {
	return db.Query[T](hostName, dbType, statement, params...)
}

func QueryScalar[T any](hostName, statement string, params ...any) (T, uint, error) {
	return db.QueryScalar[T](hostName, dbType, statement, params...)
}

func Execute(hostName, statement string, params ...any) (uint, error) {
	return db.Execute(hostName, dbType, statement, params...)
}
