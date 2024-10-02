/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_SendHttp(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello, World!")
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	content, err := sendHttp(req)
	if err != nil {
		t.Fatalf("Failed to get HTTP content: %v", err)
	}

	expected := "Hello, World!"
	if string(content) != expected {
		t.Errorf("Unexpected content. Got: %s, want: %s", content, expected)
	}
}

func Test_SendHttp_ErrorResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Something went wrong!", http.StatusInternalServerError)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	_, err = sendHttp(req)
	if err == nil {
		t.Error("Expected an error, but got nil")
	}

	expected := "HTTP error: 500 Internal Server Error\nSomething went wrong!\n"
	if err.Error() != expected {
		t.Errorf("Unexpected error message. Got: %s, want: %s", err.Error(), expected)
	}
}

func Test_PostHttp(t *testing.T) {
	type Payload struct {
		Message string `json:"message"`
	}

	type Response struct {
		Result string `json:"result"`
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Unexpected request method. Got: %s, want: %s", r.Method, http.MethodPost)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Unexpected content type. Got: %s, want: %s", r.Header.Get("Content-Type"), "application/json")
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}

		var payload Payload
		err = JsonDeserialize(body, &payload)
		if err != nil {
			t.Fatalf("Failed to deserialize request payload: %v", err)
		}

		expectedPayload := Payload{
			Message: "Hello, World!",
		}
		if !reflect.DeepEqual(payload, expectedPayload) {
			t.Errorf("Unexpected request payload. Got: %+v, want: %+v", payload, expectedPayload)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"result": "success"}`)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	url := server.URL
	payload := Payload{
		Message: "Hello, World!",
	}
	result, err := PostHttp[Response](context.Background(), url, payload, nil)
	if err != nil {
		t.Fatalf("Failed to make HTTP request: %v", err)
	}

	expectedResult := Response{
		Result: "success",
	}
	if !reflect.DeepEqual(result.Data, expectedResult) {
		t.Errorf("Unexpected response. Got: %+v, want: %+v", result, expectedResult)
	}
}

func Test_PostHttp_CustomContentType(t *testing.T) {
	const customContentType = "x-foo/bar"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") == customContentType {
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Invalid Content-Type", http.StatusBadRequest)
		}
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	bs := func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Content-Type", customContentType)
		return nil
	}

	url := server.URL
	_, err := PostHttp[string](context.Background(), url, nil, bs)
	if err != nil {
		t.Fatalf("Failed to make HTTP request: %v", err)
	}
}

func Test_PostHttp_StringResult(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `success`)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	url := server.URL
	result, err := PostHttp[string](context.Background(), url, nil, nil)
	if err != nil {
		t.Fatalf("Failed to make HTTP request: %v", err)
	}

	expectedResult := "success"
	if result.Data != expectedResult {
		t.Errorf("Unexpected response. Got: %s, want: %s", result, expectedResult)
	}
}

func Test_PostHttp_BytesResult(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte{1, 2, 3, 4, 5})
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	url := server.URL
	result, err := PostHttp[[]byte](context.Background(), url, nil, nil)
	if err != nil {
		t.Fatalf("Failed to make HTTP request: %v", err)
	}

	expectedResult := []byte{1, 2, 3, 4, 5}
	if !bytes.Equal(result.Data, expectedResult) {
		t.Errorf("Unexpected response. Got: %v, want: %v", result, expectedResult)
	}
}
