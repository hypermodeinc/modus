/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package http

import (
	"encoding/json"
)

type Content struct {
	data []byte
}

func NewContent(value any) *Content {
	if value == nil {
		return nil
	}
	switch t := any(value).(type) {
	case []byte:
		return &Content{t}
	case *[]byte:
		return NewContent(*t)
	case string:
		return &Content{[]byte(t)}
	case *string:
		return NewContent(*t)
	case Content:
		return &t
	case *Content:
		return t
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return &Content{bytes}
}

func (c *Content) Bytes() []byte {
	return c.data
}

func (c *Content) Text() string {
	return string(c.data)
}

func (c *Content) JSON(result any) error {
	return json.Unmarshal(c.data, result)
}
