/*
 * Copyright 2024 Hypermode, Inc.
 */

package httpclient

type HttpRequest struct {
	Url     string
	Method  string
	Headers HttpHeaders
	Body    []byte
}

type HttpResponse struct {
	Status     uint16
	StatusText string
	Headers    HttpHeaders
	Body       []byte
}

type HttpHeaders struct {
	Data map[string]HttpHeader
}

type HttpHeader struct {
	Name   string
	Values []string
}
