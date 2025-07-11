/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package hostfunctions

import (
	"fmt"

	"github.com/hypermodeinc/modus/runtime/neo4jclient"
)

func init() {
	const module_name = "modus_neo4j_client"

	registerHostFunction(module_name, "executeQuery", neo4jclient.ExecuteQuery,
		withStartingMessage("Executing Neo4j query."),
		withCompletedMessage("Completed Neo4j query."),
		withCancelledMessage("Cancelled Neo4j query."),
		withErrorMessage("Error executing Neo4j query."),
		withMessageDetail(func(hostName, dbName, query string) string {
			return fmt.Sprintf("Host: %s Database: %s Query: %s", hostName, dbName, query)
		}))
}
