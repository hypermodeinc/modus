//go:build !wasip1

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
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
