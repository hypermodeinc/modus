/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
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

func (r *Response) Text() string {
	return string(r.Body)
}

func (r *Response) JSON(result any) {
	if err := json.Unmarshal(r.Body, result); err != nil {
		panic(err)
	}
}
