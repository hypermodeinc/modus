/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang_test

import (
	"hypruntime/httpclient"
	"testing"
)

func TestHttpResponseHeaders(t *testing.T) {
	r := &httpclient.HttpResponse{
		Status:     200,
		StatusText: "OK",
		Headers: &httpclient.HttpHeaders{
			Data: map[string]*httpclient.HttpHeader{
				"content-type": {
					Name:   "Content-Type",
					Values: []string{"text/plain"},
				},
			},
		},
		Body: []byte("Hello, world!"),
	}

	if _, err := fixture.CallFunction(t, "testHttpResponseHeaders", r); err != nil {
		t.Fatal(err)
	}
}

func TestHttpHeaders(t *testing.T) {
	h := &httpclient.HttpHeaders{
		Data: map[string]*httpclient.HttpHeader{
			"content-type": {
				Name:   "Content-Type",
				Values: []string{"text/plain"},
			},
		},
	}

	if _, err := fixture.CallFunction(t, "testHttpHeaders", h); err != nil {
		t.Fatal(err)
	}
}

func TestHttpHeaderMap(t *testing.T) {
	m := map[string]*httpclient.HttpHeader{
		"content-type": {
			Name:   "Content-Type",
			Values: []string{"text/plain"},
		},
	}

	if _, err := fixture.CallFunction(t, "testHttpHeaderMap", m); err != nil {
		t.Fatal(err)
	}
}

func TestHttpHeader(t *testing.T) {
	h := httpclient.HttpHeader{
		Name:   "Content-Type",
		Values: []string{"text/plain"},
	}

	if _, err := fixture.CallFunction(t, "testHttpHeader", h); err != nil {
		t.Fatal(err)
	}
}
