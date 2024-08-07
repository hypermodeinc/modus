/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"

	"hmruntime/logger"
	"hmruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

func hostLog(ctx context.Context, mod wasm.Module, pLevel uint32, pMessage uint32) {

	// read input parameters
	var level, message string
	err := readParams(ctx, mod, param{pLevel, &level}, param{pMessage, &message})
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return
	}

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
