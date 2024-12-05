/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package utils_test

import (
	"testing"

	"github.com/hypermodeinc/modus/runtime/utils"
)

func TestNewOutputBuffers(t *testing.T) {
	buffers := utils.NewOutputBuffers()

	if buffers.StdOut() == nil {
		t.Errorf("Expected StdOut buffer to be initialized, but got nil")
	}

	if buffers.StdErr() == nil {
		t.Errorf("Expected StdErr buffer to be initialized, but got nil")
	}
}
