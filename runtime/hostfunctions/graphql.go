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

	"github.com/hypermodeinc/modus/runtime/graphqlclient"
)

func init() {
	const module_name = "modus_graphql_client"

	registerHostFunction(module_name, "executeQuery", graphqlclient.Execute,
		withStartingMessage("Executing GraphQL operation."),
		withCompletedMessage("Completed GraphQL operation."),
		withCancelledMessage("Cancelled GraphQL operation."),
		withErrorMessage("Error executing GraphQL operation."),
		withMessageDetail(func(hostName, stmt string) string {
			return fmt.Sprintf("Host: %s Query: %s", hostName, stmt)
		}))
}
