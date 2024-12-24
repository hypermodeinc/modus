//go:build windows

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
	"fmt"
	"time"
	"unsafe"

	_ "time/tzdata"
	_ "unsafe"

	"golang.org/x/sys/windows"
)

// note: the following link requires -ldflags "-checklinkname=0" on go 1.23 or later

//go:linkname loadFromEmbeddedTZData time/tzdata.loadFromEmbeddedTZData
func loadFromEmbeddedTZData(string) (string, error)

func loadTimeZoneInfo(tz string) (*tzInfo, error) {

	var data []byte
	if s, err := loadFromEmbeddedTZData(tz); err != nil {
		return nil, fmt.Errorf("could not load timezone data: %v", err)
	} else {
		data = []byte(s)
	}

	loc, err := time.LoadLocationFromTZData(tz, data)
	if err != nil {
		return nil, fmt.Errorf("could not load timezone data: %v", err)
	}

	info := &tzInfo{loc, data}
	return info, nil
}

func getSystemLocalTimeZone() (string, error) {
	// On Windows, we use the ICU library to get the default time zone in IANA format.
	// This requires Windows 10 release 1703 or later.
	// See https://learn.microsoft.com/en-us/windows/win32/intl/international-components-for-unicode--icu-

	// We also import time/tzdata to have an embedded copy of the IANA time zone database.

	var errorCode int32
	const bufferSize = 128
	buffer := make([]uint16, bufferSize)

	lib := windows.NewLazySystemDLL("icuin.dll")
	proc := lib.NewProc("ucal_getDefaultTimeZone")

	ret, _, err := proc.Call(
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(bufferSize),
		uintptr(unsafe.Pointer(&errorCode)),
	)

	if ret == 0 || errorCode != 0 {
		return "", fmt.Errorf("failed to determine system local time zone: %v [0x%x]", err, errorCode)
	}

	tz := windows.UTF16ToString(buffer)
	return tz, nil
}
