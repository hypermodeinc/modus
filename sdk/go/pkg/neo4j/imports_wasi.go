//go:build wasip1

/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package neo4j

import "unsafe"

//go:noescape
//go:wasmimport modus_neo4j_client executeQuery
func _hostExecuteQuery(hostName, dbName, query, parametersJson *string) unsafe.Pointer

//modus:import modus_neo4j_client executeQuery
func hostExecuteQuery(hostName, dbName, query, parametersJson *string) *EagerResult {
	response := _hostExecuteQuery(hostName, dbName, query, parametersJson)
	if response == nil {
		return nil
	}
	return (*EagerResult)(response)
}
