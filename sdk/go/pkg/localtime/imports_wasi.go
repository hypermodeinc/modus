//go:build wasip1

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package localtime

import (
	"fmt"
	"time"
	"unsafe"
)

//go:noescape
//go:wasmimport modus_system getTimeInZone
func _hostGetTimeInZone(tz *string) *string

func hostGetLocalTime() (time.Time, error) {
	ts := _hostGetTimeInZone(nil)
	if ts == nil {
		return time.Time{}, fmt.Errorf("failed to get the local time")
	}
	return time.Parse(time.RFC3339Nano, *ts)
}

func hostGetTimeInZone(tz string) (time.Time, error) {
	ts := _hostGetTimeInZone(&tz)
	if ts == nil {
		return time.Time{}, fmt.Errorf("failed to get time in %s", tz)
	}
	return time.Parse(time.RFC3339Nano, *ts)
}

//go:noescape
//go:wasmimport modus_system getTimeZoneData
func _hostGetTimeZoneData(tz, format *string) unsafe.Pointer

//modus:import modus_system getTimeZoneData
func hostGetTimeZoneData(tz, format *string) *[]byte {
	data := _hostGetTimeZoneData(tz, format)
	if data == nil {
		return nil
	}
	return (*[]byte)(data)
}

func hostGetTimeLocation(tz string) (*time.Location, error) {
	format := "tzif"
	data := hostGetTimeZoneData(&tz, &format)
	if data == nil {
		return nil, fmt.Errorf("timezone data not found for %s", tz)
	}

	loc, err := time.LoadLocationFromTZData(tz, *data)
	if err != nil {
		return nil, fmt.Errorf("failed to load timezone data for %s: %w", tz, err)
	}

	return loc, nil
}
