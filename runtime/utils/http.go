/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

var httpClient = &http.Client{}

type HttpError struct {
	StatusCode int
	Message    string
}

func (e *HttpError) Error() string {
	return "HTTP error: " + e.Message
}

func HttpClient() *http.Client {
	return httpClient
}

func sendHttp(req *http.Request) ([]byte, error) {
	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return body, &HttpError{
			StatusCode: response.StatusCode,
			Message:    response.Status,
		}
	}

	return body, nil
}

type HttpResult[T any] struct {
	Data      T
	StartTime time.Time
	EndTime   time.Time
}

func (r HttpResult[T]) Duration() time.Duration {
	return r.EndTime.Sub(r.StartTime)
}

func PostHttp[TResult any](ctx context.Context, url string, payload any, beforeSend func(context.Context, *http.Request) error) (*HttpResult[TResult], error) {
	var ct string
	var buf *bytes.Buffer

	switch payload := payload.(type) {
	case []byte:
		ct = "application/octet-stream"
		buf = bytes.NewBuffer(payload)
	case string:
		ct = "text/plain"
		buf = bytes.NewBuffer([]byte(payload))
	default:
		ct = "application/json"
		jsonPayload, err := JsonSerialize(payload)
		if err != nil {
			return nil, fmt.Errorf("error serializing payload: %w", err)
		}
		buf = bytes.NewBuffer(jsonPayload)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	if beforeSend != nil {
		err = beforeSend(ctx, req)
		if err != nil {
			return nil, err
		}
	}

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", ct)
	}

	startTime := GetTime()
	content, err := sendHttp(req)
	endTime := GetTime()

	// NOTE: Unlike most functions, the result and error are BOTH returned.
	// This is because some error messages are returned in the body of the response.

	var result TResult
	if content != nil {
		switch any(result).(type) {
		case []byte:
			result = any(content).(TResult)
		case string:
			result = any(string(content)).(TResult)
		default:
			if err := JsonDeserialize(content, &result); err != nil {
				return nil, fmt.Errorf("error deserializing response: %w", err)
			}
		}
	}

	return &HttpResult[TResult]{
		Data:      result,
		StartTime: startTime,
		EndTime:   endTime,
	}, err
}

func WriteJsonContentHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
