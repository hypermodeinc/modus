/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
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
