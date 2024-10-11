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

	"github.com/hypermodeinc/modus/runtime/models"
)

func init() {
	const module_name = "modus_models"

	registerHostFunction(module_name, "getModelInfo", models.GetModelInfo,
		withCancelledMessage("Cancelled getting model info."),
		withErrorMessage("Error getting model info."),
		withMessageDetail(func(modelName string) string {
			return fmt.Sprintf("Model: %s", modelName)
		}))

	registerHostFunction(module_name, "invokeModel", models.InvokeModel,
		withStartingMessage("Invoking model."),
		withCompletedMessage("Completed model invocation."),
		withCancelledMessage("Cancelled model invocation."),
		withErrorMessage("Error invoking model."),
		withMessageDetail(func(modelName string) string {
			return fmt.Sprintf("Model: %s", modelName)
		}))
}
