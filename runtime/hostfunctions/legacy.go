/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"fmt"

	"hypruntime/models/legacymodels"
)

func init() {
	registerHostFunction("hypermode", "invokeClassifier", legacymodels.InvokeClassifier,
		withStartingMessage("Invoking model."),
		withCompletedMessage("Completed model invocation."),
		withCancelledMessage("Cancelled model invocation."),
		withErrorMessage("Error invoking model."),
		withMessageDetail(func(modelName string) string {
			return fmt.Sprintf("Model: %s", modelName)
		}))

	registerHostFunction("hypermode", "computeEmbedding", legacymodels.ComputeEmbedding,
		withStartingMessage("Invoking model."),
		withCompletedMessage("Completed model invocation."),
		withCancelledMessage("Cancelled model invocation."),
		withErrorMessage("Error invoking model."),
		withMessageDetail(func(modelName string) string {
			return fmt.Sprintf("Model: %s", modelName)
		}))

	registerHostFunction("hypermode", "invokeTextGenerator", legacymodels.InvokeTextGenerator,
		withStartingMessage("Invoking model."),
		withCompletedMessage("Completed model invocation."),
		withCancelledMessage("Cancelled model invocation."),
		withErrorMessage("Error invoking model."),
		withMessageDetail(func(modelName string) string {
			return fmt.Sprintf("Model: %s", modelName)
		}))
}
