/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang_test

import (
	"hypruntime/httpclient"
	"testing"
)

func TestHttpResponseHeaders(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

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

	if _, err := f.InvokeFunction("testHttpResponseHeaders", r); err != nil {
		t.Fatal(err)
	}
}

func TestHttpHeaders(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	h := &httpclient.HttpHeaders{
		Data: map[string]*httpclient.HttpHeader{
			"content-type": {
				Name:   "Content-Type",
				Values: []string{"text/plain"},
			},
		},
	}

	if _, err := f.InvokeFunction("testHttpHeaders", h); err != nil {
		t.Fatal(err)
	}
}

func TestHttpHeaderMap(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	m := map[string]*httpclient.HttpHeader{
		"content-type": {
			Name:   "Content-Type",
			Values: []string{"text/plain"},
		},
	}

	if _, err := f.InvokeFunction("testHttpHeaderMap", m); err != nil {
		t.Fatal(err)
	}
}

func TestHttpHeader(t *testing.T) {
	t.Parallel()

	f := NewGoWasmTestFixture(t)
	defer f.Close()

	h := httpclient.HttpHeader{
		Name:   "Content-Type",
		Values: []string{"text/plain"},
	}

	if _, err := f.InvokeFunction("testHttpHeader", h); err != nil {
		t.Fatal(err)
	}
}
