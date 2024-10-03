//go:build wasip1

/*
 * Copyright 2024 Hypermode, Inc.
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
