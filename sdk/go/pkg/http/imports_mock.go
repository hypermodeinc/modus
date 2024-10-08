//go:build !wasip1

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package http

import "github.com/hypermodeinc/modus/sdk/go/pkg/testutils"

var FetchCallStack = testutils.NewCallStack()

func fetch(request *Request) *Response {
	FetchCallStack.Push(request)

	return &Response{
		Status:     200,
		StatusText: "OK",
		Headers: NewHeaders(map[string]string{
			"Content-Type": "text/plain",
		}),
		Body: []byte("Hello, World!"),
	}
}
