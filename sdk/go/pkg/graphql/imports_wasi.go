//go:build wasip1

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package graphql

//go:noescape
//go:wasmimport modus_graphql_client executeQuery
func hostExecuteQuery(connection, statement, variables *string) *string
