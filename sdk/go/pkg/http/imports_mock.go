//go:build !wasip1

/*
 * Copyright 2024 Hypermode, Inc.
 */

package http

import "github.com/hypermodeAI/functions-go/pkg/testutils"

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
