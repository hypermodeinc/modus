/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"fmt"

	"hypruntime/models"
)

func init() {
	registerHostFunction("hypermode", "lookupModel", models.GetModelInfo,
		withCancelledMessage("Cancelled getting model info."),
		withErrorMessage("Error getting model info."),
		withMessageDetail(func(modelName string) string {
			return fmt.Sprintf("Model: %s", modelName)
		}))

	registerHostFunction("hypermode", "invokeModel", models.InvokeModel,
		withStartingMessage("Invoking model."),
		withCompletedMessage("Completed model invocation."),
		withCancelledMessage("Cancelled model invocation."),
		withErrorMessage("Error invoking model."),
		withMessageDetail(func(modelName string) string {
			return fmt.Sprintf("Model: %s", modelName)
		}))
}
