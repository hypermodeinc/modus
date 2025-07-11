//go:build !wasip1

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package console

import (
	"fmt"

	"github.com/hypermodeinc/modus/sdk/go/pkg/testutils"
)

var LogCallStack = testutils.NewCallStack()

func hostLogMessage(level, message string) {
	LogCallStack.Push(level, message)

	if level == "" {
		fmt.Println(message)
	} else {
		fmt.Printf("[%s] %s\n", level, message)
	}
}
