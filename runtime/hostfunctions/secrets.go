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

	"github.com/hypermodeinc/modus/runtime/secrets"
)

func init() {
	const module_name = "modus_secrets"

	registerHostFunction(module_name, "getSecretValue", secrets.GetSecretValue,
		withStartingMessage("Starting secret lookup."),
		withCompletedMessage("Completed secret lookup."),
		withCancelledMessage("Cancelled secret lookup."),
		withErrorMessage("Error getting secret."),
		withMessageDetail(func(name string) string {
			return fmt.Sprintf("Secret: %s", name)
		}))
}
