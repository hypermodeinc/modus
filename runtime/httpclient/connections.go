/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package httpclient

import (
	"context"
	"fmt"
	"net/http"
	urlpkg "net/url"
	"strings"

	"github.com/hypermodeinc/modus/lib/manifest"
	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/secrets"
	"github.com/hypermodeinc/modus/runtime/utils"
)

const HypermodeConnectionName string = "hypermode"

func GetHttpConnectionInfo(name string) (*manifest.HTTPConnectionInfo, error) {
	if name == HypermodeConnectionName {
		return &manifest.HTTPConnectionInfo{Name: HypermodeConnectionName}, nil
	}

	if c, ok := manifestdata.GetManifest().Connections[name]; !ok {
		return nil, fmt.Errorf("connection [%s] was not found", name)
	} else if c, ok := c.(manifest.HTTPConnectionInfo); !ok {
		return nil, fmt.Errorf("[%s] is not an HTTP connection", name)
	} else {
		return &c, nil
	}
}

func GetHttpConnectionForUrl(url string) (*manifest.HTTPConnectionInfo, error) {

	// Ensure the url is valid
	u, err := urlpkg.ParseRequestURI(url)
	if err != nil {
		return nil, err
	}

	// Remove components not used for lookup
	u.User = nil
	u.RawQuery = ""
	u.Fragment = ""
	url = u.String()

	// Find the HTTP connection that matches the url
	// Either endpoint must match completely, or baseUrl must be a prefix of the url
	// (case insensitive comparison, either way)
	for _, connection := range manifestdata.GetManifest().Connections {
		if httpConnection, ok := connection.(manifest.HTTPConnectionInfo); ok {
			if httpConnection.Endpoint != "" && strings.EqualFold(httpConnection.Endpoint, url) {
				return &httpConnection, nil
			} else if httpConnection.BaseURL != "" && len(url) >= len(httpConnection.BaseURL) && strings.EqualFold(httpConnection.BaseURL, url[:len(httpConnection.BaseURL)]) {
				return &httpConnection, nil
			}
		}
	}

	return nil, fmt.Errorf("a connection for url [%s] was not found in the manifest", url)
}

func PostToConnectionEndpoint[TResult any](ctx context.Context, connection *manifest.HTTPConnectionInfo, payload any) (*utils.HttpResult[TResult], error) {
	if connection.Endpoint == "" {
		return nil, fmt.Errorf("connection endpoint is not defined")
	}

	bs := func(ctx context.Context, req *http.Request) error {
		return secrets.ApplySecretsToHttpRequest(ctx, connection, req)
	}

	return utils.PostHttp[TResult](ctx, connection.Endpoint, payload, bs)
}
