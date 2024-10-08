/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import "fmt"

type HttpResponse struct {
	Status     uint16
	StatusText string
	Headers    *HttpHeaders
	Body       []byte
}

type HttpHeaders struct {
	Data map[string]*HttpHeader
}

type HttpHeader struct {
	Name   string
	Values []string
}

func TestHttpResponseHeaders(r *HttpResponse) {
	if r == nil {
		fail("Response is nil")
	}

	h := r.Headers
	if h == nil {
		fail("Headers is nil")
	}

	if h.Data == nil {
		fail("Headers.Data is nil")
	}

	if len(h.Data) == 0 {
		fail("expected headers > 0, but got none")
	}

	fmt.Println("Headers:")
	for k, v := range h.Data {
		fmt.Printf("  %s: %+v\n", k, v)
		TestHttpHeader(v)
	}
}

func TestHttpHeaders(h *HttpHeaders) {
	if h == nil {
		fail("Headers is nil")
	}

	if h.Data == nil {
		fail("Headers.Data is nil")
	}

	if len(h.Data) == 0 {
		fail("expected headers > 0, but got none")
	}

	fmt.Println("Headers:")
	for k, v := range h.Data {
		fmt.Printf("  %s: %+v\n", k, v)
		TestHttpHeader(v)
	}
}

func TestHttpHeaderMap(m map[string]*HttpHeader) {
	if m == nil {
		fail("map is nil")
	}

	if len(m) == 0 {
		fail("expected headers > 0, but got none")
	}

	fmt.Println("Headers:")
	for k, v := range m {
		fmt.Printf("  %s: %+v\n", k, v)
		TestHttpHeader(v)
	}
}

func TestHttpHeader(h *HttpHeader) {
	if h == nil {
		fail("Header is nil")
	}

	if h.Name == "" {
		fail("Header.Name is empty")
	}

	if h.Values == nil {
		fail("Header.Values is nil")
	}

	fmt.Printf("Header: %+v\n", h)
}
