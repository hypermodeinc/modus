/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package postgresql

import "github.com/hypermodeinc/modus/sdk/go/pkg/db"

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
