//go:build wasip1

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package db

import "unsafe"

//go:noescape
//go:wasmimport hypermode databaseQuery
func _databaseQuery(hostName, dbType, statement, paramsJson *string) unsafe.Pointer

//hypermode:import hypermode databaseQuery
func databaseQuery(hostName, dbType, statement, paramsJson *string) *HostQueryResponse {
	response := _databaseQuery(hostName, dbType, statement, paramsJson)
	if response == nil {
		return nil
	}
	return (*HostQueryResponse)(response)
}
