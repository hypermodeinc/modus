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
		return manifest.HostInfo{Name: HypermodeHost}, nil
	}

	host, ok := manifestdata.Manifest.Hosts[hostName]
	if ok {
		return host, nil
	}

	return manifest.HostInfo{}, fmt.Errorf("a host '%s' was not found", hostName)
}

func GetHostForUrl(url string) (manifest.HostInfo, error) {

	// Ensure the url is valid
	u, err := urlpkg.ParseRequestURI(url)
	if err != nil {
		return manifest.HostInfo{}, err
	}

	// Remove components not used for lookup
	u.User = nil
	u.RawQuery = ""
	u.Fragment = ""
	url = u.String()

	// Find the host that matches the url
	// Either endpoint must match completely, or baseUrl must be a prefix of the url
	// (case insensitive comparison, either way)
	for _, host := range manifestdata.Manifest.Hosts {
		if host.Endpoint != "" && strings.EqualFold(host.Endpoint, url) {
			return host, nil
		} else if host.BaseURL != "" && len(url) >= len(host.BaseURL) && strings.EqualFold(host.BaseURL, url[:len(host.BaseURL)]) {
			return host, nil
		}
	}

	return manifest.HostInfo{}, fmt.Errorf("a host for url '%s' was not found in the manifest", url)
}

func PostToHostEndpoint[TResult any](ctx context.Context, host manifest.HostInfo, payload any) (*utils.HttpResult[TResult], error) {
	if host.Endpoint == "" {
		return nil, fmt.Errorf("host endpoint is not defined")
	}

	bs := func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Content-Type", "application/json")
		return secrets.ApplyHostSecrets(ctx, host, req)
	}

	return utils.PostHttp[TResult](ctx, host.Endpoint, payload, bs)
}
