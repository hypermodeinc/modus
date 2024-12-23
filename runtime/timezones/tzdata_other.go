//go:build !windows

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
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

func loadTimeZoneInfo(tz string) (*tzInfo, error) {
	tzFile := "/usr/share/zoneinfo/" + tz
	if _, err := os.Stat(tzFile); err != nil {
		return nil, fmt.Errorf("could not find timezone file: %v", err)
	}

	bytes, err := os.ReadFile(tzFile)
	if err != nil {
		return nil, fmt.Errorf("could not read timezone file: %v", err)
	}

	loc, err := time.LoadLocationFromTZData(tz, bytes)
	if err != nil {
		return nil, fmt.Errorf("could not load timezone data: %v", err)
	}

	info := &tzInfo{loc, bytes}
	return info, nil
}

func getSystemLocalTimeZone() (string, error) {
	// On Linux and macOS, we use the system default /etc/localtime file to get the time zone.
	// It is a symlink to the time zone file within the OS's copy of the IANA time zone database.
	// The time zone identifier is the path to the file relative to the zoneinfo folder.

	p, err := os.Readlink("/etc/localtime")
	if err == nil {
		segments := strings.Split(p, string(os.PathSeparator))
		for i := len(segments) - 1; i >= 0; i-- {
			if segments[i] == "zoneinfo" {
				return path.Join(segments[i+1:]...), nil
			}
		}
	}

	return "", errors.New("failed to determine system local time zone")
}
