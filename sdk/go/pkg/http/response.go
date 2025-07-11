/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package http

import (
	"encoding/json"
)

type Response struct {
	Status     uint16
	StatusText string
	Headers    *Headers
	Body       []byte
}

func (r *Response) Ok() bool {
	return r.Status >= 200 && r.Status < 300
}

func (r *Response) Bytes() []byte {
	return r.Body
}

func (r *Response) Text() string {
	return string(r.Body)
}

func (r *Response) JSON(result any) error {
	return json.Unmarshal(r.Body, result)
}
