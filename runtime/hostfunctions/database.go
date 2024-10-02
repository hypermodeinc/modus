/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"fmt"

	"github.com/hypermodeinc/modus/runtime/sqlclient"
)

func init() {
	registerHostFunction("hypermode", "databaseQuery", sqlclient.ExecuteQuery,
		withStartingMessage("Starting database query."),
		withCompletedMessage("Completed database query."),
		withCancelledMessage("Cancelled database query."),
		withErrorMessage("Error querying database."),
		withMessageDetail(func(hostName, statement string) string {
			return fmt.Sprintf("Host: %s Query: %s", hostName, statement)
		}))
}
