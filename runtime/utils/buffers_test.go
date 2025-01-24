/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package utils_test

import (
	"os"
	"testing"

	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewOutputBuffers(t *testing.T) {
	buffers := utils.NewOutputBuffers()

	outStream := buffers.StdOut()
	errStream := buffers.StdErr()

	assert.NotNil(t, outStream)
	assert.NotNil(t, errStream)

	assert.NotSame(t, outStream, errStream)
	assert.NotSame(t, os.Stdout, outStream)
	assert.NotSame(t, os.Stderr, errStream)

	assert.Equal(t, 0, outStream.Cap())
	assert.Equal(t, 0, errStream.Cap())

	outStream.Write([]byte("Hello, World!"))
	errStream.Write([]byte("Hello, Error!"))

	assert.Equal(t, "Hello, World!", outStream.String())
	assert.Equal(t, "Hello, Error!", errStream.String())
}
