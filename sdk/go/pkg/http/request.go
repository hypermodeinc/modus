/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package http

import (
	"encoding/json"
	"strings"
)

type Request struct {
	Url     string
	Method  string
	Headers *Headers
	Body    []byte
}

type RequestOptions struct {
	Method  string
	Headers any
	Body    any
}

func NewRequest(url string, options ...*RequestOptions) *Request {
	if len(options) > 1 {
		panic("Too many arguments to Clone")
	}

	req := &Request{Url: url, Method: "GET"}

	if len(options) == 1 && options[0] != nil {
		o := options[0]
		if o.Method != "" {
			req.Method = strings.ToUpper(o.Method)
		}
		if o.Headers != nil {
			switch v := o.Headers.(type) {
			case *Headers:
				req.Headers = v
			case Headers:
				req.Headers = &v
			case [][]string:
				req.Headers = NewHeaders(v)
			case map[string]string:
				req.Headers = NewHeaders(v)
			case map[string][]string:
				req.Headers = NewHeaders(v)
			default:
				panic("Invalid headers type")
			}
		}
		if o.Body != nil {
			req.Body = NewContent(o.Body).data
		}
	}

	return req
}

func (r *Request) Clone(options ...*RequestOptions) *Request {
	if len(options) > 1 {
		panic("Too many arguments to Clone")
	}

	ro := &RequestOptions{
		Method:  r.Method,
		Headers: r.Headers,
		Body:    &Content{r.Body},
	}

	if len(options) == 1 && options[0] != nil {
		o := options[0]
		if o.Method != "" {
			ro.Method = o.Method
		}
		if o.Headers != nil {
			switch v := o.Headers.(type) {
			case *Headers:
				if len(v.data) > 0 {
					ro.Headers = v
				}
			case Headers:
				if len(v.data) > 0 {
					ro.Headers = &v
				}
			case [][]string:
				if len(v) > 0 {
					ro.Headers = NewHeaders(v)
				}
			case map[string]string:
				if len(v) > 0 {
					ro.Headers = NewHeaders(v)
				}
			case map[string][]string:
				if len(v) > 0 {
					ro.Headers = NewHeaders(v)
				}
			default:
				panic("Invalid headers type")
			}
		}
		if o.Body != nil {
			ro.Body = o.Body
		}
	}

	return NewRequest(r.Url, ro)
}
func (r *Request) Text() string {
	return string(r.Body)
}

func (r *Request) JSON(result any) {
	if err := json.Unmarshal(r.Body, result); err != nil {
		panic(err)
	}
}
