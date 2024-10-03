//go:build !wasip1

/*
 * Copyright 2024 Hypermode, Inc.
 */

package console

import (
	"fmt"

	"github.com/hypermodeAI/functions-go/pkg/testutils"
)

var LogCallStack = testutils.NewCallStack()

func log(level, message string) {
	LogCallStack.Push(level, message)

	if level == "" {
		fmt.Println(message)
	} else {
		fmt.Printf("[%s] %s\n", level, message)
	}
}
