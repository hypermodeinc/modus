/*
 * Copyright 2024 Hypermode, Inc.
 */

package httpclient

import (
	"hmruntime/plugins"
)

type HttpRequest struct {
	Url     string
	Method  string
	Headers HttpHeaders
	Body    []byte
}

func (r *HttpRequest) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "Request",
		Path: "~lib/@hypermode/functions-as/assembly/http/Request",
	}
}

type HttpResponse struct {
	Status     uint16
	StatusText string
	Headers    HttpHeaders
	Body       []byte
}

func (r *HttpResponse) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "Response",
		Path: "~lib/@hypermode/functions-as/assembly/http/Response",
	}
}

type HttpHeaders struct {
	Data map[string]HttpHeader
}

func (r *HttpHeaders) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "Headers",
		Path: "~lib/@hypermode/functions-as/assembly/http/Headers",
	}
}

type HttpHeader struct {
	Name   string
	Values []string
}

func (r *HttpHeader) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "Header",
		Path: "~lib/@hypermode/functions-as/assembly/http/Header",
	}
}
