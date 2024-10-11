//go:build !wasip1

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
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
