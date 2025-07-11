/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package timezones

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/puzpuzpuz/xsync/v4"
)

var systemTimeZone string
var tzDataCache = *xsync.NewMap[string, []byte]()
var tzLocationCache = *xsync.NewMap[string, *time.Location]()

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

func GetLocation(ctx context.Context, tz string) (*time.Location, error) {
	if tz == "" {
		return nil, errors.New("time zone cannot be empty")
	}

	loc, err := getTimeZoneLocation(tz)
	if err != nil {
		return nil, fmt.Errorf("failed to load time zone location for %s: %w", tz, err)
	}

	return loc, nil
}

func GetTimeZoneData(ctx context.Context, tz, format string) ([]byte, error) {
	if tz == "" {
		return nil, errors.New("time zone cannot be empty")
	}

	// only support tzif format for now
	// we can expand this to support other formats as needed
	if format != "tzif" {
		return nil, fmt.Errorf("unsupported time zone format: %s", format)
	}

	data, err := getTimeZoneData(tz)
	if err != nil {
		return nil, fmt.Errorf("failed to load time zone data for %s: %w", tz, err)
	}

	return data, nil
}

func getTimeZoneLocation(tz string) (*time.Location, error) {
	if loc, ok := tzLocationCache.Load(tz); ok {
		return loc, nil
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, err
	}

	tzLocationCache.Store(tz, loc)
	return loc, nil
}

func getTimeZoneData(tz string) ([]byte, error) {
	if data, ok := tzDataCache.Load(tz); ok {
		return data, nil
	}

	data, err := loadTimeZoneData(tz)
	if err != nil {
		return nil, err
	}

	tzDataCache.Store(tz, data)
	return data, nil
}
