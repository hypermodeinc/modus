/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package hostfunctions

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/sentryutils"
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

	var zoneId string
	if tz != nil && *tz != "" {
		zoneId = *tz
	} else if tz, ok := ctx.Value(utils.TimeZoneContextKey).(string); ok {
		zoneId = tz
	} else {
		const msg = "Time zone not specified."
		sentryutils.CaptureError(ctx, nil, msg)
		logger.Error(ctx).Msg(msg)
		return nil
	}

	loc, err := timezones.GetLocation(ctx, zoneId)
	if err != nil {
		const msg = "Failed to get time zone location."
		sentryutils.CaptureError(ctx, err, msg, sentryutils.WithData("tz", zoneId))
		logger.Error(ctx, err).Str("tz", zoneId).Msg(msg)
		return nil
	}

	s := now.In(loc).Format(time.RFC3339Nano)
	return &s
}

func GetTimeZoneData(ctx context.Context, tz, format *string) []byte {
	if tz == nil {
		const msg = "Time zone not specified."
		sentryutils.CaptureError(ctx, nil, msg)
		logger.Error(ctx).Msg(msg)
		return nil
	}
	if format == nil {
		const msg = "Time zone format not specified."
		sentryutils.CaptureError(ctx, nil, msg)
		logger.Error(ctx).Msg(msg)
		return nil
	}
	data, err := timezones.GetTimeZoneData(ctx, *tz, *format)
	if err != nil {
		const msg = "Failed to get time zone data."
		sentryutils.CaptureError(ctx, err, msg, sentryutils.WithData("tz", *tz), sentryutils.WithData("format", *format))
		logger.Error(ctx, err).Str("tz", *tz).Str("format", *format).Msg(msg)
		return nil
	}
	return data
}
