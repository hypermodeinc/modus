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
	"time"

	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/timezones"
	"github.com/hypermodeinc/modus/runtime/utils"
)

func init() {
	const module_name = "modus_system"

	registerHostFunction(module_name, "logMessage", LogMessage)
	registerHostFunction(module_name, "getTimeInZone", GetTimeInZone)
	registerHostFunction(module_name, "getTimeZoneData", GetTimeZoneData)
}

func LogMessage(ctx context.Context, level, message string) {

	// store messages in the context, so we can return them to the caller
	if messages, ok := ctx.Value(utils.FunctionMessagesContextKey).(*[]utils.LogMessage); ok {
		*messages = append(*messages, utils.LogMessage{
			Level:   level,
			Message: message,
		})
	}

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

func GetTimeInZone(ctx context.Context, tz *string) *string {
	now := time.Now()

	var loc *time.Location
	if tz != nil && *tz != "" {
		loc = timezones.GetLocation(*tz)
	} else if tz, ok := ctx.Value(utils.TimeZoneContextKey).(string); ok {
		loc = timezones.GetLocation(tz)
	}

	if loc == nil {
		return nil
	}

	s := now.In(loc).Format(time.RFC3339Nano)
	return &s
}

func GetTimeZoneData(tz, format *string) []byte {
	if tz == nil {
		return nil
	}

	return timezones.GetTimeZoneData(*tz, *format)
}
