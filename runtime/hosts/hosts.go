/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package hosts

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

const HypermodeHost string = "hypermode"
const OpenAIHost string = "openai"

func GetHost(hostName string) (manifest.HostInfo, error) {
	if hostName == HypermodeHost {
		return manifest.HTTPHostInfo{Name: HypermodeHost}, nil
	}

	if host, ok := manifestdata.GetManifest().Hosts[hostName]; ok {
		return host, nil
	}

	return nil, fmt.Errorf("a host '%s' was not found", hostName)
}

func GetHttpHost(hostName string) (*manifest.HTTPHostInfo, error) {
	host, err := GetHost(hostName)
	if err != nil {
		return nil, err
	}

	if httpHost, ok := host.(manifest.HTTPHostInfo); ok {
		return &httpHost, nil
	}

	return nil, fmt.Errorf("host '%s' is not an HTTP host", hostName)
}

func GetHttpHostForUrl(url string) (*manifest.HTTPHostInfo, error) {

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

	// Find the HTTP host that matches the url
	// Either endpoint must match completely, or baseUrl must be a prefix of the url
	// (case insensitive comparison, either way)
	for _, host := range manifestdata.GetManifest().Hosts {
		if httpHost, ok := host.(manifest.HTTPHostInfo); ok {
			if httpHost.Endpoint != "" && strings.EqualFold(httpHost.Endpoint, url) {
				return &httpHost, nil
			} else if httpHost.BaseURL != "" && len(url) >= len(httpHost.BaseURL) && strings.EqualFold(httpHost.BaseURL, url[:len(httpHost.BaseURL)]) {
				return &httpHost, nil
			}
		}
	}

	return nil, fmt.Errorf("a host for url '%s' was not found in the manifest", url)
}

func PostToHostEndpoint[TResult any](ctx context.Context, host *manifest.HTTPHostInfo, payload any) (*utils.HttpResult[TResult], error) {
	if host.Endpoint == "" {
		return nil, fmt.Errorf("host endpoint is not defined")
	}

	bs := func(ctx context.Context, req *http.Request) error {
		return secrets.ApplyHostSecretsToHttpRequest(ctx, host, req)
	}

	return utils.PostHttp[TResult](ctx, host.Endpoint, payload, bs)
}
