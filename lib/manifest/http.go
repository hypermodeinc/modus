/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package manifest

import (
	"regexp"
)

const ConnectionTypeHTTP ConnectionType = "http"

var templateRegex = regexp.MustCompile(`{{\s*(?:base64\((.+?):(.+?)\)|(.+?))\s*}}`)

type HTTPConnectionInfo struct {
	Name            string            `json:"-"`
	Type            ConnectionType    `json:"type"`
	Endpoint        string            `json:"endpoint"`
	BaseURL         string            `json:"baseURL"`
	Headers         map[string]string `json:"headers"`
	QueryParameters map[string]string `json:"queryParameters"`
}

func (info HTTPConnectionInfo) ConnectionName() string {
	return info.Name
}

func (info HTTPConnectionInfo) ConnectionType() ConnectionType {
	return info.Type
}

func (info HTTPConnectionInfo) Hash() string {
	return computeHash(info.Name, info.Type, info.Endpoint, info.BaseURL)
}

func (info HTTPConnectionInfo) Variables() []string {
	cap := 2 * (len(info.Headers) + len(info.QueryParameters))
	set := make(map[string]bool, cap)
	results := make([]string, 0, cap)

	for _, header := range info.Headers {
		vars := extractVariables(header)
		for _, v := range vars {
			if _, ok := set[v]; !ok {
				set[v] = true
				results = append(results, v)
			}
		}
	}

	for _, v := range info.QueryParameters {
		vars := extractVariables(v)
		for _, v := range vars {
			if _, ok := set[v]; !ok {
				set[v] = true
				results = append(results, v)
			}
		}
	}

	return results
}
