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

	"github.com/hypermodeinc/modus/runtime/modusdb"
)

func init() {
	const module_name = "modus_modusdb_client"

	registerHostFunction(module_name, "dropAll", modusdb.ExtDropAll,
		withStartingMessage("Dropping all."),
		withCompletedMessage("Completed dropping all."),
		withCancelledMessage("Cancelled dropping all."),
		withErrorMessage("Error dropping all."),
		withMessageDetail(func() string {
			return "Dropping all."
		}))

	registerHostFunction(module_name, "dropData", modusdb.ExtDropData,
		withStartingMessage("Dropping data."),
		withCompletedMessage("Completed dropping data."),
		withCancelledMessage("Cancelled dropping data."),
		withErrorMessage("Error dropping data."),
		withMessageDetail(func() string {
			return "Dropping data."
		}))

	registerHostFunction(module_name, "alterSchema", modusdb.ExtAlterSchema,
		withStartingMessage("Altering schema."),
		withCompletedMessage("Completed altering schema."),
		withCancelledMessage("Cancelled altering schema."),
		withErrorMessage("Error altering schema."),
		withMessageDetail(func(schema string) string {
			return fmt.Sprintf("Schema: %s", schema)
		}))

	registerHostFunction(module_name, "mutate", modusdb.ExtMutate,
		withStartingMessage("Executing mutation."),
		withCompletedMessage("Completed mutation."),
		withCancelledMessage("Cancelled mutation."),
		withErrorMessage("Error executing mutation."),
		withMessageDetail(func(mutationReq modusdb.MutationRequest) string {
			return fmt.Sprintf("Mutations: %s", fmt.Sprint(mutationReq))
		}))

	registerHostFunction(module_name, "query", modusdb.ExtQuery,
		withStartingMessage("Executing query."),
		withCompletedMessage("Completed query."),
		withCancelledMessage("Cancelled query."),
		withErrorMessage("Error executing query."),
		withMessageDetail(func(query string) string {
			return fmt.Sprintf("Query: %s", query)
		}))
}
