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
	"context"
	"fmt"
	"os"

	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/utils"
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
	if level == "debug" && utils.DebugModeEnabled() {
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
