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

func init() {
	addHostFunction(&hostFunctionDefinition{
		name:     "log",
		function: wasm.GoModuleFunc(hostLog),
		params:   []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
		results:  []wasm.ValueType{},
	})
}

func hostLog(ctx context.Context, mod wasm.Module, stack []uint64) {

	// read input parameters
	var level, message string
	if err := readParams(ctx, mod, stack, &level, &message); err != nil {
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
