/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package timezones

import (
	"os"
	"time"

	"github.com/puzpuzpuz/xsync/v4"
)

type tzInfo struct {
	location *time.Location
	data     []byte
}

var systemTimeZone string
var tzCache = *xsync.NewMap[string, *tzInfo]()

func init() {
	if tz, err := getSystemLocalTimeZone(); err == nil {
		systemTimeZone = tz
	} else {
		// silently fallback to UTC
		systemTimeZone = "UTC"
	}
}

func GetLocalTimeZone() string {
	// check every time, in case the env var has changed via .env file reload
	if tz := os.Getenv("TZ"); tz != "" {
		return tz
	}

	return systemTimeZone
}

func GetLocation(tz string) *time.Location {
	if tz == "" {
		return nil
	}

	info, err := getTimeZoneInfo(tz)
	if err != nil {
		return nil
	}

	return info.location
}

func GetTimeZoneData(tz, format string) []byte {
	if tz == "" {
		return nil
	}

	// only support tzif format for now
	// we can expand this to support other formats as needed
	if format != "tzif" {
		return nil
	}

	info, err := getTimeZoneInfo(tz)
	if err != nil {
		return nil
	}

	return info.data
}

func getTimeZoneInfo(tz string) (*tzInfo, error) {
	if info, ok := tzCache.Load(tz); ok {
		return info, nil
	}

	info, err := loadTimeZoneInfo(tz)
	if err != nil {
		return nil, err
	}

	tzCache.Store(tz, info)
	return info, nil
}
