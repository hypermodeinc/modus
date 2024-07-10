/*
 * Copyright 2024 Hypermode, Inc.
 */

package hosts

import (
	"context"
	"fmt"
	"net/http"
	urlpkg "net/url"
	"strings"

	"hmruntime/manifestdata"
	"hmruntime/secrets"
	"hmruntime/utils"

	"github.com/hypermodeAI/manifest"
)

const HypermodeHost string = "hypermode"
const OpenAIHost string = "openai"

func GetHost(hostName string) (manifest.HostInfo, error) {
	if hostName == HypermodeHost {
		return manifest.HTTPHostInfo{Name: HypermodeHost}, nil
	}

	if host, ok := manifestdata.Manifest.Hosts[hostName]; ok {
		return host, nil
	}

	return nil, fmt.Errorf("a host '%s' was not found", hostName)
}

func GetHttpHost(hostName string) (manifest.HTTPHostInfo, error) {
	host, err := GetHost(hostName)
	if err != nil {
		return manifest.HTTPHostInfo{}, err
	}

	if httpHost, ok := host.(manifest.HTTPHostInfo); ok {
		return httpHost, nil
	}

	return manifest.HTTPHostInfo{}, fmt.Errorf("host '%s' is not an HTTP host", hostName)
}

func GetHttpHostForUrl(url string) (manifest.HTTPHostInfo, error) {

	// Ensure the url is valid
	u, err := urlpkg.ParseRequestURI(url)
	if err != nil {
		return manifest.HTTPHostInfo{}, err
	}

	// Remove components not used for lookup
	u.User = nil
	u.RawQuery = ""
	u.Fragment = ""
	url = u.String()

	// Find the HTTP host that matches the url
	// Either endpoint must match completely, or baseUrl must be a prefix of the url
	// (case insensitive comparison, either way)
	for _, host := range manifestdata.Manifest.Hosts {
		if httpHost, ok := host.(manifest.HTTPHostInfo); ok {
			if httpHost.Endpoint != "" && strings.EqualFold(httpHost.Endpoint, url) {
				return httpHost, nil
			} else if httpHost.BaseURL != "" && len(url) >= len(httpHost.BaseURL) && strings.EqualFold(httpHost.BaseURL, url[:len(httpHost.BaseURL)]) {
				return httpHost, nil
			}
		}
	}

	return manifest.HTTPHostInfo{}, fmt.Errorf("a host for url '%s' was not found in the manifest", url)
}

func PostToHostEndpoint[TResult any](ctx context.Context, host manifest.HTTPHostInfo, payload any) (*utils.HttpResult[TResult], error) {
	if host.Endpoint == "" {
		return nil, fmt.Errorf("host endpoint is not defined")
	}

	bs := func(ctx context.Context, req *http.Request) error {
		return secrets.ApplyHostSecretsToHttpRequest(ctx, host, req)
	}

	return utils.PostHttp[TResult](ctx, host.Endpoint, payload, bs)
}
