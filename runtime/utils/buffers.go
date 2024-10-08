/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package utils

import "bytes"

type OutputBuffers interface {
	StdOut() *bytes.Buffer
	StdErr() *bytes.Buffer
}

func NewOutputBuffers() OutputBuffers {
	return &outputBuffers{
		stdOut: &bytes.Buffer{},
		stdErr: &bytes.Buffer{},
	}
}

type outputBuffers struct {
	stdOut *bytes.Buffer
	stdErr *bytes.Buffer
}

func (b *outputBuffers) StdOut() *bytes.Buffer {
	return b.stdOut
}

func (b *outputBuffers) StdErr() *bytes.Buffer {
	return b.stdErr
}
