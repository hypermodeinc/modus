/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package hostfunctions

import (
	"fmt"

	"github.com/hypermodeinc/modus/runtime/dgraphclient"
)

func init() {
	const module_name = "modus_dgraph_client"

	registerHostFunction(module_name, "executeQuery", dgraphclient.ExecuteQuery,
		withStartingMessage("Executing DQL operation."),
		withCompletedMessage("Completed DQL operation."),
		withCancelledMessage("Cancelled DQL operation."),
		withErrorMessage("Error executing DQL operation."),
		withMessageDetail(func(hostName string, req *dgraphclient.Request) string {
			return fmt.Sprintf("Host: %s Req: %s", hostName, fmt.Sprint(req))
		}))

	registerHostFunction(module_name, "alterSchema", dgraphclient.AlterSchema,
		withStartingMessage("Altering DQL schema."),
		withCompletedMessage("Completed DQL schema alteration."),
		withCancelledMessage("Cancelled DQL schema alteration."),
		withErrorMessage("Error altering DQL schema."),
		withMessageDetail(func(hostName, schema string) string {
			return fmt.Sprintf("Host: %s Schema: %s", hostName, schema)
		}))

	registerHostFunction(module_name, "dropAttribute", dgraphclient.DropAttribute,
		withStartingMessage("Dropping DQL attribute."),
		withCompletedMessage("Completed DQL attribute drop."),
		withCancelledMessage("Cancelled DQL attribute drop."),
		withErrorMessage("Error dropping DQL attribute."),
		withMessageDetail(func(hostName, attr string) string {
			return fmt.Sprintf("Host: %s Attribute: %s", hostName, attr)
		}))

	registerHostFunction(module_name, "dropAllData", dgraphclient.DropAllData,
		withStartingMessage("Dropping all DQL data."),
		withCompletedMessage("Completed DQL data drop."),
		withCancelledMessage("Cancelled DQL data drop."),
		withErrorMessage("Error dropping DQL data."),
		withMessageDetail(func(hostName string) string {
			return fmt.Sprintf("Host: %s", hostName)
		}))
}
