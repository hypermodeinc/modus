/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package hostfunctions

import (
	"fmt"

	"github.com/hypermodeinc/modus/runtime/dgraphclient"
)

func init() {
	registerHostFunction("hypermode", "executeDQL", dgraphclient.Execute,
		withStartingMessage("Executing DQL operation."),
		withCompletedMessage("Completed DQL operation."),
		withCancelledMessage("Cancelled DQL operation."),
		withErrorMessage("Error executing DQL operation."),
		withMessageDetail(func(hostName string, req *dgraphclient.Request) string {
			return fmt.Sprintf("Host: %s Req: %s", hostName, fmt.Sprint(req))
		}))

	registerHostFunction("hypermode", "dgraphAlterSchema", dgraphclient.AlterSchema,
		withStartingMessage("Altering DQL schema."),
		withCompletedMessage("Completed DQL schema alteration."),
		withCancelledMessage("Cancelled DQL schema alteration."),
		withErrorMessage("Error altering DQL schema."),
		withMessageDetail(func(hostName, schema string) string {
			return fmt.Sprintf("Host: %s Schema: %s", hostName, schema)
		}))

	registerHostFunction("hypermode", "dgraphDropAttr", dgraphclient.DropAttr,
		withStartingMessage("Dropping DQL attribute."),
		withCompletedMessage("Completed DQL attribute drop."),
		withCancelledMessage("Cancelled DQL attribute drop."),
		withErrorMessage("Error dropping DQL attribute."),
		withMessageDetail(func(hostName, attr string) string {
			return fmt.Sprintf("Host: %s Attribute: %s", hostName, attr)
		}))

	registerHostFunction("hypermode", "dgraphDropAll", dgraphclient.DropAll,
		withStartingMessage("Dropping all DQL data."),
		withCompletedMessage("Completed DQL data drop."),
		withCancelledMessage("Cancelled DQL data drop."),
		withErrorMessage("Error dropping DQL data."),
		withMessageDetail(func(hostName string) string {
			return fmt.Sprintf("Host: %s", hostName)
		}))
}
