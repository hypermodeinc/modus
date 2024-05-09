/*
 * Copyright 2024 Hypermode, Inc.
 */

package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func sendHttp(req *http.Request) ([]byte, error) {
	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", response.Status)
	}

	return io.ReadAll(response.Body)
}

func PostHttp[TResult any](url string, payload any, headers map[string]string) (TResult, error) {
	var result TResult
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
			return result, fmt.Errorf("error serializing payload: %w", err)
		}
		buf = bytes.NewBuffer(jsonPayload)
	}

	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return result, fmt.Errorf("error creating request: %w", err)
	}

	if _, ok := headers["Content-Type"]; !ok {
		req.Header.Set("Content-Type", ct)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	content, err := sendHttp(req)
	if err != nil {
		return result, fmt.Errorf("error sending HTTP request: %w", err)
	}

	switch any(result).(type) {
	case []byte:
		return any(content).(TResult), nil
	case string:
		return any(string(content)).(TResult), nil
	}

	err = JsonDeserialize(content, &result)
	if err != nil {
		return result, fmt.Errorf("error deserializing response: %w", err)
	}

	return result, nil
}

func WriteJsonContentHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
