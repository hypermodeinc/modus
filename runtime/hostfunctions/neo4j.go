/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
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
		withStartingMessage("Executing DQL operation."),
		withCompletedMessage("Completed DQL operation."),
		withCancelledMessage("Cancelled DQL operation."),
		withErrorMessage("Error executing DQL operation."),
		withMessageDetail(func(hostName, dbName, query string) string {
			return fmt.Sprintf("Host: %s Database: %s Query: %s", hostName, dbName, query)
		}))
}
