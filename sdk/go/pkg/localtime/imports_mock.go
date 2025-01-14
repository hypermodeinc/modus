//go:build !wasip1

/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package localtime

import (
	"time"

	"github.com/hypermodeinc/modus/sdk/go/pkg/testutils"
)

var GetLocalTimeCallStack = testutils.NewCallStack()
var GetTimeInZoneCallStack = testutils.NewCallStack()
var GetTimeLocationCallStack = testutils.NewCallStack()

func hostGetLocalTime() (time.Time, error) {
	GetLocalTimeCallStack.Push()
	return time.Now(), nil
}

func hostGetTimeInZone(tz string) (time.Time, error) {
	GetTimeInZoneCallStack.Push(tz)
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, err
	}

	return time.Now().In(loc), nil
}

func hostGetTimeLocation(tz string) (*time.Location, error) {
	GetTimeLocationCallStack.Push()
	return time.LoadLocation(tz)
}
