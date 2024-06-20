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

func GetHTTPHost(hostName string) (manifest.HTTPHostInfo, error) {
	if hostName == HypermodeHost {
		return manifest.HTTPHostInfo{Name: HypermodeHost}, nil
	}

	host, ok := manifestdata.Manifest.Hosts[hostName]
	if ok && host.HostType() == manifest.HostTypeHTTP {
		return host.(manifest.HTTPHostInfo), nil
	}

	return manifest.HTTPHostInfo{}, fmt.Errorf("a http host '%s' was not found", hostName)
}

func GetHTTPHostForUrl(url string) (manifest.HTTPHostInfo, error) {
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

	// Find the host that matches the url
	// Either endpoint must match completely, or baseUrl must be a prefix of the url
	// (case insensitive comparison, either way)
	for _, h := range manifestdata.Manifest.Hosts {
		if h.HostType() != manifest.HostTypeHTTP {
			continue
		}

		host := h.(manifest.HTTPHostInfo)
		if host.Endpoint != "" && strings.EqualFold(host.Endpoint, url) {
			return host, nil
		} else if host.BaseURL != "" && len(url) >= len(host.BaseURL) && strings.EqualFold(host.BaseURL, url[:len(host.BaseURL)]) {
			return host, nil
		}
	}

	return manifest.HTTPHostInfo{}, fmt.Errorf("a host for url '%s' was not found in the manifest", url)
}

func PostToHostEndpoint[TResult any](ctx context.Context, host manifest.HTTPHostInfo, payload any) (
	*utils.HttpResult[TResult], error) {

	if host.Endpoint == "" {
		return nil, fmt.Errorf("host endpoint is not defined")
	}

	bs := func(ctx context.Context, req *http.Request) error {
		return secrets.ApplyHTTPHostSecrets(ctx, host, req)
	}

	return utils.PostHttp[TResult](ctx, host.Endpoint, payload, bs)
}
