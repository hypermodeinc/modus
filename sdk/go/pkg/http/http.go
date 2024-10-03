/*
 * Copyright 2024 Hypermode, Inc.
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

	response := fetch(request)
	if response == nil {
		msg := "HTTP fetch failed. Check the Hypermode logs for more information."
		return nil, errors.New(msg)
	}

	return response, nil
}
