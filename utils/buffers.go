/*
 * Copyright 2024 Hypermode, Inc.
 */

package utils

import "bytes"

type OutputBuffers struct {
	StdOut bytes.Buffer
	StdErr bytes.Buffer
}
