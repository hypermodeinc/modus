/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package httpclient

type HttpRequest struct {
	Url     string
	Method  string
	Headers *HttpHeaders
	Body    []byte
}

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
