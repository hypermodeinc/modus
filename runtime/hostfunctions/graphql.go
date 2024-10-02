/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"fmt"

	"github.com/hypermodeinc/modus/runtime/graphqlclient"
)

func init() {
	registerHostFunction("hypermode", "executeGQL", graphqlclient.Execute,
		withStartingMessage("Executing GraphQL operation."),
		withCompletedMessage("Completed GraphQL operation."),
		withCancelledMessage("Cancelled GraphQL operation."),
		withErrorMessage("Error executing GraphQL operation."),
		withMessageDetail(func(hostName, stmt string) string {
			return fmt.Sprintf("Host: %s Query: %s", hostName, stmt)
		}))
}
