/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang_test

import (
	"testing"

	"github.com/hypermodeinc/modus/runtime/httpclient"
)

func TestHttpResponseHeaders(t *testing.T) {
	fnName := "testHttpResponseHeaders"
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

	if _, err := fixture.CallFunction(t, fnName, r); err != nil {
		t.Error(err)
	}
}

func TestHttpHeaders(t *testing.T) {
	fnName := "testHttpHeaders"
	h := &httpclient.HttpHeaders{
		Data: map[string]*httpclient.HttpHeader{
			"content-type": {
				Name:   "Content-Type",
				Values: []string{"text/plain"},
			},
		},
	}

	if _, err := fixture.CallFunction(t, fnName, h); err != nil {
		t.Error(err)
	}
}

func TestHttpHeaderMap(t *testing.T) {
	fnName := "testHttpHeaderMap"
	m := map[string]*httpclient.HttpHeader{
		"content-type": {
			Name:   "Content-Type",
			Values: []string{"text/plain"},
		},
	}

	if _, err := fixture.CallFunction(t, fnName, m); err != nil {
		t.Error(err)
	}
}

func TestHttpHeader(t *testing.T) {
	fnName := "testHttpHeader"
	h := httpclient.HttpHeader{
		Name:   "Content-Type",
		Values: []string{"text/plain"},
	}

	if _, err := fixture.CallFunction(t, fnName, h); err != nil {
		t.Error(err)
	}
}
