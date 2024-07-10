/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"hmruntime/hosts"
	"hmruntime/logger"
	"hmruntime/plugins"
	"hmruntime/secrets"
	"hmruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

type httpRequest struct {
	Url     string
	Method  string
	Headers httpHeaders
	Body    []byte
}

func (r *httpRequest) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "Request",
		Path: "~lib/@hypermode/functions-as/assembly/http/Request",
	}
}

type httpResponse struct {
	Status     uint16
	StatusText string
	Headers    httpHeaders
	Body       []byte
}

func (r *httpResponse) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "Response",
		Path: "~lib/@hypermode/functions-as/assembly/http/Response",
	}
}

type httpHeaders struct {
	Data map[string]httpHeader
}

func (r *httpHeaders) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "Headers",
		Path: "~lib/@hypermode/functions-as/assembly/http/Headers",
	}
}

type httpHeader struct {
	Name   string
	Values []string
}

func (r *httpHeader) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "Header",
		Path: "~lib/@hypermode/functions-as/assembly/http/Header",
	}
}

func hostFetch(ctx context.Context, mod wasm.Module, pRequest uint32) (pResponse uint32) {
	var request httpRequest
	if err := readParam(ctx, mod, pRequest, &request); err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	host, err := hosts.GetHttpHostForUrl(request.Url)
	if err != nil {
		logger.Err(ctx, err).
			Str("url", request.Url).
			Msg("Error getting host.")
		return 0
	}

	body := bytes.NewBuffer(request.Body)
	req, err := http.NewRequestWithContext(ctx, request.Method, request.Url, body)
	if err != nil {
		logger.Err(ctx, err).Msg("Error creating request.")
		return 0
	}
	for _, header := range request.Headers.Data {
		req.Header[header.Name] = header.Values
	}

	if err := secrets.ApplyHostSecretsToHttpRequest(ctx, host, req); err != nil {
		logger.Err(ctx, err).
			Str("host", host.HostName()).
			Msg("Error applying host secrets.")
		return 0
	}

	resp, err := utils.HttpClient().Do(req)
	if err != nil {
		logger.Err(ctx, err).Msg("Error sending request.")
		return 0
	}
	defer resp.Body.Close()

	// Don't check status code here, just pass it back to the caller.

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading response body.")
		return 0
	}

	var response httpResponse
	response.Status = uint16(resp.StatusCode)
	response.StatusText = resp.Status[4:] // Remove the status code from the status text.

	response.Headers.Data = make(map[string]httpHeader, len(resp.Header))
	for name, values := range resp.Header {
		header := httpHeader{
			Name:   name,
			Values: values,
		}
		response.Headers.Data[strings.ToLower(name)] = header
	}

	response.Body = content

	offset, err := writeResult(ctx, mod, response)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}
	return offset
}
