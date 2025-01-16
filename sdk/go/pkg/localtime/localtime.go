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
	"errors"
	"os"
	"time"
)

// Now returns the current local time.  The time zone is determined in the following order of precedence:
// - If the X-Time-Zone header is present in the request, the time zone is set to the value of the header.
// - If the TZ environment variable is set on the host, the time zone is set to the value of the variable.
// - Otherwise, the time zone is set to the host's local time zone.
func Now() (time.Time, error) {

	// Try to return a time with the proper full time zone.
	if tz := os.Getenv("TZ"); tz != "" {
		if loc, err := hostGetTimeLocation(tz); err == nil {
			return time.Now().In(loc), nil
		}
	}

	// Otherwise, let the host get the local time with the current offset applied.
	return hostGetLocalTime()
}

// NowInZone returns the current time in the given time zone.
// The time zone should be a valid IANA time zone identifier, such as "America/New_York".
func NowInZone(tz string) (time.Time, error) {
	if tz == "" {
		return time.Time{}, errors.New("a time zone is required")
	}

	// Try to return a time with the proper full time zone.
	if loc, err := hostGetTimeLocation(tz); err == nil {
		return time.Now().In(loc), nil
	}

	// Otherwise, let the host get the local time with the current offset applied.
	return hostGetTimeInZone(tz)
}

// GetTimeZone returns the local time zone identifier, in IANA format.
func GetTimeZone() string {
	return os.Getenv("TZ")
}

// GetLocation returns the time.Location for the given time zone.
// The time zone should be a valid IANA time zone identifier, such as "America/New_York".
func GetLocation(tz string) (*time.Location, error) {
	if tz == "" {
		return nil, errors.New("time zone is required")
	}
	return hostGetTimeLocation(tz)
}

// IsValidTimeZone returns true if the given time zone is valid.
func IsValidTimeZone(tz string) bool {
	_, err := hostGetTimeInZone(tz)
	return err == nil
}
