//go:build !wasip1

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package http

import "github.com/hypermodeinc/modus/sdk/go/pkg/testutils"

var FetchCallStack = testutils.NewCallStack()

func hostFetch(request *Request) *Response {
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
