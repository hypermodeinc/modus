//go:build wasip1

/*
 * Copyright 2024 Hypermode, Inc.
 */

package graphql

//go:noescape
//go:wasmimport hypermode executeGQL
func hostExecuteGQL(hostName, statement, variables *string) *string
