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
	"errors"
	"net/url"
)

func Fetch[T *Request | string](requestOrUrl T, options ...*RequestOptions) (*Response, error) {

	if len(options) > 1 {
		panic("Too many arguments to Fetch")
	}

	var request *Request
	switch t := any(requestOrUrl).(type) {
	case *Request:
		request = t.Clone(options...)
	case string:
		request = NewRequest(t, options...)
	}

	if _, err := url.ParseRequestURI(request.Url); err != nil {
		return nil, errors.New("Invalid URL")
	}

	response := hostFetch(request)
	if response == nil {
		msg := "HTTP fetch failed. Check the logs for more information."
		return nil, errors.New(msg)
	}

	return response, nil
}
