/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"

	"hypruntime/logger"
	"hypruntime/utils"
)

func init() {
	registerHostFunction("hypermode", "log", logFunctionMessage)
}

func logFunctionMessage(ctx context.Context, level, message string) {

	// write to the logger
	logger.Get(ctx).
		WithLevel(logger.ParseLevel(level)).
		Str("text", message).
		Bool("user_visible", true).
		Msg("Message logged from function.")

	// also store messages in the context, so we can return them to the caller
	messages := ctx.Value(utils.FunctionMessagesContextKey).(*[]utils.LogMessage)
	*messages = append(*messages, utils.LogMessage{
		Level:   level,
		Message: message,
	})
}
