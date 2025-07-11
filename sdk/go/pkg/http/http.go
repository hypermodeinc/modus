/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
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
		return nil, errors.New("invalid URL")
	}

	response := hostFetch(request)
	if response == nil {
		return nil, errors.New("HTTP fetch failed - check the logs for more information")
	}

	return response, nil
}
