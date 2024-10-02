/*
 * Copyright 2024 Hypermode, Inc.
 */

package httpclient

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"hypruntime/hosts"
	"hypruntime/secrets"
	"hypruntime/utils"
)

func HttpFetch(ctx context.Context, request *HttpRequest) (*HttpResponse, error) {
	host, err := hosts.GetHttpHostForUrl(request.Url)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(request.Body)
	req, err := http.NewRequestWithContext(ctx, request.Method, request.Url, body)
	if err != nil {
		return nil, err
	}

	if request.Headers != nil {
		for _, header := range request.Headers.Data {
			req.Header[header.Name] = header.Values
		}
	}

	if err := secrets.ApplyHostSecretsToHttpRequest(ctx, host, req); err != nil {
		return nil, err
	}

	resp, err := utils.HttpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Don't check status code here, just pass it back to the caller.

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	headers := make(map[string]*HttpHeader, len(resp.Header))
	for name, values := range resp.Header {
		header := &HttpHeader{
			Name:   name,
			Values: values,
		}
		headers[strings.ToLower(name)] = header
	}

	response := &HttpResponse{
		Status:     uint16(resp.StatusCode),
		StatusText: resp.Status[4:], // Remove the status code from the status text.
		Headers:    &HttpHeaders{Data: headers},
		Body:       content,
	}

	return response, nil
}
