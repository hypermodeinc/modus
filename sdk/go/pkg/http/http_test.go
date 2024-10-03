//go:build !wasip1

/*
 * Copyright 2024 Hypermode, Inc.
 */

package http_test

import (
	"reflect"
	"testing"

	"github.com/hypermodeAI/functions-go/pkg/http"
)

func TestFetchWithUrl(t *testing.T) {
	response, err := http.Fetch("https://example.com")
	if err != nil {
		t.Fatalf("Expected no error, but received: %s", err)
	}
	if !response.Ok() {
		t.Errorf("Expected OK, but received: %d %s", response.Status, response.StatusText)
	}

	expected := &http.Request{
		Url:    "https://example.com",
		Method: "GET",
	}
	values := http.FetchCallStack.Pop()
	if values == nil {
		t.Error("Expected a request, but none was found.")
	} else {
		if !reflect.DeepEqual(expected, values[0]) {
			t.Errorf("Expected request: %v, but received: %v", expected, values[0])
		}
	}
}

func TestFetchWithRequestObject(t *testing.T) {
	request := http.NewRequest("https://example.com")
	response, err := http.Fetch(request)
	if err != nil {
		t.Fatalf("Expected no error, but received: %s", err)
	}
	if !response.Ok() {
		t.Errorf("Expected OK, but received: %d %s", response.Status, response.StatusText)
	}

	expected := &http.Request{
		Url:    "https://example.com",
		Method: "GET",
	}
	values := http.FetchCallStack.Pop()
	if values == nil {
		t.Error("Expected a request, but none was found.")
	} else {
		if !reflect.DeepEqual(expected, values[0]) {
			t.Errorf("Expected request: %v, but received: %v", expected, values[0])
		}
	}
}

func TestFetchWithInvalidUrl(t *testing.T) {
	sizeBefore := http.FetchCallStack.Size()
	response, err := http.Fetch("invalid-url")
	if err == nil {
		t.Error("Expected an error, but received none")
	}

	sizeAfter := http.FetchCallStack.Size()
	if sizeBefore != sizeAfter {
		r := http.FetchCallStack.Pop()
		t.Errorf("Expected no request, but received: %v", r)
	}

	if response != nil {
		t.Errorf("Expected no response, but received: %v", response)
	}
}

func TestFetchWithOptions(t *testing.T) {
	response, err := http.Fetch("https://example.com", &http.RequestOptions{
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: []byte(`{"key": "value"}`),
	})
	if err != nil {
		t.Fatalf("Expected no error, but received: %s", err)
	}
	if !response.Ok() {
		t.Errorf("Expected OK, but received: %d %s", response.Status, response.StatusText)
	}

	expected := &http.Request{
		Url:    "https://example.com",
		Method: "POST",
		Headers: http.NewHeaders(map[string]string{
			"Content-Type": "application/json",
		}),
		Body: []byte(`{"key": "value"}`),
	}
	values := http.FetchCallStack.Pop()
	if values == nil {
		t.Errorf("Expected a request, but none was found")
	} else {
		if !reflect.DeepEqual(expected, values[0]) {
			t.Errorf("Expected request: %v, but received: %v", expected, values[0])
		}
	}
}

func TestFetchTooManyOptionsParameters(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected a panic, but no panic occurred")
		}
	}()
	_ = http.NewRequest("https://example.com", nil, nil)
}
