/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package testutils

type CallStack struct {
	Items [][]any
}

func NewCallStack() *CallStack {
	return &CallStack{}
}

func (s *CallStack) Size() int {
	return len(s.Items)
}

func (s *CallStack) Push(values ...any) {
	s.Items = append(s.Items, values)
}

func (s *CallStack) Pop() []any {
	size := len(s.Items)
	if size == 0 {
		return nil
	}

	values := s.Items[size-1]
	s.Items = s.Items[:size]
	return values
}
