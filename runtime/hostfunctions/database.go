/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package hostfunctions

import (
	"fmt"

	"github.com/hypermodeinc/modus/runtime/sqlclient"
)

func init() {
	const module_name = "modus_sql_client"

	registerHostFunction(module_name, "executeQuery", sqlclient.ExecuteQuery,
		withStartingMessage("Starting database query."),
		withCompletedMessage("Completed database query."),
		withCancelledMessage("Cancelled database query."),
		withErrorMessage("Error querying database."),
		withMessageDetail(func(hostName, statement string) string {
			return fmt.Sprintf("Host: %s Query: %s", hostName, statement)
		}))
}
