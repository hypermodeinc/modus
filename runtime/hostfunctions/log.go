/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"fmt"
	"os"

	"hypruntime/logger"
	"hypruntime/utils"
)

func init() {
	registerHostFunction("hypermode", "log", LogFunctionMessage)
}

func LogFunctionMessage(ctx context.Context, level, message string) {

	// store messages in the context, so we can return them to the caller
	messages := ctx.Value(utils.FunctionMessagesContextKey).(*[]utils.LogMessage)
	*messages = append(*messages, utils.LogMessage{
		Level:   level,
		Message: message,
	})

	// If debugging, write debug messages to stderr instead of the logger
	if level == "debug" && utils.HypermodeDebugEnabled() {
		fmt.Fprintln(os.Stderr, message)
		return
	}

	// write to the logger
	logger.Get(ctx).
		WithLevel(logger.ParseLevel(level)).
		Str("text", message).
		Bool("user_visible", true).
		Msg("Message logged from function.")
}
