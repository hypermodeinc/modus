/*
 * Copyright 2024 Hypermode, Inc.
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
