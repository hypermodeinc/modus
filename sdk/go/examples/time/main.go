/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"fmt"
	"time"

	"github.com/hypermodeinc/modus/sdk/go/pkg/localtime"
)

/*
 * NOTES:
 *
 * - Typically in Go, `time.Now()` returns the current time in the local time zone.  But in Modus (due to Go's WASI implementation),
 *   it returns the current time in UTC - even though `time.Now().Location()` still returns "Local".
 *   Thus, to get the local time, you should use the `localtime` package provided by Modus - as shown in the examples below.
 *   If you actually want UTC, you should prefer `time.Now().UTC()` - which is unambiguous.
 *
 * - If you return a `time.Time` object to Modus, it will be converted to UTC - regardless of the time zone.
 *   Thus, if you need to return a time object to Modus, you should return a formatted string instead - Preferably in RFC3339 (ISO 8601) format.
 *   You can use `time.RFC3339``, or `time.RFC3339Nano`` for more precision, or any other format you prefer.
 *
 * - The `localTime.GetLocation` function is very similar to `time.LoadLocation` built into Go. For efficiency, `localtime.GetLocation`
 *   gets time zone data from the Modus runtime - which gets it from the host OS in most cases. However, if you are importing a third-party
 *   library that expects `time.LoadLocation` to work correctly, you can import Go's `time/tzdata` package, which will embeds the time zone
 *   database in your wasm binary. This will allow `time.LoadLocation` to work as expected, at the cost of binary size and a slight increase
 *   in execution time.
 */

// Returns the current time in UTC.
func GetUtcTime() time.Time {
	return time.Now().UTC()
}

// Returns the current local time.
func GetLocalTime() (string, error) {
	now, err := localtime.Now()
	if err != nil {
		return "", err
	}
	return now.Format(time.RFC3339), nil
}

// Returns the current time in a specified time zone.
func GetTimeInZone(tz string) (string, error) {
	now, err := localtime.NowInZone(tz)
	if err != nil {
		return "", err
	}
	return now.Format(time.RFC3339), nil
}

// Returns the local time zone identifier.
func GetLocalTimeZone() string {
	return localtime.GetTimeZone()

	// Alternatively, you can use the following:
	// return os.Getenv("TZ")
}

type TimeZoneInfo struct {
	StandardName   string
	StandardOffset string
	DaylightName   string
	DaylightOffset string
}

// Returns some basic information about the time zone specified.
func GetTimeZoneInfo(tz string) (*TimeZoneInfo, error) {
	loc, err := localtime.GetLocation(tz)
	if err != nil {
		return nil, err
	}

	janName, janOffset := time.Date(2024, 1, 1, 0, 0, 0, 0, loc).Zone()
	julName, julOffset := time.Date(2024, 7, 1, 0, 0, 0, 0, loc).Zone()

	var stdName, dltName string
	var stdOffset, dltOffset int
	if janOffset <= julOffset {
		stdName, stdOffset = janName, janOffset
		dltName, dltOffset = julName, julOffset
	} else {
		stdName, stdOffset = julName, julOffset
		dltName, dltOffset = janName, janOffset
	}

	info := &TimeZoneInfo{
		StandardName:   stdName,
		StandardOffset: formatOffset(stdOffset),
		DaylightName:   dltName,
		DaylightOffset: formatOffset(dltOffset),
	}

	return info, nil
}

// Formats the offset in hours and minutes (used by GetTimeZoneInfo).
func formatOffset(offset int) string {
	sign := "+"
	if offset < 0 {
		sign = "-"
		offset = -offset
	}
	offset /= 60
	hours := offset / 60
	minutes := offset % 60
	return fmt.Sprintf("%s%02d:%02d", sign, hours, minutes)
}
